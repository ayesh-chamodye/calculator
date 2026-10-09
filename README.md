# Go Qt Calculator

A simple calculator GUI built with Go and Qt (via `github.com/therecipe/qt`).

## Features

- Basic arithmetic: `+`, `-`, `*`, `/`, `%`
- Parentheses support
- Negative numbers
- Clear (C) and backspace (⌫) buttons
- Modern dark theme UI

## Requirements

- Go 1.18+
- Qt 5/6 development libraries installed on system

## Build & Run

```bash
go mod tidy
go build -o calculator.exe .
./calculator.exe
```

Or run directly:

```bash
go run .
```

## Dependencies

- `github.com/therecipe/qt` - Go bindings for Qt