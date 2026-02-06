package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/timbaileyjones/cursorchats/internal/db"
	"golang.org/x/term"
)

func main() {
	chatsDir := flag.String("chats-dir", "", "Directory containing Cursor chat *.db files (default: $HOME/.cursor/chats)")
	outputPath := flag.String("output", "", "Write output to file (default: stdout)")
	colorFlag := flag.Bool("color", false, "Colorize output (roles and extracted text context)")
	flag.Parse()

	if *chatsDir == "" {
		*chatsDir = filepath.Join(os.Getenv("HOME"), ".cursor", "chats")
	}

	out := os.Stdout
	if *outputPath != "" {
		f, err := os.Create(*outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open output: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	opts := &db.DumpOptions{Color: *colorFlag}
	// Use terminal width for byte dump when stdout is a TTY (leave room for "  │ " + "0000 " offset prefix)
	if *outputPath == "" {
		if fd := int(os.Stdout.Fd()); term.IsTerminal(fd) {
			if w, _, err := term.GetSize(fd); err == nil && w > 9 {
				opts.ByteDumpWidth = (w - 9) & ^1 // 9 = offset prefix; even so hex pairs never split
			}
		}
	}
	if err := db.DumpAll(*chatsDir, out, opts); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
