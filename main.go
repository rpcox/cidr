package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

func ShowVersion() {
	ver := "0.1.2"
	fmt.Fprintf(os.Stdout, "summarize v%s\n", ver)
	os.Exit(0)
}

func IP2Int(ip string) (int, error) {
	var integer int

	if net.ParseIP(ip) == nil {
		return 0, fmt.Errorf("Invalid IPv4 address: %s", ip)
	}

	quads := strings.Split(ip, ".")

	shift := 24
	for _, quad := range quads {
		i, _ := strconv.Atoi(quad)
		tmp := i << shift
		integer += tmp
		shift -= 8
	}

	return integer, nil
}

func Int2IP(integer int) string {
	s := make([]string, 4)

	for i := 0; i < 4; i++ {
		tmp := integer & 0xFF
		quad := strconv.Itoa(tmp)
		s[3 - i] = quad
		integer = integer >> 8
	}

	return strings.Join(s, ".")
}

type Range struct {
	Submit     string
	Block      string
	Netmask    string
	Compliment string
	FirstIP    string
        LastIP     string
        Count      int
}

func ParseBlock(b string) (int, int, error) {
	s := strings.Split(b, "/")
	n := len(s)
	if n < 2 || n > 2 {
		return 0, 0, fmt.Errorf("Invalid format: %s", b)
	}

	i, err := IP2Int(s[0])
	if err != nil {
		return 0, 0, err
	}

	m, err := strconv.Atoi(s[1])
	if err != nil {
		return 0, 0, err
	}

	if m < 8 || m > 32 {
		return 0, 0, fmt.Errorf("Mask out of range: %d", m)
	}

	return i, m, nil
}

func ToRange(blocks []string) []Range {
	var list []Range

	for _, v := range blocks {
		var r Range

		i, mask, err := ParseBlock(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}

		r.Submit = v

		netmask := ( 0xFFFFFFFF & ( 0xFFFFFFFF << (32 - mask)))
		compliment := netmask ^ 0xFFFFFFFF
		firstIp := i & netmask
		lastIp := i | compliment

		r.Netmask = Int2IP(netmask)
		r.Compliment = Int2IP(compliment)
		r.FirstIP = Int2IP(firstIp)
		r.LastIP = Int2IP(lastIp)
		r.Count = (lastIp - firstIp) + 1
		r.Block = fmt.Sprintf("%s/%d", r.FirstIP, mask)
		list = append(list, r)
	}

	return list
}

func PrintRange(r []Range) {
	fmt.Println()
	for _, v := range r {
		fmt.Fprintf(os.Stdout, " Submitted : %s\n", v.Submit)
		fmt.Fprintf(os.Stdout, "     Block : %s\n", v.Block)
		fmt.Fprintf(os.Stdout, "   Netmask : %s\n", v.Netmask)
		fmt.Fprintf(os.Stdout, "Compliment : %s\n", v.Compliment)
		fmt.Fprintf(os.Stdout, "  First IP : %s\n", v.FirstIP)
		fmt.Fprintf(os.Stdout, "   Last IP : %s\n", v.LastIP)
		fmt.Fprintf(os.Stdout, "     Count : %d\n", v.Count)
		fmt.Println()
	}
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
		i, err := IP2Int(item)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		i = i & bitMask
		c[i]++
	}
	f.Close()

	return c
}

func PrintSummary(c map[int]int, mask int) {
	var keys []int
	for k := range c {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, v := range keys {
		fmt.Fprintf(os.Stdout, "%6d\t-  %s/%d\n", c[v], Int2IP(v), mask)
	}
}

func PrintUsage() {
	doc := `

  NAME
	vlsm - variable length subnet masking

  SYNOPSIS
	vlsm SUBCOMMAND [COMMAND OPTS] [COMMAND ARGS]

  DESCRIPTION
	vlsm is a tool variable length subnet masking

    SUBCOMMANDS
	summarize - summarize a list of IP addresses contained in a file. The -mask option
        	is required.

		Example: vlsm summarize -mask 23 myIPList.txt

	to_range - take a subnet string in format IP/MASK and print the characteristics of
		that subnet (i.e., first IP, last IP, etc)

		Example: vlsm to_range 10.23.45.67/23
`

	fmt.Println(doc)
}

func main() {
	sumCmd   := flag.NewFlagSet("summarize", flag.ExitOnError)
	sumMask  := sumCmd.Int("mask", 32, "Length of summarization mask")
	rangeCmd := flag.NewFlagSet("to_range", flag.ExitOnError)
	verCmd   := flag.NewFlagSet("version", flag.ExitOnError)

	if flag.Arg(1) == "" && os.Args[1] != "version" {
		PrintUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "to_range":
		rangeCmd.Parse(os.Args[2:])
		r := ToRange(rangeCmd.Args())
		PrintRange(r)
	case "summarize":
		sumCmd.Parse(os.Args[2:])
		if *sumMask <= 7 || *sumMask > 32 {
			log.Fatal("mask must be > 8 or <= 32")
		}
		file := sumCmd.Args()
		m := Summarize(*sumMask, file[0])
		PrintSummary(m, *sumMask)
	case "version":
		verCmd.Parse(os.Args[2:])
		ShowVersion()
	default:
		fmt.Println("Invalid subcommand")
		PrintUsage()
		os.Exit(1)
	}
}
//SDG
