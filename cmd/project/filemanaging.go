// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package project

import (
	"archive/zip"
	"fmt"
	"io"
	"kncli/internal"
	"os"
	"path/filepath"
	"strings"
)

func writeFile(filename, content string) {
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		internal.LogError(fmt.Errorf("failed to write file %s: %v", filename, err))
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

func moveFiles(srcDir, destDir string) error {
	return filepath.WalkDir(srcDir, func(Path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".md" || ext == ".pdf" || ext == ".h" {
			destPath := filepath.Join(destDir, entry.Name())

			if Path == destPath {
				return nil
			}

			return copyFile(Path, destPath)
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

	// Resolve destination to clean absolute path for Zip Slip protection
	destAbs, err := filepath.Abs(Destination)
	if err != nil {
		return fmt.Errorf("failed to resolve destination path: %w", err)
	}

	for _, File := range ZipFile.File {
		// Prevent Zip Slip: ensure extracted path stays within destination
		targetPath := filepath.Join(destAbs, File.Name)
		if !strings.HasPrefix(filepath.Clean(targetPath), destAbs+string(os.PathSeparator)) && targetPath != destAbs {
			return fmt.Errorf("illegal file path in zip: %s", File.Name)
		}

		if File.FileInfo().IsDir() {
			_ = os.MkdirAll(targetPath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		SourceFile, err := File.Open()
		if err != nil {
			return err
		}

		DestinationFile, err := os.Create(targetPath)
		if err != nil {
			SourceFile.Close()
			return err
		}

		_, err = io.Copy(DestinationFile, SourceFile)
		SourceFile.Close()
		DestinationFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func createCodeBlocksProject(ProjectName string) {
	codeBlocksFilename := fmt.Sprintf("%s.cbp", ProjectName)
	content := fmt.Sprintf(internal.XMLCBPStruct, ProjectName, ProjectName, ProjectName)
	writeFile(codeBlocksFilename, content)
}

func createCMakeProjectFile(ProjectName string) {
	content := fmt.Sprintf(internal.CMAKEStruct, ProjectName, ProjectName)
	writeFile(internal.CMakeFilename, content)
}

func createSourceFile(cwd, language string) {
	sourcePath := filepath.Join(cwd, "Source.")
	program, ext := getProgramByLanguage(language)

	sourcePath += ext
	writeFile(sourcePath, program)
}

// Maps Kilonova language identifiers to HelloWorld template programs.
// See internal.HelloWorldPrograms for the template array.
func getProgramByLanguage(language string) (program, extension string) {
	switch language {
	case "c":
		return internal.HelloWorldPrograms[0], "c"
	case "golang":
		return internal.HelloWorldPrograms[2], "go"
	case "kotlin":
		return internal.HelloWorldPrograms[3], "kt"
	case "nodejs":
		return internal.HelloWorldPrograms[4], "js"
	case "pascal":
		return internal.HelloWorldPrograms[5], "pas"
	case "php":
		return internal.HelloWorldPrograms[6], "php"
	case "python3":
		return internal.HelloWorldPrograms[7], "py"
	case "rust":
		return internal.HelloWorldPrograms[8], "rs"
	default:
		return internal.HelloWorldPrograms[1], "cpp"
	}
}
