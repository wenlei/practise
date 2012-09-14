package main

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	ident = "    "
	sub   = "⤷ "
	//sub  = "↳"
)

type entry struct {
	mod     string
	files   int
	creator string
	role    string
	size    int
	date    string
	name    string
}

func (e entry) String() string {
	return fmt.Sprintln(e.mod, e.files, e.creator, e.role, e.size, e.date, e.name)
}

var (
	self_name    = os.Args[0]
	current_dict = "."
)

func main() {
//-----validate the input variables // may hit err when the file name is like ^.*/$ ie, foo/
	if len(os.Args) > 1 {
		current_dict = os.Args[1]
//		if string(current_dict[len(current_dict)-1]) == "/" {
//			current_dict = current_dict[:len(current_dict)-1]
//		}
	}
//------
	lines := show_dict(".")
	if current_dict != "" {
		lines = show_dict("./" + current_dict)
	}
	lst := fmt_to_list(lines)
	fmt.Println(current_dict)
	printEntry(current_dict, lst)
}

func printEntry(uplevel_dict string, el *list.List) {
	if el.Len() == 0 {
		return
	}
	ident_number := strings.Count(uplevel_dict, "/")
	preIdent := ""
	for i := 0; i <= ident_number; i++ {
		preIdent += ident
	}

	for e := el.Front(); e != nil; e = e.Next() {

		elem, isDict := infof(e.Value.(entry))

		if isDict {
			current_dict = uplevel_dict + "/" + elem
			fmt.Printf("%v%v<%v>\n", preIdent, sub, elem)
			printEntry(current_dict, fmt_to_list(show_dict(current_dict))) // the dictionary name error happens in ln100
		} else {
			fmt.Printf("%v%v%v\n", preIdent, sub, elem)
		}
	}
}

func infof(e entry) (name string, is_dict bool) {
	return e.name, is_dictionary(e)
}

func is_dictionary(e entry) bool {
	return string((e.mod)[0]) == "d"
}

func fmt_to_list(lines []string) *list.List {
	lst := list.New()
	for _, v := range lines { //MENTION HERE
		fields := strings.Fields(v)
		//coping with precedent spacing before file name like ^ filename$
		re, _ := regexp.Compile("^[\\S]+\\s+\\w+\\s+\\w+\\s+\\w+\\s+\\w+\\s+\\w+\\s+\\w+\\s+\\S+\\s")
		r := re.ReplaceAllLiteralString(v, "")

		last := len(fields) - 1
		if (fields[last] == "." && len(fields[last]) == 1) || (fields[last] == ".." && len(fields[last]) == 2) {
			continue
		}
		lst.PushBack(stuff(fields, r))
	}
	return lst
}

func stuff(fields []string, name string) entry { //blank character problem
	size, _ := strconv.Atoi(fields[4])
	files_number, _ := strconv.Atoi(fields[1])
	return entry{fields[0], files_number, fields[2], fields[3], size, strings.Join(fields[5:8], "-"), name}
}

func show_dict(d string) []string {
	//	fmt.Println("ls -l ", d)
	cmd := exec.Command("ls", "-la", d)
	buf, err := cmd.Output()
	if err != nil {
		fmt.Println("Error: ", err.Error(), "in line 100 func show_dict")
		log.Fatal(self_name, ": ", err.Error())
	}
	lines := strings.Split(string(buf), "\n")
	return lines[1 : len(lines)-1] //ok
}