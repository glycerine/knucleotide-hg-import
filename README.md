# k-nucleotide in Go

I've read this on [Benchmarkgames](http://benchmarksgame.alioth.debian.org/u64q/compare.php?lang=java&lang2=go):

![Java vs. Go](https://bytebucket.org/s_l_teichmann/knucleotide/raw/default/images/java-go.png)

and can't believe it. So I wrote my own Go version.

## Reproducing the timings on my hardware

The timings of the current upstream versions on my laptop (GNU/Linux, Intel(R) Core(TM) i7-5600U CPU @ 2.60GHz):

* Java: `/usr/bin/time -p java -server knucleotide < knucleotide-input25000000.txt` -> `real 5.80`
* Go: `/usr/bin/time -p ./knucleotide.go-3 < knucleotide-input25000000.txt` -> `real 16.86`

So the Java version is ~2.9 times faster then the Go one. This resembles the scale of the reported measurement.

## Timing of my version

* `/usr/bin/time -p ./knucleotide < knucleotide-input25000000.txt` -> `real 6.09`

So Java is only ~1.05 times as fast as the Go version.  
This is in the expected range. Further optimization is possible (by tuning the hash map implementation).

## Build

You need a working Go 1.7 environment. No other external library needed.

    $ go get bitbucket.org/s_l_teichmann/knucleotide

# License
This is Free software covered by the terms of the Apache 2.0 [LICENSE](LICENSE)  
(c) 2016 by Sascha L. Teichmann