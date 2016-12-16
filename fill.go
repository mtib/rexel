package main

import (
	"fmt"
	"github.com/mtib/rexel/reader"
	"os"
	"text/template"
)

var (
	tmpl, _ = template.ParseFiles("example/tmpl.rtf")
)

func FillTmpl(r *reader.Rexel, fileout string) error {
	file, err := os.OpenFile(fileout, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	tmpl.Execute(file, r)
	return nil
}

func Automate(inputFile string) error {
	r, err := reader.Extract(inputFile)
	if err != nil {
		return err
	}
	r.Prepare()
	return FillTmpl(&r, fmt.Sprintf("%s.rtf", inputFile))
}
