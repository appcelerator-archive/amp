package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
)

// Theme is struct for terminal color functions
type Theme struct {
	Normal  *color.Color
	Info    *color.Color
	Warn    *color.Color
	Error   *color.Color
	Success *color.Color
}

const (
	// Padding for tabwriter
	Padding = 3
)

var (
	// DarkTheme defines colors appropriate for a dark terminal
	DarkTheme = &Theme{
		Normal:  color.New(),
		Info:    color.New(color.FgHiBlack),
		Warn:    color.New(color.FgYellow),
		Error:   color.New(color.FgRed),
		Success: color.New(color.FgGreen),
	}

	// LightTheme defines colors appropriate for a light terminal
	LightTheme = &Theme{
		Normal:  color.New(),
		Info:    color.New(color.FgBlack),
		Warn:    color.New(color.FgYellow),
		Error:   color.New(color.FgRed),
		Success: color.New(color.FgGreen),
	}

	successPrefix = "Succes: "
	warningPrefix = "Warning: "
	errorPrefix   = "Error: "
	fatalPrefix   = "Fatal: "
)

// Console augments basic logging functions with a theme that can be applied
// for various standardized CLI output functions, such as Success and Error.
type Console struct {
	Logger
	theme *Theme
}

// NewConsole creates a CLI Console instance that writes to the provided stream.
func NewConsole(out *OutStream, verbose bool) *Console {
	// TODO: investigate why cursor color isn't restored after printing colorized output
	// Uncomment the following if necessary to disable color output
	// color.NoColor = true

	return &Console{
		Logger: *NewLogger(out, verbose),
		theme:  DarkTheme,
	}
}

// Theme returns the current theme.
func (c *Console) Theme() *Theme {
	return c.theme
}

// SetTheme sets the console theme.
func (c *Console) SetTheme(theme *Theme) {
	c.theme = theme
}

// SetThemeName sets the console theme by name.
func (c *Console) SetThemeName(name string) {
	switch strings.TrimSuffix(strings.ToLower(strings.TrimSpace(name)), "theme") {
	case "light":
		c.SetTheme(LightTheme)
	default:
		c.SetTheme(DarkTheme)
	}
}

// OutStream returns the underlying OutStream that wraps stdout.
func (c *Console) OutStream() *OutStream {
	return c.Logger.OutStream()
}

// Print prints args using Theme.Normal().
func (c *Console) Print(args ...interface{}) {
	c.theme.Normal.Fprint(c.OutStream(), args...) // nolint
}

// Printf prints a formatted string using Theme.Normal().
func (c *Console) Printf(format string, args ...interface{}) {
	c.theme.Normal.Fprintf(c.OutStream(), format, args...) // nolint
}

// Println prints args using Theme.Normal() and appends a newline.
func (c *Console) Println(args ...interface{}) {
	c.theme.Normal.Fprintln(c.OutStream(), args...) // nolint
}

// Info prints args using Theme.Info().
func (c *Console) Info(args ...interface{}) {
	c.theme.Info.Fprint(c.OutStream(), args...) // nolint
}

// Infof prints a formatted string using Theme.Info().
func (c *Console) Infof(format string, args ...interface{}) {
	c.theme.Info.Fprintf(c.OutStream(), format, args...) // nolint
}

// Infoln prints args using Theme.Info() and appends a newline.
func (c *Console) Infoln(args ...interface{}) {
	c.theme.Info.Fprintln(c.OutStream(), args...) // nolint
}

// Warn prints args using Theme.Warn().
func (c *Console) Warn(args ...interface{}) {
	c.warn()
	c.theme.Warn.Fprint(c.OutStream(), args...) // nolint
}

// Warnf prints a formatted string using Theme.Warn().
func (c *Console) Warnf(format string, args ...interface{}) {
	c.warn()
	c.theme.Warn.Fprintf(c.OutStream(), format, args...) // nolint
}

// Warnln prints args using Theme.Warn() and appends a newline.
func (c *Console) Warnln(args ...interface{}) {
	c.warn()
	c.theme.Warn.Fprintln(c.OutStream(), args...) // nolint
}

// Error prints args using Theme.Error().
func (c *Console) Error(args ...interface{}) {
	c.error()
	c.theme.Error.Fprint(c.OutStream(), args...) // nolint
}

// Errorf prints a formatted string using Theme.Error().
func (c *Console) Errorf(format string, args ...interface{}) {
	c.error()
	c.theme.Error.Fprintf(c.OutStream(), format, args...) // nolint
}

// Errorln prints args using Theme.Error() and appends a newline.
func (c *Console) Errorln(args ...interface{}) {
	c.error()
	c.theme.Error.Fprintln(c.OutStream(), args...) // nolint
}

// Success prints args Theme.Success().
func (c *Console) Success(args ...interface{}) {
	c.success()
	c.theme.Success.Fprint(c.OutStream(), args...) // nolint
}

// Successf prints a formatted string using Theme.Success()
func (c *Console) Successf(format string, args ...interface{}) {
	c.success()
	c.theme.Success.Fprintf(c.OutStream(), format, args...) // nolint
}

// Successln prints args using Theme.Success() and appends a newline.
func (c *Console) Successln(args ...interface{}) {
	c.success()
	c.theme.Success.Fprintln(c.OutStream(), args...) // nolint
}

// Fatal prints args Theme.Error() and exits with code 1.
func (c *Console) Fatal(args ...interface{}) {
	c.fatal()
	c.theme.Error.Fprint(c.OutStream(), args...) // nolint
	os.Exit(1)
}

// Fatalf prints a formatted string using Theme.Error() and exits with code 1
func (c *Console) Fatalf(format string, args ...interface{}) {
	c.fatal()
	c.theme.Error.Fprintf(c.OutStream(), format, args...) // nolint
	os.Exit(1)
}

// Fatalln prints args using Theme.Error(),appends a newline and exits with code 1
func (c *Console) Fatalln(args ...interface{}) {
	c.fatal()
	c.theme.Error.Fprintln(c.OutStream(), args...) // nolint
	os.Exit(1)
}

// prints success prefix
func (c *Console) success() {
	c.theme.Success.Fprint(c.OutStream(), successPrefix) // nolint
}

// prints warning prefix
func (c *Console) warn() {
	c.theme.Warn.Fprint(c.OutStream(), warningPrefix) // nolint
}

// prints error prefix
func (c *Console) error() {
	c.theme.Error.Fprint(c.OutStream(), errorPrefix) // nolint
}

// prints fatal prefix
func (c *Console) fatal() {
	c.theme.Error.Fprint(c.OutStream(), fatalPrefix) // nolint
}

// GetInput gets input from standard input and returns it
func (c *Console) GetInput(prompt string) (in string) {
	c.Printf("%s: ", prompt)
	fmt.Scanln(&in)
	in = strings.TrimSpace(in)
	return in
}

// GetSilentInput gets input from standard input, without displaying characters, and returns it
func (c *Console) GetSilentInput(prompt string) (in string) {
	c.Printf("%s: ", prompt)
	bytes, err := gopass.GetPasswd()
	if err != nil {
		c.Fatalln(err.Error())
	}
	in = string(bytes)
	return in
}
