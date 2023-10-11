package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	magicMajor = `\$[a-zA-Z0-9\/\.]+\$`
	magicMinor = `[\/a-zA-Z0-9\.]+`
	perm       = 0622
)

var (
	rexMajor = regexp.MustCompile(magicMajor)
	rexMinor = regexp.MustCompile(magicMinor)
)

// Warning: Recursive function
func compose(str string) string {
	list := rexMajor.FindAllString(str, -1)
	if len(list) < 1 {
		return str
	}
	for _, match := range list {
		var (
			sts string
			err error
		)
		if sts, err = readFile(rexMinor.FindString(match)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		str = strings.ReplaceAll(str, match, compose(sts))
	}
	return str
}

// Like os.ReadFile(), but returns string
func readFile(path string) (string, error) {
	if buf, err := os.ReadFile(path); err != nil {
		return "", err
	} else {
		return string(buf), nil
	}
}

func usage() {
	fmt.Println("fcmpose: Compose files with macros")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	var (
		inf  = flag.String("if", "", "Specify input file")
		outf = flag.String("of", "/dev/stdout", "Specify output file")
		fstr string
		err  error
	)

	flag.Usage = usage
	flag.Parse()

	if fstr, err = readFile(*inf); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.WriteFile(*outf, []byte(compose(fstr)), perm)
	os.Exit(0)
}
