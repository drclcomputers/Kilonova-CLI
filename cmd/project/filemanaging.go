// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package project

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	utility "kilocli/cmd/utility"
)

func copyFile(src, dest string) error {
	SourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer SourceFile.Close()

	DestinationFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer DestinationFile.Close()

	_, err = io.Copy(DestinationFile, SourceFile)
	return err
}

func moveFiles(RootDir string) error {
	return filepath.WalkDir(RootDir, func(Path string, EntryReadFromDir os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if EntryReadFromDir.IsDir() {
			return nil
		}

		Extension := strings.ToLower(filepath.Ext(EntryReadFromDir.Name()))
		if Extension == ".md" || Extension == ".pdf" || Extension == ".h" {
			DestinationPath := filepath.Join(RootDir, EntryReadFromDir.Name())

			if Path == DestinationPath {
				return nil
			}

			return copyFile(Path, DestinationPath)
		}
		return nil
	})
}

func unzip(Source string, Destination string) error {
	ZipFile, err := zip.OpenReader(Source)
	if err != nil {
		return err
	}
	defer ZipFile.Close()

	for _, File := range ZipFile.File {
		FilePathNotThePackage := filepath.Join(Destination, File.Name)

		if File.FileInfo().IsDir() {
			os.MkdirAll(FilePathNotThePackage, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(FilePathNotThePackage), os.ModePerm); err != nil {
			return err
		}

		SourceFile, err := File.Open()
		if err != nil {
			return err
		}
		defer SourceFile.Close()

		DestinationFile, err := os.Create(FilePathNotThePackage)
		if err != nil {
			return err
		}
		defer DestinationFile.Close()

		_, err = io.Copy(DestinationFile, SourceFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func createCodeBlocksProject(ProjectName string) {
	XMLCodeBlocksProjectFile := fmt.Sprintf(utility.XMLCBPStruct, ProjectName, ProjectName, ProjectName)

	File, err := os.Create(fmt.Sprintf("%s.cbp", ProjectName))
	if err != nil {
		utility.LogError(fmt.Errorf("error creating .cbp file: %v", err))
		return
	}
	defer File.Close()

	File.WriteString(XMLCodeBlocksProjectFile)
}

func createCMakeProjectFile(ProjectName string) {
	CMakeProjectFileTXT := fmt.Sprintf(utility.CMAKEStruct, ProjectName, ProjectName)

	File, err := os.Create("CMakeLists.txt")
	if err != nil {
		utility.LogError(fmt.Errorf("error creating \"CMakeLists.txt\": %v", err))
		return
	}
	defer File.Close()

	File.WriteString(CMakeProjectFileTXT)
}

// 0-C 1-CPP 2-GO 3-Kotlin 4-JS 5-PAS 6-PHP 7-Python 8-Rust
func createSourceFile(CurrentWorkingDir, ProgrammingLanguage string) {
	SourceFilename := filepath.Join(CurrentWorkingDir, "Source.")
	switch ProgrammingLanguage {
	case "c":
		SourceFilename += "c"
		os.WriteFile(SourceFilename, []byte(utility.HelloWorldPrograms[0]), 0644)
	case "golang":
		SourceFilename += "go"
		os.WriteFile(SourceFilename, []byte(utility.HelloWorldPrograms[2]), 0644)
	case "kotlin":
		SourceFilename += "kt"
		os.WriteFile(SourceFilename, []byte(utility.HelloWorldPrograms[3]), 0644)
	case "nodejs":
		SourceFilename += "js"
		os.WriteFile(SourceFilename, []byte(utility.HelloWorldPrograms[4]), 0644)
	case "pascal":
		SourceFilename += "pas"
		os.WriteFile(SourceFilename, []byte(utility.HelloWorldPrograms[5]), 0644)
	case "php":
		SourceFilename += "php"
		os.WriteFile(SourceFilename, []byte(utility.HelloWorldPrograms[6]), 0644)
	case "python3":
		SourceFilename += "py"
		os.WriteFile(SourceFilename, []byte(utility.HelloWorldPrograms[7]), 0644)
	case "rust":
		SourceFilename += "rs"
		os.WriteFile(SourceFilename, []byte(utility.HelloWorldPrograms[8]), 0644)
	default:
		SourceFilename += "cpp"
		os.WriteFile(SourceFilename, []byte(utility.HelloWorldPrograms[1]), 0644)
	}
}
