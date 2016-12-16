package reader

import (
	"fmt"
	"strings"
)

func keyword(r *Rexel, line int) (label, display string, nextLine int) {
	switch r.row[line][0] {
	case "text":
		label, display = text(r.row[line][1:])
		nextLine = line + 1
		return
	case "table":
		label, display, nextLine = table(r, 1, line)
		return
	}
	return "<empty line>", "<empty>", line + 1
}

func text(param []string) (string, string) {
	return param[0], strings.Join(trim(param[1:]), " ")
}

func table(param *Rexel, left, top int) (string, string, int) {
	end := top
	textrep := "<table>"
	for _, r := range param.row[top:] {
		if end > top && r[left-1] != "" {
			break
		} else {
			textrep += "<tr>"
			for _, e := range r[left:] {
				textrep += fmt.Sprintf("<td>%s</td>", e)
			}
			textrep += "</tr>"
		}
		end += 1
	}
	return "table", textrep, end
}
