# 🚀 Kilonova-CLI

A sleek command-line client for the [Kilonova](https://kilonova.ro) competitive programming platform, letting you search problems, submit solutions, and track results — all from your terminal. ✨

---

## 📦 Project Structure

```
.
├── cmd/              # CLI command definitions (using Cobra)
├── internal/         # General Functions 
├── main.go           # Application entry point
├── go.mod            # Module dependencies
├── go.sum            # Dependency checksums
├── .gitignore        # Git ignored files
├── LICENSE           # MIT License
└── README.md         # This file
```

---

## ✨ Features

- 🔐 **User Authentication** — Securely log into your Kilonova account.
- 🔍 **Problem Search** — Find problems by keywords or IDs.
- 📤 **Solution Submission** — Submit your code directly via CLI.
- 📈 **Submission Info** — Retrieve and display your submission results.
- 🎨 **Pretty Terminal UI** — Thanks to `glamour`, `bubbletea`, and `lipgloss`.

---

## 📥 Installation

> **Requires Go 1.24.1+**

Clone and build the project:
```bash
git clone https://github.com/drclcomputers/Kilonova-CLI
cd Kilonova-CLI
go build
./kncli
```

---

## 🛠️ Usage

List all available commands:
```bash
./kncli help
```

Examples:
```bash
./kilocli signin      # Log in to Kilonova
./kilocli search      # Search for problems
./kilocli submit      # Submit a solution
./kilocli submission  # Check submission status
```

---

## 📚 Dependencies

- [Cobra](https://github.com/spf13/cobra) 🐍 — CLI framework
- [Bubbletea](https://github.com/charmbracelet/bubbletea) 🫖 — Terminal UI toolkit
- [Glamour](https://github.com/charmbracelet/glamour) ✨ — Markdown rendering in terminal
- [Lipgloss](https://github.com/charmbracelet/lipgloss) 💅 — Style definitions for Bubbletea apps

---

## 👨‍💻 Contributing

Contributions welcome! 💙  
- Fork this repo  
- Create a feature branch  
- Commit your changes  
- Open a pull request 🚀

Or report bugs/suggestions via [issues](https://github.com/drclcomputers/Kilonova-CLI/issues).

---

## 📄 License

This project is licensed under the MIT License. See `LICENSE` for details.

---

## 📌 Notes

- This is an independent open-source project and not officially affiliated with Kilonova.
- Built entirely in **Go** 💙.

