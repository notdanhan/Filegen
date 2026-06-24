package main

import (
	"os"

	"github.com/danielh2942/filegen/languages"
)

func main() {
	args := os.Args[1:]

	languages.MakeCppFilePairs(args)
}
