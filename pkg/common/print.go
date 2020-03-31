package common

import "fmt"

// PrintFn is a wrapper function to print using an ansi color
type PrintFn func(...interface{}) string

// PrintColor represents a color manager
type PrintColor struct {
	Colors   []PrintFn
	ColorIdx int
}

// NewPrintColor creates a new color manager
func NewPrintColor() *PrintColor {
	pc := &PrintColor{}
	pc.load()

	return pc
}

func (p *PrintColor) load() {
	p.Colors = []PrintFn{
		PrintClr("\033[1;32m%s\033[0m"), // Green
		PrintClr("\033[1;33m%s\033[0m"), // Yellow
		PrintClr("\033[1;34m%s\033[0m"), // Purple
		PrintClr("\033[1;35m%s\033[0m"), // Magenta
		PrintClr("\033[1;36m%s\033[0m"), // Teal
		PrintClr("\033[1;31m%s\033[0m"), // Red
	}
}

// Get the next color
func (p *PrintColor) Get() PrintFn {
	fn := p.Colors[p.ColorIdx]
	p.ColorIdx++

	if p.ColorIdx == len(p.Colors) {
		p.ColorIdx = 0
	}

	return fn
}

// GetNoColor disable ansi color
func (p *PrintColor) GetNoColor() PrintFn {
	return func(args ...interface{}) string {
		return fmt.Sprint(args...)
	}
}

// PrintClr returns a the wrapper func using a given color
func PrintClr(color string) PrintFn {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(color,
			fmt.Sprint(args...))
	}
	return sprint
}
