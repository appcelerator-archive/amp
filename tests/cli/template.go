package cli

import (
	"bytes"
	"math/rand"
	"strconv"
	"text/template"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().Unix())
}

// create, parse and execute a template to generate unique values
func templating(input string, cache map[string]string) (output string, err error) {
	var tmpl *template.Template
	tmpl, err = template.New("Command").Parse(input)
	if err != nil {
		return
	}
	// custom function to create a unique name with a randomly generated string
	name := func(in string) string {
		if val, ok := cache[in]; ok {
			return val
		}
		out := in + randUniq(10)
		cache[in] = out
		return out
	}
	// custom function to randomly generate a port number
	port := func(in string, min, max int) string {
		if val, ok := cache[in]; ok {
			return val
		}
		out := randPort(min, max)
		cache[in] = out
		return out
	}
	var doc bytes.Buffer
	// add the custom functions to template for execution
	var fMap = template.FuncMap{
		"uniq": func(in string) string { return name(in) },
		"port": func(in string, min, max int) string { return port(in, min, max) },
	}
	// execute the parsed template
	err = tmpl.Execute(&doc, fMap)
	if err != nil {
		return
	}
	output = doc.String()
	return
}

// generate a random string consisting of uppercase and lowercase characters
func randUniq(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// generate a random number between the minimum and maximum values
func randPort(min int, max int) string {
	return strconv.Itoa(rand.Intn(max-min) + min)
}
