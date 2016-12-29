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

const (
	RTF_FORMAT  = "rtf"
	HTML_FORMAT = "html"
)

// Extract csv file
// Input: file <- csv file
// Output: Rexel <- object with informations about row, labels, prepared, origin struct as defined above
//         error <- error flag
func Extract(file string) (Rexel, error) {
	byteData, err := ioutil.ReadFile(file)         // Pure text with no format
	if err != nil {                                // if an error occured, set in rexel object -> prepared to false
		return Rexel{prepared: false}, err
	}
	stringData := string(byteData)                 // convert text from csv file to strings
	rows       := strings.Split(stringData, "\n")  // split string variable into rows
	numRows    := len(rows)                        // get number of rows; needed for number of replacements in template
    //fmt.Printf("stringData\n%s\n", stringData)
    //fmt.Printf("rows\n%s\n", rows)
	//fmt.Printf("numRows %d\n", numRows)
	rex        := Rexel{row: make([][]string, numRows), labels: nil, prepared: false, origin: file}   // define variable rex as a rexel object with following settings: row is set to a new created object of type string and size numRows, no labels, not prepared and origin is set to csv file
	for k, _   := range rows {                     // loop through all rows in variable row and write those lines in rex object
		rex.row[k] = strings.Split(rows[k], ",")
        //fmt.Printf("Extract rows[%d]: %s\n", k, rows[k])
        //fmt.Printf("\n")
	}
	return rex, nil
}

func (r Rexel) String() string {
	rString := ""
	for _, row := range r.row {
        sl := len(row)
        //fmt.Printf("row len:%d\n", sl)
		for j, entry := range row {
            //fmt.Printf("row (%d) - entry number (%d): %s\n", u, j, entry)
			if entry == "" {
				rString += "NONE"
			} else {
				rString += entry
			}
            if j < sl - 1 {
			rString += ", "
            }
            //fmt.Printf("rString:%s\n", rString)
		}
		rString += "\n"
	}

	rString += "\n"

	if r.prepared {
		rString += fmt.Sprintln("Rexel file is prepared:")
		for k, v := range r.labels {
			rString += fmt.Sprintf("Label('%s') := %s\n", k, v)
		}
	} else {
		rString += fmt.Sprintln("Rexel file is not prepared")
	}
	return rString
}

func (r *Rexel) Prepare(format string) (err error) {
	r.labels = make(map[string]string, len(r.row))
	// TODO support more keywords
    //fmt.Printf("row length:%d\n", len(r.row))
	k := 0
	for k < len(r.row) {
		var label string
		var display string
        //fmt.Printf("rows:%s\n", r.row[k])
		label, display, k = keyword(r, k, format) //search for keywords. Input variables are rexel object, current row and reader.RTF_FORMAT
        if r.labels[label] != ""{
            r.labels[label] += "\n" + display
        }else{
        r.labels[label] = display
    }
		//fmt.Printf("found label '%s' with content: %s\n", label, display)
        //fmt.Printf("Display für Label 'Anmerkungen': %s\n", r.labels["Anmerkungen"])
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
