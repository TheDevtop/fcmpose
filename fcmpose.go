package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	magicMajor = `\$[a-zA-Z0-9\/\.]+\$`
	magicMinor = `[\/a-zA-Z0-9\.]+`
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

func Errorln(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

// Like os.ReadFile(), but returns string
func readFile(path string) (string, error) {
	if buf, err := os.ReadFile(path); err != nil {
		return "", err
	} else {
		return string(buf), nil
	}
}

func main() {
	var (
		fstr string
		err  error
	)
	if len(os.Args) < 2 {
		fmt.Println("usage: fcmpose [FILE]")
		os.Exit(0)
	}
	if fstr, err = readFile(os.Args[1]); err != nil {
		Errorln(err)
	}
	fmt.Print(compose(fstr))
	os.Exit(0)
}
