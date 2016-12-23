package reader

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

func keyword(r *Rexel, line int, format string) (label, display string, nextLine int) {
	switch r.row[line][0] {
	case "text":
		label, display = text(r.row[line][1:], format)
		nextLine = line + 1
		return
	case "table":
		label, display, nextLine = table(r, 1, line, format)
		return
	}
	return "<empty line>", "<empty>", line + 1
}

func text(param []string, format string) (string, string) {
	switch format {
	case HTML_FORMAT:
		return param[0], fmt.Sprintf("<span>%s</span>", strings.Join(trim(param[1:]), " "))
	case RTF_FORMAT:
		return param[0], strings.Join(trim(param[1:]), " ")
	}
	return param[0], "UNKNOWN FORMAT"
}

func table(param *Rexel, left, top int, format string) (string, string, int) {
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
	textrep += "</table>"
	switch format {
	case HTML_FORMAT:
		return "table", textrep, end
	case RTF_FORMAT:
		// textrep ist HTML rep von der Tabelle
		//
		// todo
		cmd := exec.Command("pandoc", "-f", "html", "-t", "rtf")
		in, _ := cmd.StdinPipe()
		out, _ := cmd.StdoutPipe()
		fmt.Println(cmd)
		io.WriteString(in, textrep)
		in.Close()
		xerr := cmd.Start()
		fmt.Println(xerr)
		byterep, _ := ioutil.ReadAll(out)
		fmt.Println(byterep)
		return "table", string(byterep), end
	}
	return "UNKNOWN", "UNKNOWN", top + 1
}
