package main

import (
	"flag"
	"fmt"
	"github.com/mtib/rexel/reader"
	"io"
	"os"
	"strings"
	"text/template"
)

const useArgs = `Use rexel like this:
    rexel <template> <data> [-v] [-o <output>]`

func main() {
	output := flag.String("o", "stdout", "file to output to")
	verbose := flag.Bool("v", false, "be verbose")
	flag.Parse()
	tmplFile := flag.Arg(0)
	dataFile := flag.Arg(1)

	if flag.NArg() < 2 {
		fmt.Println(useArgs)
		os.Exit(1)
	}
	verbf(fmt.Sprintf("reading %s to fill %s", dataFile, tmplFile), verbose)

	rex, extractError := reader.Extract(dataFile)
	if extractError != nil {
		fmt.Printf("%s\n", extractError)
		os.Exit(2)
	}

	verbf("preparing rexel object", verbose)
	rex.Prepare((func(outf string) string {
		if outf == "stdout" {
			return reader.RTF_FORMAT
		}
		switch {
		case strings.HasSuffix(outf, ".rtf"):
			return reader.RTF_FORMAT
		case strings.HasSuffix(outf, ".html"):
			return reader.HTML_FORMAT
		}
		return reader.RTF_FORMAT
	})(tmplFile))

	var file io.WriteCloser
	switch *output {
	case "stdout":
		file = os.Stdout
	default:
		file, _ = os.OpenFile(*output, os.O_CREATE|os.O_RDWR, 0644)
		defer file.Close()
	}

	tmpl, _ := template.ParseFiles(tmplFile)
	fillErr := tmpl.Execute(file, &rex)

	if fillErr != nil {
		fmt.Errorf("was not able to fill the template")
	}
	verbf(fmt.Sprint(rex), verbose)
}

func Fill(tmplFile, dataFile, outputFile string) error {
	rex, extractError := reader.Extract(dataFile)
	if extractError != nil {
		return extractError
	}
	rex.Prepare((func(outf string) string {
		switch {
		case strings.HasSuffix(outf, ".rtf"):
			return reader.RTF_FORMAT
		case strings.HasSuffix(outf, ".html"):
			return reader.HTML_FORMAT
		}
		return reader.RTF_FORMAT
	})(tmplFile))

	var file io.WriteCloser
	switch outputFile {
	case "stdout":
		file = os.Stdout
	default:
		file, _ = os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR, 0644)
		defer file.Close()
	}

	tmpl, tmplErr := template.ParseFiles(tmplFile)

	if tmplErr != nil {
		return tmplErr
	}
	fillErr := tmpl.Execute(file, &rex)

	if fillErr != nil {
		return fillErr
	}
	return nil
}

func verbf(str string, v *bool) {
	if *v {
		fmt.Printf("%s\n", str)
	}
}

func test() {
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
