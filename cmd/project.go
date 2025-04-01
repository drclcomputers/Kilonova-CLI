// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var iniProjectCmd = &cobra.Command{
	Use:   "initproject [Problem ID]",
	Short: "Create a project, consisting of the statement and assets already downloaded.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initProject(args[0])
	},
}

func init() {
	rootCmd.AddCommand(iniProjectCmd)
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
		if ext == ".md" || ext == ".pdf" {
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

func initProject(problem_id string) {
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

	getAssets(problem_id)

	archiveFilename := fmt.Sprintf("%s.zip", problem_id)
	unzipedDir := problem_id

	Unzip(archiveFilename, unzipedDir)

	_ = os.Remove(archiveFilename)

	cwd, err = os.Getwd()
	if err != nil {
		fmt.Printf("Could not get current working directory! Err: %v\n", err)
		return
	}

	moveFiles(cwd)
}
