// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package utility

type RequestType int

const (
	API_URL = "https://kilonova.ro/api/"

	URL_LOGIN          = API_URL + "auth/login"
	URL_LOGOUT         = API_URL + "auth/logout"
	URL_EXTEND_SESSION = API_URL + "auth/extendSession"

	URL_SEARCH  = API_URL + "problem/search"
	URL_PROBLEM = API_URL + "problem/%s/"

	URL_SELF          = API_URL + "user/self/"
	URL_SELF_PROBLEMS = API_URL + "user/self/solvedProblems"
	URL_SELF_SET_BIO  = API_URL + "user/self/setBio"

	URL_CHANGE_EMAIL  = API_URL + "user/changeEmail"
	URL_CHANGE_PASS   = API_URL + "user/changePassword"
	URL_CHANGE_NAME   = API_URL + "user/updateName"
	URL_RESEND_MAIL   = API_URL + "user/resendEmail"
	URL_DELETE_USER   = API_URL + "user/moderation/deleteUser"
	URL_USER          = API_URL + "user/byID/%s"
	URL_USER_PROBLEMS = API_URL + "user/byID/%s/solvedProblems"

	URL_LANGS_PB = API_URL + "problem/%s/languages"

	URL_SUBMIT                     = API_URL + "submissions/submit"
	URL_LATEST_SUBMISSION          = API_URL + "submissions/getByID?id=%d"
	URL_SUBMISSION_LIST            = API_URL + "submissions/get?ascending=false&limit=50&offset=%d&ordering=id&problem_id=%s&user_id=%s"
	URL_SUBMISSION_LIST_NO_FILTER  = API_URL + "submissions/get?ascending=false&limit=50&offset=%d&ordering=id"
	URL_SUBMISSION_LIST_NO_PROBLEM = API_URL + "submissions/get?ascending=false&limit=50&offset=%d&ordering=id&user_id=%s"
	URL_SUBMISSION_LIST_NO_USER    = API_URL + "submissions/get?ascending=false&limit=50&offset=%d&ordering=id&problem_id=%s"

	URL_CONTEST                     = API_URL + "contest/%s"
	URL_CONTEST_UPDATE              = API_URL + "contest/%s/update"
	URL_CONTEST_CREATE              = API_URL + "contest/create"
	URL_CONTEST_DELETE              = API_URL + "contest/%s/delete"
	URL_CONTEST_REGISTER            = API_URL + "contest/%s/register"
	URL_CONTEST_START               = API_URL + "contest/%s/startRegistration"
	URL_CONTEST_ANNOUNCEMENTS       = API_URL + "contest/%s/announcements"
	URL_CONTEST_CREATE_ANNOUNCEMENT = API_URL + "contest/%s/createAnnouncement"
	URL_CONTEST_UPDATE_ANNOUNCEMENT = API_URL + "contest/%s/updateAnnouncement"
	URL_CONTEST_DELETE_ANNOUNCEMENT = API_URL + "contest/%s/deleteAnnouncement"
	URL_CONTEST_ASK_QUESTION        = API_URL + "contest/%s/askQuestion"
	URL_CONTEST_RESPOND_QUESTION    = API_URL + "contest/%s/answerQuestion"
	URL_CONTEST_YOUR_QUESTIONS      = API_URL + "contest/%s/questions"
	URL_CONTEST_ALL_QUESTIONS       = API_URL + "contest/%s/allQuestions"
	URL_CONTEST_UPDATE_PROBLEMS     = API_URL + "contest/%s/update/problems"
	URL_CONTEST_PROBLEMS            = API_URL + "contest/%s/problems"
	URL_CONTEST_LEADERBOARD         = API_URL + "contest/%s/leaderboard"

	STAT_FILENAME_RO = "statement-ro.md"
	STAT_FILENAME_EN = "statement-en.md"

	URL_STATEMENT      = API_URL + "problem/%s/get/attachmentByName/%s"
	URL_ASSETS         = "https://kilonova.ro/assets/problem/%s/problemArchive?tests=true&attachments=true&private_attachments=false&details=true&tags=true&editors=true&submissions=false&all_submissions=false"
	URL_CONTEST_ASSETS = "https://kilonova.ro/assets/contest/%s/leaderboard.csv"

	UserAgent = "KilonovaCLIClient/0.2"

	XMLCBPStruct = `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<CodeBlocks_project_file>
	<FileVersion major="1" minor="6" />
	<Project>
		<Option title="%s" />
		<Option pch_mode="2" />
		<Option compiler="gcc" />
		<Build>
			<Target title="Debug">
				<Option output="bin/Debug/%s" prefix_auto="1" extension_auto="1" />
				<Option object_output="obj/Debug/" />
				<Option type="1" />
				<Option compiler="gcc" />
				<Compiler>
					<Add option="-g" />
				</Compiler>
			</Target>
			<Target title="Release">
				<Option output="bin/Release/%s" prefix_auto="1" extension_auto="1" />
				<Option object_output="obj/Release/" />
				<Option type="1" />
				<Option compiler="gcc" />
				<Compiler>
					<Add option="-O2" />
				</Compiler>
				<Linker>
					<Add option="-s" />
				</Linker>
			</Target>
		</Build>
		<Compiler>
			<Add option="-Wall" />
			<Add option="-fexceptions" />
		</Compiler>
		<Unit filename="Source.cpp" />
		<Extensions />
	</Project>
</CodeBlocks_project_file>
`
	CMAKEStruct = `cmake_minimum_required(VERSION 3.10)
project(%s VERSION 1.0 LANGUAGES CXX)
add_executable(%s Source.cpp)`

	SubmissionTemplate = `Submission ID: #{{.ID}}
Created: {{.CreatedAt}}
Language: {{.Language}}
Score: {{.Score}}

Max memory: {{.MaxMemory}}KB
Max time: {{.MaxTime}}s
Compile error: {{.CompileError}}
Compile message: {{.CompileMessage}}

Code:
{{.Code}}
`

	TemplatePattern = `Name: {{.Name}}
ID: #{{.ID}}
Time Limit: {{.TimeLimit}}s
Memory Limit: {{.MemoryLimit}}KB
Source Size: {{.SourceSize}}KB
Credits: {{.Credits}}
`

	CMakeFilename = "CMakeLists.txt"

	// Request Types

	RequestNone          RequestType = iota
	RequestFormAuth                  // 1
	RequestJSON                      // 2
	RequestFormGuest                 // 3
	RequestDownloadZip               // 4
	RequestInfo                      // 5
	RequestMultipartForm             // 6

	// Colors
	RED   = "\033[31m"
	WHITE = "\033[0m"

	// Others
	BOOLTRUE  = "true"
	BOOLFALSE = "false"
	NOLANG    = "nolang"
	SUCCESS   = "success"
	ERROR     = "error"
)

var HelloWorldPrograms = []string{
	// C
	`#include <stdio.h>
int main() {
	printf("Hello, World!\n");
	return 0;
}`,

	// C++
	`#include <iostream>
int main() {
	std::cout << "Hello, World!" << std::endl;
	return 0;
}`,

	// GO
	`package main
import "fmt"
func main() {
	fmt.Println("Hello, World!")
}`,

	// Kotlin
	`fun main() {
	println("Hello, World!")
	}`,

	// JavaScipt
	`console.log("Hello, World!");`,

	// Pascal
	`program HelloWorld;
begin
	writeln('Hello, World!');
end.`,

	// PHP
	`<?php
	echo "Hello, World!\n";
	?>`,

	// Python
	`print("Hello, World!")`,

	// Rust
	`fn main() {
	println!("Hello, World!");
	}`,

	// C with file I/O
	`#include <stdio.h>

int main() {
	FILE *file = fopen("example.txt", "r");
	char ch;
	while ((ch = fgetc(file)) != EOF) {
		putchar(ch);
	}
	return 0;
}`,

	// C++ with file I/O
	`#include <iostream>
#include <fstream>

int main() {
	std::ifstream file("example.txt");
	char ch;
	while (file.get(ch)) {
		std::cout << ch;
	}
	return 0;
}`,
}

var Replacements = map[string]string{
	"$":             "",
	` \ `:           "",
	`\ldots`:        "...",
	`\leq`:          "≤",
	`\geq`:          "≥",
	`\el`:           "",
	`\in`:           "∈",
	`\le`:           "≤",
	`\qe`:           "≥",
	`\pm`:           "±",
	`\cdot`:         "•",
	`\sum_`:         "Σ ",
	`\displaystyle`: "",
	`\times`:        "x",
	`\%`:            "%",
}

var ReplacementsRegex = []string{
	`\\text{(.*?)}`,
	`\\texttt{(.*?)}`,
	`\\bm{(.*?)}`,
	`\\textit{(.*?)}`,
	`\\rule\{[^}]+\}\{[^}]+\}`,
}
