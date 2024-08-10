package main

import (
	"testing"
    "reflect"
)

func TestLineToFields(t *testing.T) {
	name := "f 0444 root sys path/to/some.o path/to/some.o sum(64) size(64) postop(\"touch randomfile\") nostrip mach(MAYBEENV=DATA) f(1234) cmpsize(0) some.dev.base"
	blah := idb_line_fields(name)
    expect := []string {
        "f",
        "0444",
        "root",
        "sys",
        "path/to/some.o",
        "path/to/some.o",
        "sum(64)",
        "size(64)",
        "postop(\"touch randomfile\")",
        "nostrip",
        "mach(MAYBEENV=DATA)",
        "f(1234)",
        "cmpsize(0)",
        "some.dev.base",
    }
	if !reflect.DeepEqual(blah, expect) {
		t.Fatalf("err")
	}
}

func TestFieldToKeyValue(t *testing.T) {
	name := "sum(64)"
	key, value := idb_field_key_value(name)
	if key != "sum" {
		t.Fatalf("err")
	}
    if value != "64" {
		t.Fatalf("err")
	}
}

func TestFieldToKeyQuotedValue(t *testing.T) {
	name := "postop(\"touch randomfile\")"
	key, value := idb_field_key_value(name)
	if key != "postop" {
		t.Fatalf("err")
	}
    if value != "touch randomfile" {
		t.Fatalf("err")
	}
}

func TestFieldToKeyQuotedValueOnly(t *testing.T) {
	name := "nostrip"
	key, value := idb_field_key_value(name)
	if key != "" {
		t.Fatalf("err")
	}
    if value != "nostrip" {
		t.Fatalf("err")
	}
}

func TestLineToEntry(t *testing.T) {
	name := "f 0444 root sys path/to/some.o path/to/some.o sum(64) size(100) postop(\"touch randomfile\") nostrip symval(asympath) mach(MAYBEENV=DATA) f(1234) cmpsize(55) some.dev.base"
	blah := idb_line_entry(name, 0)
    expect := entry2 {
        idb_entry_type: "f",
        path: "path/to/some.o",
        size_in_archive: 55,
		final_size: 100,
        symval: "asympath",
        compressed: true,
		offset: 55,
    }
	if !reflect.DeepEqual(blah, expect) {
		t.Fatalf("err")
	}
}

func TestLineToEntry2(t *testing.T) {
	name := "f 0444 root sys path/to/some.o path/to/some.o sum(64) size(100) postop(\"touch randomfile\") nostrip mach(MAYBEENV=DATA) f(1234) some.dev.base"
	blah := idb_line_entry(name, 0)
    expect := entry2 {
        idb_entry_type: "f",
        path: "path/to/some.o",
        size_in_archive: 100,
		final_size: 100,
        symval: "",
        compressed: false,
		offset: 100,
    }
	if !reflect.DeepEqual(blah, expect) {
		t.Fatalf("err")
	}
}

func TestLineToEntryUseSizeWhenCmpSize0(t *testing.T) {
	name := "f 0444 root sys path/to/some.o path/to/some.o sum(64) size(100) postop(\"touch randomfile\") nostrip symval(asympath) mach(MAYBEENV=DATA) f(1234) cmpsize(0) some.dev.base"
	blah := idb_line_entry(name, 0)
    expect := entry2 {
        idb_entry_type: "f",
        path: "path/to/some.o",
        size_in_archive: 100,
		final_size: 100,
        symval: "asympath",
        compressed: false,
		offset: 100,
    }
	if !reflect.DeepEqual(blah, expect) {
		t.Fatalf("err")
	}
}

func TestLineToEntryAddsPassedOffsetToSize(t *testing.T) {
	name := "f 0444 root sys path/to/some.o path/to/some.o sum(64) size(100) postop(\"touch randomfile\") nostrip symval(asympath) mach(MAYBEENV=DATA) f(1234) cmpsize(0) some.dev.base"
	blah := idb_line_entry(name, 13)
    expect := entry2 {
        idb_entry_type: "f",
        path: "path/to/some.o",
        size_in_archive: 100,
		final_size: 100,
        symval: "asympath",
        compressed: false,
		offset: 113,
    }
	if !reflect.DeepEqual(blah, expect) {
		t.Fatalf("err")
	}
}