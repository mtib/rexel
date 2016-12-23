package main

import (
	"fmt"
	"github.com/mtib/rexel/reader"
)

func main() {
	rex, err := reader.Extract("example/data.csv")
	if err != nil {
		panic(err)
	}
	fmt.Println(rex)
	rex.Prepare(reader.RTF_FORMAT)
	fmt.Println(rex)
	err = FillTmpl(&rex, "example/filled.rtf")
	if err != nil {
		panic(err)
	}
	rex.Prepare(reader.HTML_FORMAT)
	fmt.Println(rex)
	err = FillTmpl(&rex, "example/filled.html")
	if err != nil {
		panic(err)
	}
}
