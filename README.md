# NEXUS Scientific Calculator

A modern scientific calculator written in Go using the Fyne UI toolkit. It provides a sleek dark interface, mathematical expression evaluation, scientific functions, history tracking, and memory operations in a desktop application.

## Features

- Basic arithmetic: `+`, `-`, `*`, `/`, `%`
- Parentheses and nested expressions
- Exponentiation with `^`
- Scientific functions: `sin`, `cos`, `tan`, `sqrt`, `log`, `ln`, `abs`, `exp`
- Constants: `pi`, `e`
- Degree and radian mode toggle
- Memory functions: `MC`, `MR`, `M+`, `M-`
- Calculation history with quick reload
- Clean dark theme and responsive display sizing
- Keyboard input support for fast calculations

## Project Overview

This project is a single-file Go application with a custom parser and Fyne-based UI. The calculator evaluates expressions directly from user input and formats results for display.

## Requirements

- Go 1.27.2 (as specified in `go.mod`) or a compatible newer Go version
- Fyne dependencies managed automatically by Go modules

## Installation

Clone the repository:

```bash
git clone https://github.com/ayesh-chamodye/calculator.git
cd calculator
```

Install dependencies:

```bash
go mod tidy
```

Run the app:

```bash
go run .
```

Build a binary:

```bash
go build -o calculator.exe .
```

## Usage

Enter expressions using the numeric keypad or keyboard:

```text
2 + 3 * 4
sin(30)
sqrt(81)
2 ^ 8
pi * 2
log(100)
```

Examples:

- `sin(30)`
- `cos(60)`
- `tan(45)`
- `sqrt(16)`
- `ln(2.718281828)`
- `abs(-42)`
- `exp(1)`

Use the `DEG/RAD` button to switch angle units when working with trigonometric functions.

## Keyboard Support

- Number keys and operators are accepted directly
- `Enter` evaluates the expression
- `Backspace` removes the last character
- `Delete` or `Esc` clears the current input
- `x` or `X` is interpreted as multiplication

## Project Structure

```text
.
├── main.go
├── go.mod
├── go.sum
├── README.md
└── ModernCalc.exe
```

- `main.go` contains the calculator parser, UI, and event handling.
- `go.mod` defines the module and dependencies.
- `ModernCalc.exe` is a compiled Windows executable included in the repo.

## License

This project does not currently declare a license. If you plan to distribute or reuse it, add an appropriate open-source license before publishing.

## Notes

The app is designed as a desktop calculator and is especially suited for expression-heavy scientific math workflows. The parser supports a broad set of scientific operations while keeping the interface simple and responsive.
