// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var cbp bool = false
var cmake bool = false

var initProjectCmd = &cobra.Command{
	Use:   "init [Problem ID] [Language]",
	Short: "Create a project (statement, assets and source file for your chosen language)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { initProject(args[0], args[1]) }
		if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
			logErr(err)
		}
	},
}

var getRandPbCmd = &cobra.Command{
	Use:   "random",
	Short: "Get random problem to solve.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		getRandPb()
	},
}

func init() {
	rootCmd.AddCommand(initProjectCmd)
	rootCmd.AddCommand(getRandPbCmd)

	initProjectCmd.Flags().BoolVarP(&cbp, "cbp_project", "b", false, "Create a codeblocks project.")
	initProjectCmd.Flags().BoolVarP(&cmake, "cmake_project", "c", false, "Create a cmake project.")
}

func CopyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func moveFiles(rootDir string) error {
	return filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext == ".md" || ext == ".pdf" || ext == ".h" {
			destPath := filepath.Join(rootDir, d.Name())

			if path == destPath {
				return nil
			}

			//fmt.Println("Copying:", path, "→", destPath)
			return CopyFile(path, destPath)
		}
		return nil
	})
}

func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		filePath := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func createCBPProject(name string) {
	xmlCBP := fmt.Sprintf(XMLCBPStruct, name)

	File, err := os.Create(fmt.Sprintf("%s.cbp", name))
	if err != nil {
		logErr(fmt.Errorf("error creating .cbp file: %v", err))
	}
	defer File.Close()

	File.WriteString(xmlCBP)
}

func createCMAKEProject(name string) {
	cmakeTXT := fmt.Sprintf(CMAKEStruct, name, name)

	File, err := os.Create("CMakeLists.txt")
	if err != nil {
		logErr(fmt.Errorf("error creating \"CMakeLists.txt\": %v", err))
		return
	}
	defer File.Close()

	File.WriteString(cmakeTXT)
}

// 0-C 1-CPP 2-GO 3-Kotlin 4-JS 5-PAS 6-PHP 7-Py 8-Rust
func createSource(cwd, lang string) {
	filename := filepath.Join(cwd, "Source.")
	switch lang {
	case "c":
		filename += "c"
		os.WriteFile(filename, []byte(helloWorldprog[0]), 0644)
	case "golang":
		filename += "go"
		os.WriteFile(filename, []byte(helloWorldprog[2]), 0644)
	case "kotlin":
		filename += "kt"
		os.WriteFile(filename, []byte(helloWorldprog[3]), 0644)
	case "nodejs":
		filename += "js"
		os.WriteFile(filename, []byte(helloWorldprog[4]), 0644)
	case "pascal":
		filename += "pas"
		os.WriteFile(filename, []byte(helloWorldprog[5]), 0644)
	case "php":
		filename += "php"
		os.WriteFile(filename, []byte(helloWorldprog[6]), 0644)
	case "python3":
		filename += "py"
		os.WriteFile(filename, []byte(helloWorldprog[7]), 0644)
	case "rust":
		filename += "rs"
		os.WriteFile(filename, []byte(helloWorldprog[8]), 0644)
	default:
		filename += "cpp"
		os.WriteFile(filename, []byte(helloWorldprog[1]), 0644)
	}
}

func getName(problem_id string) string {
	url := fmt.Sprintf(URL_PROBLEM, problem_id)
	body, err := makeRequest("GET", url, nil, "0")
	if err != nil {
		logErr(err)
	}

	var info info
	if err := json.Unmarshal(body, &info); err != nil {
		logErr(err)
	}
	return info.Data.Name
}

func extractFunctionDeclarations(hFileContent string) []string {
	re := regexp.MustCompile(`\s*(int|float|char|double)\s+(\w+)\s*\((.*)\);`)
	matches := re.FindAllStringSubmatch(hFileContent, -1)

	var declarations []string
	for _, match := range matches {
		if len(match) > 3 {
			declaration := fmt.Sprintf("%s %s(%s)", match[1], match[2], match[3])
			declarations = append(declarations, declaration)
		}
	}

	return declarations
}

func initProject(problem_id, lang string) {
	cwd, err := os.Getwd()
	if err != nil {
		logErr(fmt.Errorf("could not get current working directory! error: %v", err))
	}

	newFolder := filepath.Join(cwd, fmt.Sprintf("Problem_%s_Proj_%s", problem_id, lang))
	err = os.MkdirAll(newFolder, os.ModePerm)
	if err != nil {
		logErr(fmt.Errorf("could not create project directory! error: %v", err))
	}

	err = os.Chdir(newFolder)
	if err != nil {
		logErr(fmt.Errorf("could not change directory to project dir! error: %v", err))
	}

	cwd, err = os.Getwd()
	if err != nil {
		logErr(fmt.Errorf("could not get current working directory! error: %v", err))
	}

	ok := false

	supportedLangs := checklangs(problem_id, 2)
	for i := range supportedLangs {
		if lang == supportedLangs[i] {
			ok = true
			createSource(cwd, lang)
		}
	}

	if !ok {
		os.Chdir("..")
		os.Remove(newFolder)
		logErr(errors.New("problem is not available in the selected language"))
	}

	getAssets(problem_id)

	archiveFilename := fmt.Sprintf("%s.zip", problem_id)
	unzipedDir := problem_id

	Unzip(archiveFilename, unzipedDir)

	_ = os.Remove(archiveFilename)

	moveFiles(cwd)

	if strings.Contains(getName(problem_id), "interactiv") {
		_ = os.Remove("Source.cpp")
		files, err := os.ReadDir(cwd)
		if err != nil {
			logErr(err)
		}
		var headerFilename string

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".h" {
				headerFilename = file.Name()
				break
			}
		}

		hFileContent, err := os.ReadFile(headerFilename)
		if err != nil {
			os.Chdir("..")
			os.Remove(newFolder)
			logErr(fmt.Errorf("error reading header file: %v", err))
		}

		funcDecls := extractFunctionDeclarations(string(hFileContent))

		cppFile, err := os.Create("Source.cpp")
		if err != nil {
			os.Chdir("..")
			os.Remove(newFolder)
			logErr(fmt.Errorf("error creating .cpp file: %v", err))
		}
		defer cppFile.Close()

		cppFile.WriteString("#include<iostream>\n#include \"myfunc.h\"\n\n")

		for _, decl := range funcDecls {
			cppFile.WriteString(decl + " {\n")
			cppFile.WriteString("\n")
			cppFile.WriteString("}\n\n")
		}

		cppFile.WriteString("\nint main(){\n\n\treturn 0;\n}")
	} else {
		auxStat := printStatement(problem_id, "RO", 2)
		if auxStat == "nolang" {
			auxStat = printStatement(problem_id, "EN", 2)
		}

		if !strings.Contains(auxStat, "stdin") && (lang == "cpp11" || lang == "cpp14" || lang == "cpp17" || lang == "cpp20") {
			_ = os.Remove("Source.cpp")
			cppFile, err := os.Create("Source.cpp")
			if err != nil {
				os.Chdir("..")
				os.Remove(newFolder)
				logErr(fmt.Errorf("error creating .cpp file: %v", err))
			}
			defer cppFile.Close()
			cppFile.WriteString(helloWorldprog[10])
		} else if !strings.Contains(auxStat, "stdin") && lang == "c" {
			_ = os.Remove("Source.c")
			cppFile, err := os.Create("Source.c")
			if err != nil {
				os.Chdir("..")
				os.Remove(newFolder)
				logErr(fmt.Errorf("error creating .c file: %v", err))
			}
			defer cppFile.Close()
			cppFile.WriteString(helloWorldprog[9])
		}
	}

	if cbp {
		fmt.Println(cbp)
		createCBPProject(problem_id)
	}

	if cmake {
		fmt.Println(cmake)
		createCMAKEProject(problem_id)
	}

}

func nrPbs() int {
	searchData := map[string]interface{}{
		"name_fuzzy": "",
		"offset":     0,
	}

	jsonData, err := json.Marshal(searchData)
	if err != nil {
		logErr(fmt.Errorf("error marshaling JSON: %v", err))
	}

	body, err := makeRequest("POST", URL_SEARCH, bytes.NewBuffer(jsonData), "2")
	if err != nil {
		logErr(err)
	}

	var data search
	err = json.Unmarshal(body, &data)
	if err != nil {
		logErr(fmt.Errorf("error unmarshaling JSON: %v", err))
		os.Exit(1)
	}

	return data.Data.Count

}

func getRandPb() {
	//min := 1
	max := nrPbs()
	fmt.Printf("Your random problem's ID: #%d\n", rand.IntN(max)+1)
}
