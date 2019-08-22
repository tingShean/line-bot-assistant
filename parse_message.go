package main

import (
	"strings"
	"log"
	"bytes"
)

func ParseMessage(msg string) (string, error) {
	var buf bytes.Buffer
	log.Println(msg)

	if strings.Compare("help", msg) == 0 {
		res := "=======Help======\nhelp	使用說明"
		return res, nil
	}

	if strings.Contains(msg, "taskadd") {
		log.Println("task add")
		taskFunc(msg)

		return "Ok", nil
	}

	if strings.Compare("task", msg) == 0 {
		log.Println("task list")
		buf.WriteString("Your Task List >>>>\n")
		res := get_google_sheet(msg)
		buf.WriteString(res)
		return buf.String(), nil
	}

	return msg, nil
}

func taskFunc(msg string) {
	add_google_sheet_row("task", msg)
}
