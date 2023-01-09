package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func ShowVersion() {
	ver := "0.1.0"
	fmt.Fprintf(os.Stdout, "summarize v%s\n", ver)
	os.Exit(0)
}

func Ip2Int(ip string) int {
	var integer int
	quads := strings.Split(ip, ".")

	shift := 24
	for _, quad := range quads {
		i, _ := strconv.Atoi(quad)
		tmp := i << shift
		integer += tmp
		shift -= 8
	}

	return integer
}

func Int2Ip(integer int) string {
	s := make([]string, 4)

	for i := 0; i < 4; i++ {
		tmp := integer & 0xFF
		quad := strconv.Itoa(tmp)
		s[3 - i] = quad
		integer = integer >> 8
	}

	return strings.Join(s, ".")
}

func Summarize(mask int, file string) map[int]int {
	c := make(map[int]int)
	bitMask := ( 0xFFFFFFFF & ( 0xFFFFFFFF << (32 - mask)))

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		item := scanner.Text()
		i := Ip2Int(item)
		i = i & bitMask
		c[i]++
	}
	f.Close()

	return c
}

func ShowSummary(c map[int]int, mask int) {
	var keys []int
	for k := range c {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, v := range keys {
		fmt.Fprintf(os.Stdout, "%6d\t-  %s/%d\n", c[v], Int2Ip(v), mask)
	}
}

func main() {
	sumCmd   := flag.NewFlagSet("summarize", flag.ExitOnError)
	sumMask  := sumCmd.Int("mask", 32, "Length of summarization mask")
	rangeCmd := flag.NewFlagSet("to_range", flag.ExitOnError)
	verCmd   := flag.NewFlagSet("version", flag.ExitOnError)

	switch os.Args[1] {
	case "to_range":
		rangeCmd.Parse(os.Args[2:])
		fmt.Println(rangeCmd.Args())
	case "summarize":
		sumCmd.Parse(os.Args[2:])
		file := sumCmd.Args()
		m := Summarize(*sumMask, file[0])
		ShowSummary(m, *sumMask)
	case "version":
		verCmd.Parse(os.Args[2:])
		ShowVersion()
	default:
		fmt.Println("Invalid selection")
		os.Exit(1)
	}
}
//SDG
