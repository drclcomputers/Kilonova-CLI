// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	problem "kilocli/cmd/problems"
	"kilocli/cmd/submission"

	utility "kilocli/cmd/utility"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
)

var codeBlocksProjectFile bool = false
var cMakeProjectFile bool = false

var InitProjectCmd = &cobra.Command{
	Use:   "init [Problem ID] [Language]",
	Short: "Create a project (statement, assets and source file for your chosen language)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		action := func() { initProject(args[0], args[1]) }
		if err := spinner.New().Title("Waiting ...").Action(action).Run(); err != nil {
			utility.LogError(err)
			return
		}
	},
}

func init() {

	InitProjectCmd.Flags().BoolVarP(&codeBlocksProjectFile, "codeblocksproject", "b", false, "Create a codeblocks project.")
	InitProjectCmd.Flags().BoolVarP(&cMakeProjectFile, "cmakeproject", "c", false, "Create a CMake project.")
}

func extractFunctionDeclarations(HeaderFileContent string) []string {
	Regexp := regexp.MustCompile(`\s*(int|float|char|double)\s+(\w+)\s*\((.*)\);`)
	matches := Regexp.FindAllStringSubmatch(HeaderFileContent, -1)

	var Declarations []string
	for _, match := range matches {
		if len(match) > 3 {
			Declaration := fmt.Sprintf("%s %s(%s)", match[1], match[2], match[3])
			Declarations = append(Declarations, Declaration)
		}
	}

	return Declarations
}

func configInteractiveProblem(CurrentWorkingDir, NewFolder string) {
	_ = os.Remove("Source.cpp")

	FilesInCWD, err := os.ReadDir(CurrentWorkingDir)
	if err != nil {
		utility.LogError(err)
		return
	}

	var headerFilename string
	for _, File := range FilesInCWD {
		if !File.IsDir() && filepath.Ext(File.Name()) == ".h" {
			headerFilename = File.Name()
			break
		}
	}

	HeaderFileContent, err := os.ReadFile(headerFilename)
	if err != nil {
		os.Chdir("..")
		os.Remove(NewFolder)
		utility.LogError(fmt.Errorf("error reading header file: %v", err))
		return
	}

	funcDecls := extractFunctionDeclarations(string(HeaderFileContent))

	CPPFile, err := os.Create("Source.cpp")
	if err != nil {
		os.Chdir("..")
		os.Remove(NewFolder)
		utility.LogError(fmt.Errorf("error creating .cpp file: %v", err))
		return
	}
	defer CPPFile.Close()

	CPPFile.WriteString("#include<iostream>\n#include \"myfunc.h\"\n\n")
	for _, decl := range funcDecls {
		CPPFile.WriteString(decl + " {\n")
		CPPFile.WriteString("\n")
		CPPFile.WriteString("}\n\n")
	}

	CPPFile.WriteString("\nint main(){\n\n\treturn 0;\n}")
}

func keyboardIOProblem(ProblemStatement, ProgrammingLanguage, NewFolder string) {
	if !strings.Contains(ProblemStatement, "stdin") && isCppLang(ProgrammingLanguage) {
		_ = os.Remove("Source.cpp")
		CPPFile, err := os.Create("Source.cpp")
		if err != nil {
			os.Chdir("..")
			os.Remove(NewFolder)
			utility.LogError(fmt.Errorf("error creating .cpp file: %v", err))
			return
		}
		defer CPPFile.Close()
		CPPFile.WriteString(utility.HelloWorldPrograms[10])
	} else if !strings.Contains(ProblemStatement, "stdin") && ProgrammingLanguage == "c" {
		_ = os.Remove("Source.c")
		CPPFile, err := os.Create("Source.c")
		if err != nil {
			os.Chdir("..")
			os.Remove(NewFolder)
			utility.LogError(fmt.Errorf("error creating .c file: %v", err))
			return
		}
		defer CPPFile.Close()
		CPPFile.WriteString(utility.HelloWorldPrograms[9])
	}
}

func GetCWDandCreateNewFolder(problemID, ProgrammingLanguage string) (string, string) {
	CurrentWorkingDir, err := os.Getwd()
	if err != nil {
		utility.LogError(fmt.Errorf("could not get current working directory! error: %v", err))
		return utility.ERROR, utility.ERROR
	}

	NewFolder := filepath.Join(CurrentWorkingDir, fmt.Sprintf("Problem_%s_Proj_%s", problemID, ProgrammingLanguage))
	if err := os.MkdirAll(NewFolder, os.ModePerm); err != nil {
		utility.LogError(fmt.Errorf("could not create project directory! error: %v", err))
		return utility.ERROR, utility.ERROR
	}

	if err := os.Chdir(NewFolder); err != nil {
		utility.LogError(fmt.Errorf("could not change directory to project dir! error: %v", err))
		return utility.ERROR, utility.ERROR
	}

	CurrentWorkingDir, err = os.Getwd()
	if err != nil {
		utility.LogError(fmt.Errorf("could not get current working directory! error: %v", err))
		return utility.ERROR, utility.ERROR
	}

	return CurrentWorkingDir, NewFolder

}

func AuxiliaryModifications(problemID, ProgrammingLanguage, CurrentWorkingDir, NewFolder string) {
	if strings.Contains(utility.GetAProblemName(problemID), "interactiv") {
		configInteractiveProblem(CurrentWorkingDir, NewFolder)
	} else {
		ProblemStatement, err := problem.PrintStatement(problemID, "RO", 2)
		if err != nil && err.Error() == utility.NOLANG {
			ProblemStatement, err = problem.PrintStatement(problemID, "EN", 2)
			if err != nil {
				utility.LogError(fmt.Errorf("error fetching problem statement: %v", err))
				return
			}
		}

		keyboardIOProblem(ProblemStatement, ProgrammingLanguage, NewFolder)
	}

	if codeBlocksProjectFile {
		createCodeBlocksProject(problemID)
	}

	if cMakeProjectFile {
		createCMakeProjectFile(problemID)
	}
}

func initProject(problemID, ProgrammingLanguage string) {
	CurrentWorkingDir, NewFolder := GetCWDandCreateNewFolder(problemID, ProgrammingLanguage)

	if !isLanguageSupported(problemID, ProgrammingLanguage) {
		os.Chdir("..")
		os.Remove(NewFolder)
		utility.LogError(errors.New("problem is not available in the selected language"))
		return
	}

	problem.GetAssets(problemID)

	archiveFilename := fmt.Sprintf("%s.zip", problemID)
	unzipedDir := problemID
	if err := unzip(archiveFilename, unzipedDir); err != nil {
		utility.LogError(fmt.Errorf("error unzipping file: %v", err))
		return
	}

	_ = os.Remove(archiveFilename)

	moveFiles(CurrentWorkingDir)

	AuxiliaryModifications(problemID, ProgrammingLanguage, CurrentWorkingDir, NewFolder)
}

func isLanguageSupported(problemID, ProgrammingLanguage string) bool {
	CurrentWorkingDir, err := os.Getwd()
	if err != nil {
		utility.LogError(fmt.Errorf("could not get current working directory: %v", err))
		return false
	}

	SupportedLangs := submission.CheckLanguages(problemID, 2)
	for _, SupportedLang := range SupportedLangs {
		if ProgrammingLanguage == SupportedLang {
			createSourceFile(CurrentWorkingDir, ProgrammingLanguage)
			return true
		}
	}
	return false
}

func isCppLang(ProgrammingLanguage string) bool {
	return ProgrammingLanguage == "cpp11" || ProgrammingLanguage == "cpp14" || ProgrammingLanguage == "cpp17" || ProgrammingLanguage == "cpp20"
}
