package cli

import (
	"bytes"
	"math/rand"
	"strconv"
	"text/template"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// int initialises the seed in math/rand.
func init() {
	rand.Seed(time.Now().Unix())
}

// templating creates, parses and executes a template to generate unique values.
func templating(input string, cache map[string]string) (output string, err error) {
	// Create a new template and parse the input using it.
	var tmpl *template.Template
	tmpl, err = template.New("Command").Parse(input)
	if err != nil {
		return
	}

	// Custom function to create a unique name with a randomly generated string.
	name := func(in string) string {
		if val, ok := cache[in]; ok {
			return val
		}
		out := in + randUniq(10)
		cache[in] = out
		return out
	}

	// Custom function to randomly generate a port number.
	port := func(in string, min, max int) string {
		if val, ok := cache[in]; ok {
			return val
		}
		out := randPort(min, max)
		cache[in] = out
		return out
	}

	// Buffer to store output of template execution.
	var doc bytes.Buffer

	// Add custom functions to templates function map for execution.
	var fMap = template.FuncMap{
		"uniq": func(in string) string { return name(in) },
		"port": func(in string, min, max int) string { return port(in, min, max) },
	}

	// Execute the parsed template.
	err = tmpl.Execute(&doc, fMap)
	if err != nil {
		return
	}

	// Get the output from the template execution.
	output = doc.String()
	return
}

// randUniq generates a random string from the input, consisting of uppercase and lowercase characters.
func randUniq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// randPort generates a random number between the minimum and maximum values.
func randPort(min int, max int) string {
	return strconv.Itoa(rand.Intn(max-min) + min)
}
