package main

import (
	"flag"
	"github.com/notdanhan/filegen/languages"
)

func main() {
	language := flag.String("lang", "cpp", "language to use (Current available are cpp and cmake")
	flag.Parse()
	if *language == "cpp" {
		languages.MakeCppFilePairs(flag.Args())
	} else if *language == "cmake" {
		languages.MakeCMakeProject(flag.Arg(0))
	}
}
