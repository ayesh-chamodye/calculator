package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Premium dark theme

// calcTheme overrides Fyne's default palette with a dark, high-contrast scheme
// used for backgrounds, buttons, text, accents and error states.
type calcTheme struct{ fyne.Theme }

// Color maps Fyne's semantic color names to this app's custom palette values.
func (t calcTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 11, G: 17, B: 32, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 27, G: 39, B: 57, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 226, G: 232, B: 240, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 59, G: 130, B: 246, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 45, G: 60, B: 82, A: 255}
	case theme.ColorNameError:
		return color.NRGBA{R: 251, G: 113, B: 133, A: 255}
	}
	return t.Theme.Color(n, v)
}

// Size customizes the base text size while deferring every other metric
// to the wrapped default theme.
func (t calcTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameText {
		return 15
	}
	return t.Theme.Size(n)
}

// Scientific expression parser
// Supports + - * / % ^, parentheses, functions, pi, e and degree mode.

type Parser struct {
	s       string
	pos     int
	degrees bool
}

func (p *Parser) spaces() {
	for p.pos < len(p.s) && unicode.IsSpace(rune(p.s[p.pos])) {
		p.pos++
	}
}

func (p *Parser) match(c byte) bool {
	p.spaces()
	if p.pos < len(p.s) && p.s[p.pos] == c {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) expression() (float64, error) {
	v, err := p.term()
	if err != nil {
		return 0, err
	}
	for {
		if p.match('+') {
			r, e := p.term()
			if e != nil {
				return 0, e
			}
			v += r
		} else if p.match('-') {
			r, e := p.term()
			if e != nil {
				return 0, e
			}
			v -= r
		} else {
			return v, nil
		}
	}
}

func (p *Parser) term() (float64, error) {
	v, err := p.unary()
	if err != nil {
		return 0, err
	}
	for {
		if p.match('*') {
			r, e := p.unary()
			if e != nil {
				return 0, e
			}
			v *= r
		} else if p.match('/') {
			r, e := p.unary()
			if e != nil {
				return 0, e
			}
			if r == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			v /= r
		} else if p.match('%') {
			r, e := p.unary()
			if e != nil {
				return 0, e
			}
			if r == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			v = math.Mod(v, r)
		} else {
			return v, nil
		}
	}
}

func (p *Parser) unary() (float64, error) {
	if p.match('+') {
		return p.unary()
	}
	if p.match('-') {
		v, err := p.unary()
		return -v, err
	}
	return p.power()
}

func (p *Parser) power() (float64, error) {
	v, err := p.primary()
	if err != nil {
		return 0, err
	}
	if p.match('^') {
		r, e := p.unary()
		if e != nil {
			return 0, e
		}
		v = math.Pow(v, r)
	}
	return v, nil
}

func (p *Parser) primary() (float64, error) {
	p.spaces()

	if p.match('(') {
		v, err := p.expression()
		if err != nil {
			return 0, err
		}
		if !p.match(')') {
			return 0, fmt.Errorf("missing closing parenthesis")
		}
		return v, nil
	}

	if p.pos >= len(p.s) {
		return 0, fmt.Errorf("incomplete expression")
	}

	c := p.s[p.pos]

	// Numeric literal: consume consecutive digits with at most one decimal point.
	if (c >= '0' && c <= '9') || c == '.' {
		start := p.pos
		dots := 0
		for p.pos < len(p.s) {
			ch := p.s[p.pos]
			if ch == '.' {
				dots++
				if dots > 1 {
					return 0, fmt.Errorf("invalid number")
				}
				p.pos++
			} else if ch >= '0' && ch <= '9' {
				p.pos++
			} else {
				break
			}
		}
		if start == p.pos || p.s[start:p.pos] == "." {
			return 0, fmt.Errorf("invalid number")
		}
		return strconv.ParseFloat(p.s[start:p.pos], 64)
	}

	// Identifier: either a math constant or a function that takes a parenthesized argument.
	if unicode.IsLetter(rune(c)) {
		start := p.pos
		for p.pos < len(p.s) && unicode.IsLetter(rune(p.s[p.pos])) {
			p.pos++
		}
		name := strings.ToLower(p.s[start:p.pos])

		switch name {
		case "pi":
			return math.Pi, nil
		case "e":
			return math.E, nil
		}

		if !p.match('(') {
			return 0, fmt.Errorf("expected '(' after %s", name)
		}
		v, err := p.expression()
		if err != nil {
			return 0, err
		}
		if !p.match(')') {
			return 0, fmt.Errorf("missing closing parenthesis")
		}

		switch name {
		case "sqrt":
			if v < 0 {
				return 0, fmt.Errorf("sqrt requires a non-negative value")
			}
			return math.Sqrt(v), nil
		case "abs":
			return math.Abs(v), nil
		case "ln":
			if v <= 0 {
				return 0, fmt.Errorf("ln requires a positive value")
			}
			return math.Log(v), nil
		case "log":
			if v <= 0 {
				return 0, fmt.Errorf("log requires a positive value")
			}
			return math.Log10(v), nil
		case "sin", "cos", "tan":
			if p.degrees {
				v *= math.Pi / 180
			}
			switch name {
			case "sin":
				return math.Sin(v), nil
			case "cos":
				return math.Cos(v), nil
			default:
				if math.Abs(math.Cos(v)) < 1e-14 {
					return 0, fmt.Errorf("tan is undefined here")
				}
				return math.Tan(v), nil
			}
		case "exp":
			return math.Exp(v), nil
		default:
			return 0, fmt.Errorf("unknown function: %s", name)
		}
	}

	return 0, fmt.Errorf("unexpected character: %c", c)
}

// calculate validates the input, evaluates it using the parser, and formats the result.
// If evaluation succeeds, it returns a formatted string, the raw float value, and nil.
func calculate(input string, degrees bool) (string, float64, error) {
	input = strings.TrimSpace(input)
	if input == "" || len(input) > 256 {
		return "", 0, fmt.Errorf("enter a valid expression")
	}

	// Parse and evaluate the trimmed expression.
	p := Parser{s: input, degrees: degrees}
	v, err := p.expression()
	if err != nil {
		return "", 0, err
	}
	// Ensure we consumed the entire input (no trailing tokens).
	p.spaces()

	if p.pos != len(p.s) {
		return "", 0, fmt.Errorf("invalid expression")
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "", 0, fmt.Errorf("result is not finite")
	}

	return strconv.FormatFloat(v, 'g', 12, 64), v, nil
}

// main initializes the Fyne application, builds the UI, and wires input handling
// (button presses and keyboard events) to the expression evaluator.
func main() {
	a := app.New()
	a.Settings().SetTheme(&calcTheme{Theme: theme.DefaultTheme()})

	w := a.NewWindow("NEXUS // Scientific Calculator")
	w.Resize(fyne.NewSize(850, 720))
	w.SetFixedSize(false)

	// Calculator state.
	var (
		expr        []rune
		fresh       bool
		degrees     = true
		memory      float64
		history     []string
		historyExpr []string
	)

	display := canvas.NewText("0", color.NRGBA{R: 248, G: 250, B: 252, A: 255})
	display.Alignment = fyne.TextAlignTrailing
	display.TextStyle = fyne.TextStyle{Bold: true}
	display.TextSize = 38

	preview := canvas.NewText("SCIENTIFIC MODE", color.NRGBA{R: 100, G: 116, B: 139, A: 255})
	preview.Alignment = fyne.TextAlignTrailing
	preview.TextSize = 12

	status := canvas.NewText("READY", color.NRGBA{R: 96, G: 165, B: 250, A: 255})
	status.TextSize = 11

	modeButton := widget.NewButton("DEG", nil)

	historyList := widget.NewList(
		func() int { return len(history) },
		func() fyne.CanvasObject {
			return widget.NewLabel("0 + 0 = 0")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(history[id])
		},
	)

	refresh := func() {
		s := string(expr)
		if s == "" {
			s = "0"
		}
		switch {
		case len(s) > 28:
			display.TextSize = 17
		case len(s) > 19:
			display.TextSize = 22
		case len(s) > 12:
			display.TextSize = 29
		default:
			display.TextSize = 38
		}
		display.Text = s
		display.Refresh()
	}

	setStatus := func(s string, isError bool) {
		if isError {
			status.Color = color.NRGBA{R: 251, G: 113, B: 133, A: 255}
		} else {
			status.Color = color.NRGBA{R: 96, G: 165, B: 250, A: 255}
		}
		status.Text = strings.ToUpper(s)
		status.Refresh()
	}

	var press func(string)

	press = func(key string) {
		switch key {
		case "C":
			expr = nil
			fresh = false
			setStatus("READY", false)

		case "⌫":
			if len(expr) > 0 {
				expr = expr[:len(expr)-1]
			}
			fresh = false
			setStatus("READY", false)

		case "=":
			input := string(expr)
			result, value, err := calculate(input, degrees)
			if err != nil {
				setStatus(err.Error(), true)
				return
			}

			history = append([]string{
				input + " = " + result,
			}, history...)
			historyExpr = append([]string{result}, historyExpr...)

			if len(history) > 100 {
				history = history[:100]
				historyExpr = historyExpr[:100]
			}
			historyList.Refresh()

			expr = []rune(result)
			fresh = true
			preview.Text = fmt.Sprintf("RESULT  •  %.12g", value)
			preview.Refresh()
			setStatus("CALCULATED", false)

		case "MC":
			memory = 0
			setStatus("MEMORY CLEARED", false)

		case "MR":
			expr = []rune(strconv.FormatFloat(memory, 'g', 12, 64))
			fresh = false
			setStatus("MEMORY RECALLED", false)

		case "M+":
			_, value, err := calculate(string(expr), degrees)
			if err != nil {
				setStatus(err.Error(), true)
				return
			}
			memory += value
			setStatus("MEMORY UPDATED", false)

		case "M-":
			_, value, err := calculate(string(expr), degrees)
			if err != nil {
				setStatus(err.Error(), true)
				return
			}
			memory -= value
			setStatus("MEMORY UPDATED", false)

		case "DEG/RAD":
			degrees = !degrees
			if degrees {
				modeButton.SetText("DEG")
			} else {
				modeButton.SetText("RAD")
			}
			setStatus("ANGLE MODE CHANGED", false)

		default:
			if fresh && !strings.ContainsAny(key, "+-*/%^") {
				expr = nil
			}
			fresh = false
			if len(expr)+len(key) <= 256 {
				expr = append(expr, []rune(key)...)
			}
			setStatus("READY", false)
		}

		refresh()
	}

	modeButton.OnTapped = func() { press("DEG/RAD") }

	historyList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(historyExpr) {
			expr = []rune(historyExpr[id])
			fresh = true
			refresh()
			setStatus("HISTORY LOADED", false)
		}
	}

	makeButton := func(label string, action string) *widget.Button {
		b := widget.NewButton(label, func() { press(action) })
		return b
	}

	// Scientific controls.
	scientificRows := [][]string{
		{"sin(", "cos(", "tan(", "sqrt("},
		{"ln(", "log(", "abs(", "exp("},
		{"x²", "^", "pi", "e"},
		{"MC", "MR", "M+", "M-"},
	}

	var scientificObjects []fyne.CanvasObject
	for _, row := range scientificRows {
		for _, label := range row {
			action := label
			switch label {
			case "x²":
				action = "^2"
			case "pi":
				action = "pi"
			case "e":
				action = "e"
			}
			scientificObjects = append(scientificObjects, makeButton(label, action))
		}
	}
	scientificGrid := container.NewGridWithColumns(4, scientificObjects...)

	// Basic keypad.
	basicRows := [][]string{
		{"C", "⌫", "(", ")"},
		{"7", "8", "9", "/"},
		{"4", "5", "6", "*"},
		{"1", "2", "3", "-"},
		{"0", ".", "%", "+"},
	}

	var basicObjects []fyne.CanvasObject
	for _, row := range basicRows {
		for _, label := range row {
			b := makeButton(label, label)
			if label == "C" {
				b.Importance = widget.DangerImportance
			}
			basicObjects = append(basicObjects, b)
		}
	}
	basicGrid := container.NewGridWithColumns(4, basicObjects...)

	equals := widget.NewButton("=   CALCULATE", func() { press("=") })
	equals.Importance = widget.HighImportance

	brand := canvas.NewText("N E X U S", color.NRGBA{R: 96, G: 165, B: 250, A: 255})
	brand.TextStyle = fyne.TextStyle{Bold: true}
	brand.TextSize = 15

	subtitle := canvas.NewText("SCIENTIFIC CALCULATOR  /  GO - BY Ayesh Dev", color.NRGBA{R: 100, G: 116, B: 139, A: 255})
	subtitle.TextSize = 10

	displayCard := container.NewVBox(
		container.NewPadded(display),
		container.NewPadded(preview),
		container.NewPadded(status),
	)

	header := container.NewVBox(
		container.NewHBox(brand, widget.NewLabel(""), modeButton),
		subtitle,
		widget.NewSeparator(),
		displayCard,
		widget.NewSeparator(),
	)

	calculatorPanel := container.NewBorder(
		header,
		container.NewVBox(equals),
		nil, nil,
		container.NewVScroll(container.NewVBox(
			scientificGrid,
			widget.NewSeparator(),
			basicGrid,
		)),
	)

	historyTitle := canvas.NewText("HISTORY", color.NRGBA{R: 96, G: 165, B: 250, A: 255})
	historyTitle.TextStyle = fyne.TextStyle{Bold: true}
	historyTitle.TextSize = 13

	clearHistory := widget.NewButton("Clear history", func() {
		history = nil
		historyExpr = nil
		historyList.Refresh()
		setStatus("HISTORY CLEARED", false)
	})

	historyPanel := container.NewBorder(
		container.NewVBox(historyTitle, widget.NewSeparator()),
		clearHistory,
		nil, nil,
		historyList,
	)

	split := container.NewHSplit(calculatorPanel, historyPanel)
	split.Offset = 0.70

	w.SetContent(container.NewPadded(split))

	// Keyboard support.
	w.Canvas().SetOnTypedRune(func(r rune) {
		if strings.ContainsRune("0123456789.+-*/%^()", r) {
			press(string(r))
		} else if r == 'x' || r == 'X' {
			press("*")
		} else if r == ',' {
			press(".")
		} else if r == '=' {
			press("=")
		}
	})

	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		switch e.Name {
		case fyne.KeyReturn, fyne.KeyEnter:
			press("=")
		case fyne.KeyBackspace:
			press("⌫")
		case fyne.KeyEscape, fyne.KeyDelete:
			press("C")
		}
	})

	w.ShowAndRun()
}
