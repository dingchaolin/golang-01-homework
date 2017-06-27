package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Student struct for student info
type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var students = make(map[string]Student)

func main() {
	var line, cmd, name string
	var id int
	f := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _ = f.ReadString('\n')
		line = strings.Replace(line, "\n", "", 1)
		_, err := fmt.Sscan(line, &cmd, &name, &id)
		// EOF for no enough space-separated values, like `list` cmd
		if err != nil && err != io.EOF {
			fmt.Println(" + parse cmd or info err:", err)
			continue
		}
		if cmd == "exit" {
			break
		}
		switch cmd {
		case "add":
			addStu(name, id)
		case "list":
			listStu()
		case "load":
			loadStu(name)
		case "save":
			saveStu(name)
		case "":
			continue
		default:
			usage()
		}
		line, cmd = "", ""
	}
}

func saveStu(name string) {
	var err error
	var fd *os.File

	// file existed and override
	if checkFileExist(name) {
		if !checkOverride(name) {
			fmt.Println(" + cancel save to", name)
			return
		}
		if fd, err = os.OpenFile(name, os.O_RDWR|os.O_TRUNC, 0644); err != nil {
			log.Fatal("open file error of %s", name)
		}
	}
	// save to new file
	if fd, err = os.Create(name); err != nil {
		log.Fatal("open new file error of %s", name)
	}
	if buf, err := json.Marshal(students); err != nil {
		log.Fatal("marshal stu info error")
	} else if _, err := fd.Write(buf); err == nil {
		fmt.Println(" + save success")
	} else {
		log.Fatal("save error")
	}
}

// check file f exist or not
func checkFileExist(f string) bool {
	if _, err := os.Stat(f); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

// check override the give file name
func checkOverride(f string) bool {
	var s string
	for {
		fmt.Printf(" + override file: %s, y or n? :", f)
		fmt.Scanf("%s", &s)
		if string(s[0]) == "y" || string(s[0]) == "Y" {
			return true
		} else if string(s[0]) == "n" || string(s[0]) == "N" {
			return false
		}
	}
}

func loadStu(name string) {
	fmt.Println("params:", name)
}

func addStu(name string, id int) {
	if _, ok := students[name]; ok {
		fmt.Println(" + duplicated name:", name)
		return
	}
	students[name] = Student{ID: id, Name: name}
}

func listStu() {
	if len(students) == 0 {
		fmt.Println(" + no student info here")
		return
	}
	fmt.Println(" + Id\tName:")
	for _, val := range students {
		fmt.Printf(" + %d\t%s\n", val.ID, val.Name)
	}
}

func usage() {
	fmt.Println(" + cli usage:")
	fmt.Println(" + add name id -- add student info")
	fmt.Println(" + list \t-- list student info")
	fmt.Println(" + load file \t-- load student from file")
	fmt.Println(" + save file \t-- save student info file")
	fmt.Println(" + exit \t-- exit the cli")
}
