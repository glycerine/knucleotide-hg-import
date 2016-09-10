// The Computer Language Benchmarks Game
// http://benchmarksgame.alioth.debian.org/
//
// (c) 2016 by Sascha L. Teichmann
// This is Free Software covered by the terms of the MIT license.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"

	"bitbucket.org/s_l_teichmann/fastmap"
)

var (
	three           = []byte(">THREE")
	fragmentLengths = []int{1, 2, 3, 4, 6, 12, 18}
	nucleotides     = []byte{'A', 'C', 'G', 'T'}
	codes           = [256]byte{
		'A': 0,
		'a': 0,
		'C': 1,
		'c': 1,
		'G': 2,
		'g': 2,
		'T': 3,
		't': 3,
	}
)

type result struct {
	c      *sync.Cond
	m      *fastmap.Hash64uToInt
	keyLen int
	sync.Mutex
}

func newResult(keyLen int) *result {
	r := &result{keyLen: keyLen}
	r.c = sync.NewCond(r)
	return r
}

func (r *result) getM() (h *fastmap.Hash64uToInt) {
	r.Lock()
	for r.m == nil {
		r.c.Wait()
	}
	h = r.m
	r.Unlock()
	return
}

func (r *result) setM(m *fastmap.Hash64uToInt) {
	r.Lock()
	r.m = m
	r.c.Signal()
	r.Unlock()
}

func read(r io.Reader) ([]byte, error) {
	s := bufio.NewScanner(r)
	if s.Scan() {
		for !bytes.HasPrefix(s.Bytes(), three) {
			s.Scan()
		}
	}
	var buf bytes.Buffer
	for s.Scan() {
		b := s.Bytes()
		if len(b) > 0 && b[0] == '>' {
			break
		}
		buf.Write(encode(b))
	}
	return buf.Bytes(), s.Err()
}

func encode(seq []byte) []byte {
	for i, b := range seq {
		seq[i] = codes[b]
	}
	return seq
}

func key(arr []byte) uint64 {
	var k uint64
	for _, v := range arr {
		k = (k << 2) | uint64(v)
	}
	return k
}

func createFragmentMap(seq []byte, ofs, length int) *fastmap.Hash64uToInt {
	m := fastmap.NewHash64uToInt()
	lastIndex := len(seq) - length + 1
	for i := ofs; i < lastIndex; i += length {
		// Manually inlined for performance
		// m.Add(key(seq[i:i+length]), 1)
		var k uint64
		for _, v := range seq[i : i+length] {
			k = (k << 2) | uint64(v)
		}
		m.Inc(k)
	}
	return m
}

func (r *result) add(o *result) {
	rm := r.getM()
	o.getM().Visit(func(k uint64, v int) {
		rm.Add(k, v)
	})
}

type keyFreq struct {
	key string
	cnt int
}

type sortByFreq []keyFreq

func (s sortByFreq) Len() int {
	return len(s)
}

func (s sortByFreq) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByFreq) Less(i, j int) bool {
	if d := s[i].cnt - s[j].cnt; d != 0 {
		return d > 0
	}
	return s[i].key < s[j].key
}

func keyToString(k uint64, length int) string {
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[length-i-1] = nucleotides[k&0x3]
		k >>= 2
	}
	return string(res)
}

func writeFrequencies(total int, frequencies *result) string {
	m := frequencies.getM()
	freq := make(sortByFreq, m.Size())
	i := 0
	m.Visit(func(k uint64, v int) {
		freq[i] = keyFreq{keyToString(k, frequencies.keyLen), v}
		i++
	})
	sort.Sort(freq)
	var buf bytes.Buffer
	for _, f := range freq {
		fmt.Fprintf(&buf, "%s %.3f\n", f.key,
			float32(f.cnt*100)/float32(total))
	}
	return buf.String()
}

func createFragments(seq []byte) []*result {

	type job struct {
		ofs int
		r   *result
	}

	jobs := make(chan job)

	for i, n := 0, runtime.NumCPU(); i < n; i++ {
		go func() {
			for j := range jobs {
				j.r.setM(createFragmentMap(seq, j.ofs, j.r.keyLen))
			}
		}()
	}

	var results []*result

	for _, l := range fragmentLengths {
		for i := 0; i < l; i++ {
			r := newResult(l)
			results = append(results, r)
			jobs <- job{ofs: i, r: r}
		}
	}
	close(jobs)
	return results
}

func writeCount(results []*result, fragment string) string {

	ks := encode([]byte(fragment))
	k := key(ks)
	count := 0
	for _, r := range results {
		if r.keyLen == len(ks) {
			count += r.getM().Get(k)
		}
	}
	return fmt.Sprintf("%d\t%s", count, fragment)
}

func main() {

	sequence, err := read(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	results := createFragments(sequence)

	var buf bytes.Buffer

	fmt.Fprintln(&buf, writeFrequencies(len(sequence), results[0]))
	results[1].add(results[2])
	fmt.Fprintln(&buf, writeFrequencies(len(sequence)-1, results[1]))

	for _, fragment := range []string{
		"GGT", "GGTA", "GGTATT",
		"GGTATTTTAATT",
		"GGTATTTTAATTTATAGT"} {
		fmt.Fprintln(&buf, writeCount(results, fragment))
	}

	buf.WriteTo(os.Stdout)
}
