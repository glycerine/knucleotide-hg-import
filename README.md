# k-nucleotide in Go

I've read this on [Benchmarkgames](http://benchmarksgame.alioth.debian.org/u64q/compare.php?lang=java&lang2=go):

![Java vs. Go](https://bytebucket.org/s_l_teichmann/knucleotide/raw/default/images/java-go.png)

and can't believe it. So I wrote my own Go version.

## Reproducing the old timings on my hardware

For reference the timings of the current upstream versions on my laptop (GNU/Linux, Intel(R) Core(TM) i7-5600U CPU @ 2.60GHz):

* Java: `/usr/bin/time -p java -server knucleotide < knucleotide-input25000000.txt` -> `real 5.80`
* Go: `/usr/bin/time -p ./knucleotide.go-3 < knucleotide-input25000000.txt` -> `real 16.86`

So the Java version is ~2.9 times faster then the Go one. This resembles the scale of the reported measurement.

## Timings of my new version

### On my laptop:

* `/usr/bin/time -p ./knucleotide < knucleotide-input25000000.txt`
* `/usr/bin/time -p java -XX:+TieredCompilation -XX:+AggressiveOpts -server knucleotide < knucleotide-input25000000.txt`

Values are best of 5.

* My Go version: `real 5.37`
* Java: `real 5.78`

So on this machine my Go version is ~1.076 times as fast as the Java one.

On this machine the values of the Java version build a broader range.  
Comparison of the mean values:

* My Go version: `5.406`
* Java: `6.202`

Go is ~1.147 as fast the Java one.

### On my desktop system (GNU/Linux, Intel(R) Core(TM) i7 CPU 860 @ 2.80GHz):

Calls as above.

Values are best of 5.

* My Go version: `real 3.90`
* Java: `real 4.79`

So on this machine my Go version is ~1.2 times as fast as the Java one.

The values for Java are more stable on this machine.

## Java and Go versions

    $ java -version
    openjdk version "1.8.0_102"
    OpenJDK Runtime Environment (build 1.8.0_102-b14)
    OpenJDK 64-Bit Server VM (build 25.102-b14, mixed mode)

    $ go version
    go version go1.7 linux/amd64

## Build

You need a working Go 1.7 environment.

    $ go get bitbucket.org/s_l_teichmann/knucleotide

# License
This is Free software covered by the terms of the MIT [LICENSE](LICENSE)  
(c) 2016 by Sascha L. Teichmann
