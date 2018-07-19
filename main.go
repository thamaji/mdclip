package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/russross/blackfriday.v2"
)

func usage() {
	output := flag.CommandLine.Output()
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Usage: "+os.Args[0]+" [OPTIONS] FILE [FILE...]")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Render markdown to clipboard")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Options:")
	flag.CommandLine.PrintDefaults()
}

func main() {
	flag.Usage = usage
	output := flag.CommandLine.Output()

	var version, help bool

	flag.BoolVar(&version, "v", false, "show version")
	flag.BoolVar(&help, "h", false, "show help")
	flag.Parse()

	if help {
		usage()
		return
	}

	if version {
		fmt.Fprintln(output, "v1.0.0")
		return
	}

	args := flag.Args()
	if len(args) <= 0 {
		args = append(args, "-")
	}

	cmd := exec.Command("xclip", "-t", "text/html", "-selection", "clipboard", "-i")
	w, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, arg := range args {
		html, err := RenderMarkdown(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Fprintln(w, string(html))
	}

	w.Close()

	if err := cmd.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func RenderMarkdown(filename string) ([]byte, error) {
	var file io.ReadCloser
	var err error

	if filename == "-" {
		file = os.Stdin

	} else {
		file, err = os.Open(filename)
		if err != nil {
			return nil, err
		}
	}

	bytes, err := ioutil.ReadAll(file)
	file.Close()
	if err != nil {
		return nil, err
	}

	return blackfriday.Run(bytes), nil
}
