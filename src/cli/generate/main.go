package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func genFloat32(r *rand.Rand, min, max float64) float32 {
	v := r.Float64()
	v *= max - min
	v += min
	return float32(v)
}

func main() {
	var help bool
	flag.BoolVar(&help, "help", false, "show the usage")

	var dimension, kline uint
	var output string

	flag.UintVar(&dimension, "dimension", 1024, "the dimension of each vector")
	flag.UintVar(&kline, "kline", 100, "kilo lines of vectors")
	flag.StringVar(&output, "output", "output.csv", "the output file of test data")

	var min, max float64
	flag.Float64Var(&min, "min", -1, "min value")
	flag.Float64Var(&max, "max", 1, "min value")

	flag.Parse()
	if help {
		fmt.Printf("Usage: of %s:\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	if output == "" {
		fmt.Println("The `output` flag is empty.")
		os.Exit(-1)
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))

	outputf, err := os.Create(output)
	if err != nil {
		fmt.Println("Create output file error:", err)
		os.Exit(-1)
	}
	defer outputf.Close()

	buff := bufio.NewWriter(outputf)
	defer buff.Flush()

	for n, maxline := 0, int(kline)*1000; n < maxline; n++ {
		for d := 0; d < int(dimension); d++ {
			if d != 0 {
				fmt.Fprint(buff, ",")
			}
			fmt.Fprintf(buff, "%.10f", genFloat32(r, min, max))
		}
		fmt.Fprint(buff, "\n")

		if n%100 == 0 {
			fmt.Printf("generated %d lines...\n", n)
		}
	}
}
