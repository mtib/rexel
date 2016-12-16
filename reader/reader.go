package reader

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Rexel struct {
	row      [][]string
	labels   map[string]string
	prepared bool
	origin   string
}

func Extract(file string) (Rexel, error) {
	byteData, err := ioutil.ReadFile(file)
	if err != nil {
		return Rexel{prepared: false}, err
	}
	stringData := string(byteData)
	rows := strings.Split(stringData, "\n")
	numRows := len(rows)

	rex := Rexel{row: make([][]string, numRows), labels: nil, prepared: false, origin: file}
	for k, _ := range rows {
		rex.row[k] = strings.Split(rows[k], ",")
	}
	return rex, nil
}

func (r Rexel) String() string {
	rString := ""
	for _, row := range r.row {
		for _, entry := range row {
			if entry == "" {
				rString += "NONE"
			} else {
				rString += entry
			}
			rString += ", "
		}
		rString += "\n"
	}

	rString += "\n"

	if r.prepared {
		rString += fmt.Sprintln("Rexel file is prepared:")
		for k, v := range r.labels {
			rString += fmt.Sprintf("Label('%s') := '%s'\n", k, v)
		}
	} else {
		rString += fmt.Sprintln("Rexel file is not prepared")
	}
	return rString
}

func (r *Rexel) Prepare() (err error) {
	r.labels = make(map[string]string, len(r.row))
	// TODO support more keywords
	k := 0
	for k < len(r.row) {
		var label string
		var display string
		label, display, k = keyword(r, k)
		r.labels[label] = display
		// fmt.Printf("found label '%s' with content: '%s'\n", label, display)
	}
	r.prepared = true
	return err
}

func (r *Rexel) Label(label string) string {
	return r.labels[label]
}

func trim(row []string) []string {
	if len(row) == 0 {
		return row
	}
	if row[0] == "" {
		return trim(row[1:])
	}
	if row[len(row)-1] == "" {
		return trim(row[:len(row)-1])
	}
	return row
}
