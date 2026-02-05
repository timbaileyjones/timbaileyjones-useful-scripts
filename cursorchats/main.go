package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/timbaileyjones/cursorchats/internal/db"
)

func main() {
	chatsDir := flag.String("chats-dir", "", "Directory containing Cursor chat *.db files (default: $HOME/.cursor/chats)")
	outputPath := flag.String("output", "", "Write output to file (default: stdout)")
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

	if err := db.DumpAll(*chatsDir, out); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
