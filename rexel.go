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
	rex.Prepare()
	fmt.Println(rex)
	err = FillTmpl(&rex, "example/filled.rtf")
	if err != nil {
		panic(err)
	}
}
