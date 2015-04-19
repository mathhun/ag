package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	re       *regexp.Regexp
	reIgnore *regexp.Regexp = regexp.MustCompile(`(^|/)\.git($|/)`)
)

func main() {
	argRe := os.Args[1]
	argDir := os.Args[2]

	re, err := regexp.Compile(argRe)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Failed to compile regexp: %s\n", argRe)
		os.Exit(1)
	}

	if _, err := os.Stat(argDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stdout, "Not exist: %s\n", argDir)
		os.Exit(1)
	}

	var list []string
	filepath.Walk(argDir, func(path string, info os.FileInfo, err error) error {
		if reIgnore.FindString(path) == "" && (info != nil && !info.IsDir()) {
			list = append(list, path)
		} else {
			//fmt.Println("ignored: ", path)
		}

		return nil
	})

	for _, file := range list {
		doGrep(re, file)
	}
}

func doGrep(re *regexp.Regexp, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		file.Close()
	}()

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if m := re.Find(line); len(m) > 0 {
			fmt.Printf("%s\t%s\n", path, string(line))
		}
	}
}
