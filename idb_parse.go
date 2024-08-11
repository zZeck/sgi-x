package main

import (
	"strconv"
	"strings"
)

//Could try using a csv parsing library
//that allows you to set the delimiter to ' '
//The go standard one will not work though
//Because of how quotes appear within a field.
//LazyQuotes does not do what is needed

// won't handle escaped quotes correctly within a quoted string
func idb_line_fields(line string) []string {
	quoted := false
	a := strings.FieldsFunc(line, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})
	return a
}

// won't handle escaped quotes correctly within a quoted string
func idb_field_key_value(field string) (string, string) {
	paren_index := strings.IndexByte(field, '(')
	if paren_index == -1 {
		return "", field
	}
	key := field[:paren_index]

	quoted := false
	a := strings.FieldsFunc(field, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == '(' || r == ')'
	})

	return key, strings.Trim(a[1], "\"")
}

func idb_line_entry(line string, entry_offset int) entry2 {
	blah := idb_line_fields(line)

	m := make(map[string]string)
	for _, field := range blah {
		key, value := idb_field_key_value(field)
		if key != "" {
			m[key] = value
		}
	}

	size_in_archive := 0
	compressed := false
	if val, ok := m["cmpsize"]; ok {
		num, _ := strconv.ParseInt(val, 10, 0)
		if num > 0 {
			size_in_archive = int(num)
			compressed = true
		}
	}

	final_size := 0
	if val, ok := m["size"]; ok {
		num, _ := strconv.ParseInt(val, 10, 0)
		//cmpsize can be present, but set to 0, size should be used
		if size_in_archive == 0 {
			size_in_archive = int(num)
		}
		final_size = int(num)
	}

	symval := ""
	if val, ok := m["symval"]; ok {
		symval = val
	}

	data_offset := entry_offset
	//directory entries will have no size
	//do not want to accidentally incrment offset by their path
	if final_size > 0 {
		data_offset = entry_offset + 2 + len(blah[4]) //end of last file, + 2 bytes of unknown, + path
	}

	return entry2{
		idb_entry_type:  blah[0],
		path:            blah[4],
		size_in_archive: size_in_archive,
		final_size:      final_size,
		symval:          symval,
		compressed:      compressed,
		data_offset:     data_offset,
	}
}

type entry2 struct {
	idb_entry_type  string
	path            string
	size_in_archive int
	final_size      int
	symval          string
	compressed      bool
	data_offset     int
}
