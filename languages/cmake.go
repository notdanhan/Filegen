// cmake.go
// This is for outlining CMake generation
// This is very opinionated.
package languages

import (
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/notdanhan/filegen/common"
)

type cmakeFileFiller struct {
	ProjectName string // name of the project
	AuthorName  string // name of the author
}

const rootCmakeListsFile = `# {{.ProjectName}}
cmake_minimum_required(VERSION 4.2)
set(CMAKE_CXX_STANDARD 23)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

project({{.ProjectName}} CXX)

if(CMAKE_CXX_COMPILER_ID STREQUAL "MSVC")
	add_compile_options(/W4)
else()
	add_compile_options(-Wall -Wextra -Wpedantic)
endif()

set(CMAKE_BINARY_DIR ${CMAKE_CURRENT_LIST_DIR}/build)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)

`

func MakeCMakeProject(projectName string) {
	tRoot := template.Must(template.New("root_cmake").Parse(rootCmakeListsFile))
	// Remove illegal characters from the filenames provided
	re := regexp.MustCompile(`[\\ -.\/:*?\"<>|$?{}\[\]]`)
	projectName = re.ReplaceAllString(projectName, "")
	projectName = strings.ToLower(projectName)
	authName := common.GetAuthorName()

	myFileDetails := cmakeFileFiller{
		ProjectName: projectName,
		AuthorName:  authName,
	}

	if _, err := os.Stat("CMakeLists.txt"); err != nil {
		log.Println("Creating root cmake file")

		cmakeFile, err := os.Create("CMakeLists.txt")
		if err != nil {
			log.Fatal(err.Error())
		}
		tRoot.Execute(cmakeFile, myFileDetails)
		cmakeFile.Close()
	}
	return
}
