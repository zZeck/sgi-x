package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func mainE() error {
	args := os.Args
	if len(args) < 2 || len(args) > 4 {
		return errors.New("usage: sgix <file.idb> [<data> [<dir>]]")
	}
	idbfile := args[1]
	datafile := args[2]
	dest := args[3]

	idb_file, _ := os.Open(idbfile)
	defer idb_file.Close()

	sc := bufio.NewScanner(idb_file)
	lines := make([]string, 0)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	entries := make([]entry2, 0)
	//im001V530P00 + nul term
	entry_offset := 13
	for _, line := range lines {
		entry := idb_line_entry(line, entry_offset)
		entries = append(entries, entry)
		entry_offset = entry.data_offset + entry.size_in_archive
	}

	data_file, _ := os.Open(datafile)
	defer idb_file.Close()
	for _, entry := range entries {
		output_entry(entry, data_file, dest)
	}

	return nil
}

func main() {
	if err := mainE(); err != nil {
		os.Exit(1)
	}
}

func output_entry(entry entry2, src *os.File, out_dir string) error {
	name := path.Clean(entry.path)
	dest := path.Join(out_dir, name)

	src.Seek(int64(entry.data_offset), io.SeekStart)

	switch entry.idb_entry_type {
	case "f":
		os.MkdirAll(filepath.Dir(dest), 0770)
		fp, err := os.Create(dest)
		if err != nil {
			return err
		}

		if entry.compressed {
			fmt.Println("uncompress ", entry.path)
			exe := exec.Command("uncompress")
			exe.Stdin = &io.LimitedReader{R: src, N: int64(entry.size_in_archive)}
			exe.Stdout = fp
			exe.Stderr = os.Stderr
			return exe.Run()
		}

		_, err = io.CopyN(fp, src, int64(entry.final_size))
		return err
	case "d":
		if dest == "" {
			return nil
		}
		return os.Mkdir(dest, 0777)
	case "l":
		//symlink
	default:
	}

	return fmt.Errorf("unknown type")
}
