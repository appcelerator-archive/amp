package cli

import (
	"strings"

	"github.com/fatih/color"
)

// Theme is struct for terminal color functions
type Theme struct {
	Normal  *color.Color
	Info    *color.Color
	Warn    *color.Color
	Error   *color.Color
	Success *color.Color
}

var (
	// DarkTheme defines colors appropriate for a dark terminal
	DarkTheme = &Theme{
		Normal:  color.New(color.FgWhite),
		Info:    color.New(color.FgHiBlack),
		Warn:    color.New(color.FgYellow),
		Error:   color.New(color.FgRed),
		Success: color.New(color.FgGreen),
	}

	// LightTheme defines colors appropriate for a light terminal
	LightTheme = &Theme{
		Normal:  color.New(color.FgBlue),
		Info:    color.New(color.FgBlack),
		Warn:    color.New(color.FgYellow),
		Error:   color.New(color.FgRed),
		Success: color.New(color.FgGreen),
	}

	successPrefix = "Succes: "
	warningPrefix = "Warning: "
	errorPrefix   = "Error: "
)

// Console augments basic logging functions with a theme that can be applied
// for various standardized CLI output functions, such as Success and Error.
type Console struct {
	Logger
	theme *Theme
}

// NewConsole creates a CLI Console instance that writes to the provided stream.
func NewConsole(out *OutStream, verbose bool) *Console {
	// TODO: disabled pending investigation why cursor color isn't restored after printing colorized output
	color.NoColor = true

	return &Console{
		Logger: *NewLogger(out, verbose),
		theme:  DarkTheme,
	}
}

// Theme returns the current theme.
func (c Console) Theme() *Theme {
	return c.theme
}

// SetTheme sets the console theme.
func (c Console) SetTheme(theme *Theme) {
	c.theme = theme
}

func (c Console) SetThemeName(name string) {
	switch strings.TrimSuffix(strings.ToLower(strings.TrimSpace(name)), "theme") {
	case "light":
		c.SetTheme(LightTheme)
	default:
		c.SetTheme(DarkTheme)
	}
}

// OutStream returns the underlying OutStream that wraps stdout.
func (c Console) OutStream() *OutStream {
	return c.Logger.OutStream()
}

// Print prints args using Theme.Normal().
func (c Console) Print(args ...interface{}) {
	c.theme.Normal.Fprint(c.OutStream(), args...)
}

// Printf prints a formatted string using Theme.Normal().
func (c Console) Printf(format string, args ...interface{}) {
	c.theme.Normal.Fprintf(c.OutStream(), format, args...)
}

// Println prints args using Theme.Normal() and appends a newline.
func (c Console) Println(args ...interface{}) {
	c.theme.Normal.Fprintln(c.OutStream(), args...)
}

// Info prints args using Theme.Info().
func (c Console) Info(args ...interface{}) {
	c.theme.Info.Fprint(c.OutStream(), args...)
}

// Infof prints a formatted string using Theme.Info().
func (c Console) Infof(format string, args ...interface{}) {
	c.theme.Info.Fprintf(c.OutStream(), format, args...)
}

// Infoln prints args using Theme.Info() and appends a newline.
func (c Console) Infoln(args ...interface{}) {
	c.theme.Info.Fprintln(c.OutStream(), args...)
}

// Warn prints args using Theme.Warn().
func (c Console) Warn(args ...interface{}) {
	c.warn()
	c.theme.Warn.Fprint(c.OutStream(), args...)
}

// Warnf prints a formatted string using Theme.Warn().
func (c Console) Warnf(format string, args ...interface{}) {
	c.warn()
	c.theme.Warn.Fprintf(c.OutStream(), format, args...)
}

// Warnln prints args using Theme.Warn() and appends a newline.
func (c Console) Warnln(args ...interface{}) {
	c.warn()
	c.theme.Warn.Fprintln(c.OutStream(), args...)
}

// Error prints args using Theme.Error().
func (c Console) Error(args ...interface{}) {
	c.error()
	c.theme.Error.Fprint(c.OutStream(), args...)
}

// Errorf prints a formatted string using Theme.Error().
func (c Console) Errorf(format string, args ...interface{}) {
	c.error()
	c.theme.Error.Fprintf(c.OutStream(), format, args...)
}

// Errorln prints args using Theme.Error() and appends a newline.
func (c Console) Errorln(args ...interface{}) {
	c.error()
	c.theme.Error.Fprintln(c.OutStream(), args...)
}

// Success prints args Theme.Success().
func (c Console) Success(args ...interface{}) {
	c.success()
	c.theme.Success.Fprint(c.OutStream(), args...)
}

// Successf prints a formatted string using Theme.Success()
func (c Console) Successf(format string, args ...interface{}) {
	c.success()
	c.theme.Success.Fprintf(c.OutStream(), format, args...)
}

// Successln prints args using Theme.Success() and appends a newline.
func (c Console) Successln(args ...interface{}) {
	c.success()
	c.theme.Success.Fprintln(c.OutStream(), args...)
}

// prints success prefix
func (c Console) success() {
	c.theme.Success.Fprint(c.OutStream(), successPrefix)
}

// prints warning prefix
func (c Console) warn() {
	c.theme.Warn.Fprint(c.OutStream(), warningPrefix)
}

// prints error prefix
func (c Console) error() {
	c.theme.Error.Fprint(c.OutStream(), errorPrefix)
}
