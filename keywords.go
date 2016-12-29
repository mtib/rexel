package reader

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

func keyword(r *Rexel, line int, format string) (label, display string, nextLine int) {
	//fmt.Printf("START KEYWORD FUNCTION:%s\n", r.row[line])
	if strings.TrimSpace(strings.Join(r.row[line], "")) == ""{	//check for empty []string
		return "<empty line>", "<empty>", line + 1
	}else{
		if r.row[line][0] == "" && r.row[line][1:] != nil {		// if no label is set, revert rows and find last label
			u := line
			for u >= 0 {
				u -= 1
				if r.row[u][0] != ""{
					r.row[line][0] = r.row[u][0]
					//fmt.Printf("line:%d - u:%d label:%s\n", line, u, r.row[u][0])
					break
				}
			}
		}
		//fmt.Printf("Keyword:%s\n", r.row[line][1])
	switch r.row[line][1] {
	case "text":
		//label, display = text(r.row[line][1:], format)
		label = r.row[line][0]
		display = r.row[line][2]
		//fmt.Printf("format:%s\n", format)
		//fmt.Printf("rexel:%s\n", r.row[line][0:])
		//fmt.Printf("label:%s\n", label)
		//fmt.Printf("display:%s\n", display)
		nextLine = line + 1
		return
	case "table":
		label, display, nextLine = table(r, 2, line, format)
		return
	}
	return "<empty line>", "<empty>", line + 1
	}
}

func text(param string, format string) (string) {
	switch format {
	case HTML_FORMAT:
		//return param[0], fmt.Sprintf("<span>%s</span>", strings.Join(trim(param[1:]), " ")) // for []string
		return fmt.Sprintf("<span>%s</span>", param)
	case RTF_FORMAT:
		return param
	}
	return "UNKNOWN FORMAT"
}

func table(param *Rexel, left, top int, format string) (string, string, int) {
	end := top
	textrep := "<table>"
	for _, r := range param.row[top:] {	//get all following rows (top := lines)
		//fmt.Printf("end:%d\n",end)
		//fmt.Printf("top:%d\n",top)
		//fmt.Printf("r:%s\n",r)
		//fmt.Printf("r[left-1]:%s\n",r[left-1])
		if end > top && r[left-1] != "" {	// if current row element on most left side is not zero (next label/keyword detected) or line number is greater than end from beginning (avoid start error), abort
			break
		} else {
			textrep += "<tr>"
			//for _, e := range r[left:] {
				e := r[left]
				//fmt.Printf("r[left]:%s\n",r[left])
				//fmt.Printf("e:%s\n",e)
				textrep += fmt.Sprintf("<td>%s</td>", e)
				fmt.Printf("textrep:%s\n",textrep)
			//}
			textrep += "</tr>"
		}
		end += 1
	}
	textrep += "</table>"
	fmt.Printf("textrep:%s\n",textrep)
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
		//fmt.Println(cmd)
		io.WriteString(in, textrep)
		in.Close()
		cmd.Start() // xerr := ...
		// fmt.Println(xerr)
		byterep, _ := ioutil.ReadAll(out)
		//fmt.Println(byterep)
		return "table", string(byterep), end
	}
	return "UNKNOWN", "UNKNOWN", top + 1
}
