package main

import (
	"github.com/google/uuid"

	"os"
	"os/exec"

	"time"

	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type FileFiller struct {
	FileName     string
	FileNameCaps string
	AuthorName   string
	UUID         string
	Year         int
}

func CreateFileFiller(fName string) FileFiller {

	// Remove illegal characters from the filenames
	re := regexp.MustCompile(`[\\ -.\/:*?\"<>|$?{}\[\]]`)
	fName = re.ReplaceAllString(fName, "")
	fName = strings.ToLower(fName)
	fNameUpper := strings.ToUpper(fName)

	// get the git username for generating the name
	command := exec.Command("git", "config", "user.name")
	authNameTmp, err := command.Output()
	authName := ""
	if err == nil {
		authName = string(authNameTmp)
	}
	authName = strings.ReplaceAll(authName, "\n", "")

	year := time.Now().Year()

	tUUID := uuid.New()

	tUUIDstr := strings.ToUpper(strings.ReplaceAll(tUUID.String(), "-", ""))

	return FileFiller{AuthorName: authName,
		FileName:     fName,
		FileNameCaps: fNameUpper,
		Year:         year,
		UUID:         tUUIDstr,
	}
}

func main() {
	const header = `/*
 * File:      {{.FileName}}.h
 * Author:    {{.AuthorName}}
 *
 * Copyright: {{.Year}} {{.AuthorName}}
 */
#ifndef {{.FileNameCaps}}_H_{{.UUID}}
#define {{.FileNameCaps}}_H_{{.UUID}} 1
#pragma once

#endif // {{.FileNameCaps}}_H_{{.UUID}}

`

	const cpp = `/*
 * File:      {{.FileName}}.cpp
 * Author:    {{.AuthorName}}
 *
 * Copyright: {{.Year}} {{.AuthorName}}
 */
#include "{{.FileName}}.h"

`
	args := os.Args[1:]

	tHeader := template.Must(template.New("header").Parse(header))
	tBody := template.Must(template.New("cpp").Parse(cpp))

	for _, arg := range args {
		mFile := CreateFileFiller(arg)
		// Check if header file exists, don't overwrite if it does
		if _, err := os.Stat(mFile.FileName + ".h"); err != nil {
			fmt.Println("Creating", mFile.FileName+".h")
			headerFile, err := os.Create(mFile.FileName + ".h")
			if err != nil {
				fmt.Println("ERROR:", err.Error())
				return
			}
			tHeader.Execute(headerFile, mFile)
			headerFile.Close()
		}

		// Check if source file exists, don't overwrite if it does
		if _, err := os.Stat(mFile.FileName + ".cpp"); err != nil {
			fmt.Println("Creating", mFile.FileName+".cpp")
			cppFile, err := os.Create(mFile.FileName + ".cpp")
			if err != nil {
				fmt.Println("ERROR:", err.Error())
				return
			}
			tBody.Execute(cppFile, mFile)
			cppFile.Close()
		}
	}
}
