// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var iniProjectCmd = &cobra.Command{
	Use:   "init [Problem ID] [Language]",
	Short: "Create a project, consisting of the statement, assets already downloaded and a source file for your chosen language.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		initProject(args[0], args[1])
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
	rootCmd.AddCommand(iniProjectCmd)
	rootCmd.AddCommand(getRandPbCmd)
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

var helloWorldprog = []string{
	`#include <stdio.h>
int main() {
    printf("Hello, World!\n");
    return 0;
}`,
	`#include <iostream>
int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}`,
	`package main
import "fmt"
func main() {
    fmt.Println("Hello, World!")
}`,
	`fun main() {
    println("Hello, World!")
}`,
	`console.log("Hello, World!");`,
	`program HelloWorld;
begin
    writeln('Hello, World!');
end.
`,
	`<?php
echo "Hello, World!\n";
?>
`, `print("Hello, World!")
`, `fn main() {
    println!("Hello, World!");
}
`,
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

func initProject(problem_id, lang string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not get current working directory! Err: %v\n", err)
		return
	}

	newFolder := filepath.Join(cwd, fmt.Sprintf("Problem_%s_Proj", problem_id))
	err = os.MkdirAll(newFolder, os.ModePerm)
	if err != nil {
		fmt.Printf("Could not create project directory! Err: %v\n", err)
		return
	}

	err = os.Chdir(newFolder)
	if err != nil {
		fmt.Printf("Could not change directory to project dir! Err: %v\n", err)
		return
	}

	cwd, err = os.Getwd()
	if err != nil {
		fmt.Printf("Could not get current working directory! Err: %v\n", err)
		return
	}

	ok := false

	supportedLangs := checklangs(problem_id, 2)
	for i := 0; i < len(supportedLangs); i++ {
		if lang == supportedLangs[i] {
			ok = true

			createSource(cwd, lang)
		}
	}

	if !ok {
		fmt.Println("Problem is not available in the selected language!")
		return
	}

	getAssets(problem_id)

	archiveFilename := fmt.Sprintf("%s.zip", problem_id)
	unzipedDir := problem_id

	Unzip(archiveFilename, unzipedDir)

	_ = os.Remove(archiveFilename)

	moveFiles(cwd)

}

func nrPbs() int {
	searchData := map[string]interface{}{
		"name_fuzzy": "",
		"offset":     0,
	}

	jsonData, err := json.Marshal(searchData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	body, err := makeRequest("POST", URL_SEARCH, bytes.NewBuffer(jsonData), "2")
	if err != nil {
		logErr(err)
	}

	var data search
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		os.Exit(1)
	}

	return data.Data.Count

}

func getRandPb() {
	//min := 1
	max := nrPbs()
	fmt.Printf("Your random problem's ID: #%d\n", rand.IntN(max)+1)
}
