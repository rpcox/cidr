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

func main() {
	mask    := flag.Int("mask", 32, "Summarization mask")
	version := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *version {
		ShowVersion()
	}

	c := make(map[int]int)
	bitMask := ( 0xFFFFFFFF & ( 0xFFFFFFFF << (32 - *mask)))

	list := flag.Args()
	f, err := os.Open(list[0])
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

	var keys []int
	for k := range c {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, v := range keys {
		fmt.Fprintf(os.Stdout, "%6d\t-  %s/%d\n", c[v], Int2Ip(v), *mask)
	}

}
