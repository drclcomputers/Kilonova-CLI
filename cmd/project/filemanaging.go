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

func writeFile(filename, content string) {
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		utility.LogError(fmt.Errorf("failed to write file %s: %v", filename, err))
	}
}

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

	codeBlocksFilename := fmt.Sprintf("%s.cbp", ProjectName)
	File, err := os.Create(codeBlocksFilename)
	if err != nil {
		utility.LogError(fmt.Errorf("error creating .cbp file: %v", err))
		return
	}
	defer File.Close()

	writeFile(codeBlocksFilename, XMLCodeBlocksProjectFile)
}

func createCMakeProjectFile(ProjectName string) {
	CMakeProjectFileTXT := fmt.Sprintf(utility.CMAKEStruct, ProjectName, ProjectName)

	File, err := os.Create(utility.CMakeFilename)
	if err != nil {
		utility.LogError(fmt.Errorf("error creating \"CMakeLists.txt\": %v", err))
		return
	}
	defer File.Close()

	writeFile(utility.CMakeFilename, CMakeProjectFileTXT)
}

func createSourceFile(cwd, language string) {
	sourcePath := filepath.Join(cwd, "Source.")
	program, ext := getProgramByLanguage(language)

	sourcePath += ext
	writeFile(sourcePath, program)
}

// 0-C 1-CPP 2-GO 3-Kotlin 4-JS 5-PAS 6-PHP 7-Python 8-Rust
func getProgramByLanguage(language string) (program, extension string) {
	switch language {
	case "c":
		return utility.HelloWorldPrograms[0], "c"
	case "golang":
		return utility.HelloWorldPrograms[2], "go"
	case "kotlin":
		return utility.HelloWorldPrograms[3], "kt"
	case "nodejs":
		return utility.HelloWorldPrograms[4], "js"
	case "pascal":
		return utility.HelloWorldPrograms[5], "pas"
	case "php":
		return utility.HelloWorldPrograms[6], "php"
	case "python3":
		return utility.HelloWorldPrograms[7], "py"
	case "rust":
		return utility.HelloWorldPrograms[8], "rs"
	default:
		return utility.HelloWorldPrograms[1], "cpp"
	}
}
