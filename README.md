# Kilonova-CLI

Kilonova-CLI is a command-line interface (CLI) client designed for interacting with the [Kilonova](https://kilonova.ro/) competitive programming platform. It enables users to manage their accounts, search for problems, submit solutions, and retrieve submission results directly from the terminal.

## Features

- **User Authentication**: Log in to your Kilonova account securely from the CLI.
- **Problem Search**: Find programming problems by keywords or categories.
- **Solution Submission**: Submit your code solutions to Kilonova problems.
- **Submission Status**: Check the results and status of your submissions.

## Installation

To install Kilonova-CLI, ensure you have [Go](https://golang.org/dl/) installed on your system. Then, run:

```sh
git clone https://github.com/drclcomputers/Kilonova-CLI
go build
./kilocli
```


This command will download and build the Kilonova-CLI.


## Usage

Once installed, you can use Kilonova-CLI to interact with the Kilonova platform.

### Authenticate

To authenticate with your Kilonova account:

```sh
<PROGRAM> -signin <USERNAME> <PASSWORD>
```

This command uses the credentials provided to log in.

### Search for Problems

To search for problems containing a specific keyword:

```sh
<PROGRAM> -search <NAME or ID>
```


### View Problem Statement

To view the statement of a specific problem:

```sh
<PROGRAM> -statement <ID> <RO or EN>
```

Replace `<RO or EN>` with the language you want to view the statement in.

### Submit a Solution

To submit a solution to a problem:

```sh
<PROGRAM> -submit <ID> <LANGUAGE> <FILENAME>
```

Replace `<LANGUAGE>` with one of the available languages which can be seen by running
```sh
<PROGRAM> -langs <ID>
```
and `<FILENAME>` with the path to your solution file.

### Check Submission Status

To check the status of your submissions:

```sh
<PROGRAM> -submissions <ID>
```

This command retrieves and displays the results of your submissions.

## Contributing

Contributions to Kilonova-CLI are welcome! If you find a bug or have a feature request, please open an issue on the [GitHub repository](https://github.com/drclcomputers/Kilonova-CLI). For code contributions, fork the repository and submit a pull request with your changes.

## License

Kilonova-CLI is licensed under the [MIT License](LICENSE).

---

*Note: Kilonova-CLI is an independent project and is not officially affiliated with the Kilonova platform.*
