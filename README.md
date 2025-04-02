# Kilonova-CLI

Kilonova-CLI is a command-line interface (CLI) client designed for interacting with the [Kilonova](https://kilonova.ro/) competitive programming platform. It enables users to view statements, search for problems, submit solutions, and retrieve submission results directly from the terminal.

## Features

- **User Authentication**: Log in to your Kilonova account securely from the CLI.
- **Problem Search**: Find programming problems by keywords or IDs.
- **Solution Submission**: Submit your code solutions to Kilonova problems.
- **Submission Status**: Check the results and status of your submissions.

## Installation

To install Kilonova-CLI, ensure you have [Go 1.24.1](https://golang.org/dl/) installed on your system. Then, run:

```sh
git clone https://github.com/drclcomputers/Kilonova-CLI
cd Kilonova-CLI
go build
./kilocli
```


This command will download and build the Kilonova-CLI.


## Usage

Once installed, you can use Kilonova-CLI to interact with the Kilonova platform. For detailed instructions on available commands and their usage, please run:
```sh
./kilocli help
```

## Contributing

Contributions to Kilonova-CLI are welcome! If you find a bug or have a feature request, please open an issue on the [GitHub repository](https://github.com/drclcomputers/Kilonova-CLI). For code contributions, fork the repository and submit a pull request with your changes.

## License

Kilonova-CLI is licensed under the [MIT License](LICENSE).

---

## Used projects

[Cobra](https://github.com/spf13/cobra)

[Glamour](https://github.com/charmbracelet/glamour)

[Bubbletea](https://github.com/charmbracelet/bubbletea)

[Lipgloss](https://github.com/charmbracelet/lipgloss)

---

*Note: Kilonova-CLI is an independent project and is not officially affiliated with the Kilonova platform.*
