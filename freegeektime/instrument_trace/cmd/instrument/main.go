package main

import (
	"flag"
	"fmt"
	"github.com/hgaonice/instrument_trace/instrumenter/ast"
	"os"
	"path/filepath"
)

var (
	wrote bool
)

func init() {
	flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")
}

func usage() {
	fmt.Println("instrument [-w] xxx.go")
	flag.PrintDefaults()
}

func main() {
	fmt.Println(os.Args)
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
		return
	}

	var file string
	if len(os.Args) == 3 {
		file = os.Args[2]
	}
	if filepath.Ext(file) != ".go" {
		usage()
		return
	}
	ins := ast.New("github.com/hgaonice/instrument_trace", "trance", "Trance")
	newSrc, err := ins.Instrument(file)
	if err != nil {
		panic(err)
	}
	if newSrc == nil {
		fmt.Printf("no trance added for %s\n", file)
		return
	}

	if !wrote {
		fmt.Println(string(newSrc))
		return
	}
	//if err = ioutil.WriteFile(file, newSrc, 0666); err != nil {
	if err = os.WriteFile(file, newSrc, 0666); err != nil {
		fmt.Printf("write %s error:%v\v", file, err)
		return
	}
	fmt.Printf("instrument trance %s ok \n", file)

}
