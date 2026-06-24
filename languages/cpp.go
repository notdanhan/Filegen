// cpp.go
// This is for outlining C++ source code generation
// Currently, it is not using cpp modules
// TODO: Add CPP module support when more widely available
package languages

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
)

const cpp_header = `/** \file
 * File Description goes here
 *
 * File:   {{.FileName}}.h
 * Author: {{.AuthorName}}
 *
 * Copyright: {{.Year}} {{.AuthorName}}
 */
#ifndef {{.FileNameCaps}}_H_{{.UUID}}
#define {{.FileNameCaps}}_H_{{.UUID}} 1
#pragma once

namespace foo
{

/**
 * \brief returns 42
inline int foo()
{
	return 42;
}

} // namespace foo

#endif // {{.FileNameCaps}}_H_{{.UUID}}

`

const cpp_body = `/** \file
 * File:      {{.FileName}}.cpp
 * Author:    {{.AuthorName}}
 *
 * Copyright: {{.Year}} {{.AuthorName}}
 */
#include "{{.FileName}}.h"

`

type cppFileFiller struct {
	FileName     string // name of the file in regular case
	FileNameCaps string //name of the file in Block capitals, for Header Guards
	AuthorName   string // name of the author who wrote this
	UUID         string // UUID for the header guard
	Year         int    // Year of authorship
}

// Generate the CPP file filler so I can populate the file and stuff
func createCppFileFiller(fName string) cppFileFiller {
	// Remove illegal characters from the filenames provided
	re := regexp.MustCompile(`[\\ -.\/:*?\"<>|$?{}\[\]]`)
	fName = re.ReplaceAllString(fName, "")
	fName = strings.ToLower(fName)
	fNameUpper := strings.ToUpper(fName)

	// Get the git username for attribution purposes
	command := exec.Command("git", "config", "user.name")
	authNameTmp, err := command.Output()
	// shit happens when you party naked ~Socrates
	authName := "Thurman Merman"
	if err == nil {
		authName = string(authNameTmp)
	}
	authName = strings.ReplaceAll(authName, "\n", "")

	year := time.Now().Year()

	tUUID := uuid.New()

	tUUIDstr := strings.ToUpper(strings.ReplaceAll(tUUID.String(), "-", ""))

	return cppFileFiller{AuthorName: authName,
		FileName:     fName,
		FileNameCaps: fNameUpper,
		Year:         year,
		UUID:         tUUIDstr,
	}
}

// Generate CPP file pairs
func MakeCppFilePairs(filenames []string) {
	tHeader := template.Must(template.New("cpp_header").Parse(cpp_header))
	tBody := template.Must(template.New("cpp_body").Parse(cpp_body))
	for _, filename := range filenames {
		// Create the file data
		mFile := createCppFileFiller(filename)

		// Check if the header file exists, don't overwrite it if it does.
		if _, err := os.Stat(mFile.FileName + ".h"); err != nil {
			log.Println("Creating", mFile.FileName+".h")

			headerFile, err := os.Create(mFile.FileName + ".h")
			if err != nil {
				log.Fatal(err.Error())
			}
			tHeader.Execute(headerFile, mFile)
			headerFile.Close()
		}

		// Do the same for the source file
		if _, err := os.Stat(mFile.FileName + ".cpp"); err != nil {
			log.Println("Creating", mFile.FileName+".cpp")
			cppFile, err := os.Create(mFile.FileName + ".cpp")
			if err != nil {
				log.Fatal(err.Error())
			}
			tBody.Execute(cppFile, mFile)
			cppFile.Close()
		}
	}
}
