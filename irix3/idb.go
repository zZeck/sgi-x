package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type entry struct {
	ty      rune
	mode    uint
	user    string
	group   string
	path    string
	sum     uint64
	size    uint64
	cmpsize uint64
	offset  uint64
	symval  string
}

func getField(line []byte) (x, rest []byte) {
	i := bytes.IndexByte(line, ' ')
	if i == -1 {
		x = line
	} else {
		x = line[:i]
		rest = line[i+1:]
	}
	return
}

func getFieldP(line []byte) (x, p, rest []byte, err error) {
	i := bytes.IndexByte(line, ' ')
	if i == -1 {
		x = line
	} else {
		x = line[:i]
		rest = line[i+1:]
	}
	i = bytes.IndexByte(x, '(')
	if i == -1 {
		return x, p, rest, nil
	}
	x = line[:i]
	rest = line[i+1:]
	i = bytes.IndexByte(rest, ')')
	if i == -1 {
		return nil, nil, nil, errors.New("missing ')'")
	}
	p = rest[:i]
	rest = rest[i+1:]
	switch {
	case len(rest) == 0:
	case rest[0] == ' ':
		rest = rest[1:]
	default:
		return nil, nil, nil, errors.New("no space after ')'")
	}
	return x, p, rest, nil
}

func parseEntry(line []byte, curoff *uint64) (e entry, err error) {
	//fmt.Println("parse line=", line)

	var f []byte
	f, line = getField(line)
	if f[0] == 0 {
		return e, nil
	}
	if len(f) != 1 || f[0] < 'a' || 'z' < f[0] {
		return e, fmt.Errorf("invalid type: %q", f)
	}
	e.ty = rune(f[0])
	f, line = getField(line)
	mode, err := strconv.ParseUint(string(f), 8, strconv.IntSize)
	if err != nil {
		return e, fmt.Errorf("invalid mode: %q", f)
	}
	e.mode = uint(mode)
	f, line = getField(line)
	e.user = string(f)
	f, line = getField(line)
	e.group = string(f)
	f, line = getField(line)
	e.path = string(f)
	_, line = getField(line)
	fmt.Println("parseEntry:", e.mode, e.user, e.group, e.path)

	var p []byte
	for len(line) != 0 {
		f, p, line, err = getFieldP(line)
		fmt.Println("           ", string(f), string(p))

		if err != nil {
			return e, err
		}
		if p == nil {
			continue
		}
		var x uint64
		switch string(f) {
		case "sum":
			x, err = strconv.ParseUint(string(p), 10, 64)
			if err != nil {
				return e, fmt.Errorf("invalid sum: %q", p)
			}
			e.sum = x
		case "size":
			x, err = strconv.ParseUint(string(p), 10, 64)
			if err != nil {
				return e, fmt.Errorf("invalid size: %q", p)
			}
			e.size = x
		case "cmpsize":
			x, err = strconv.ParseUint(string(p), 10, 64)
			if err != nil {
				return e, fmt.Errorf("invalid cmpsize: %q", p)
			}
			e.cmpsize = x
		case "symval":
			e.symval = string(p)

		case "f", "exitop", "nohist", "nostrip", "mach", "postop", "config":
		default:
			fmt.Printf("UNKNOWN: %q\n", f)
		}
	}
	if e.ty == 'f' {
		e.offset = *curoff
		fmt.Println("            offset start=", *curoff)
		*curoff += uint64(len(e.path)) + 2 // MAGIC: was 2
		fmt.Println("                   + path+2=", e.path, *curoff)
		if e.cmpsize > 0 {
			*curoff += e.cmpsize
			fmt.Println("                   + cmpsize=", e.cmpsize, *curoff)
		} else {
			*curoff += e.size
			fmt.Println("                   + size=", e.size, *curoff)
		}
	}
	return e, nil
}

func readIDB(name string) ([]entry, error) {
	fp, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	sc := bufio.NewScanner(fp)
	var curoff uint64 = 2 // MAGIC: was 13
	var r []entry
	for lineno := 1; sc.Scan(); lineno++ {
		line := sc.Bytes()
		if line[0] != 0 {
			e, err := parseEntry(line, &curoff)
			if err != nil {
				return nil, fmt.Errorf("%s:%d: %v", name, lineno, err)
			}
			r = append(r, e)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return r, nil
}
