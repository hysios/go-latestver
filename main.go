package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	reDate   = regexp.MustCompile(`(?m)Date:\s+(.*?)$`)
	reCommit = regexp.MustCompile(`(?m)commit\s+(.*?)\s`)
)

type Git string

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("missing git repos argument.\n")
		os.Exit(1)
	}

	repos := os.Args[1:]
	for _, repo := range repos {
		func(repo string) {
			git, err := clone(repo)
			if err != nil {
				fmt.Printf("git error %s\n", err)
			}
			defer git.Close()
			t, err := git.Time()
			if err != nil {
				fmt.Printf("time error %s\n", err)
			}
			co, err := git.Commit()
			if err != nil {
				fmt.Printf("commit error %s\n", err)
			}

			fmt.Printf("v0.0.0-%s-%s\n", t.Format("20060102150405"), co[:12])
		}(repo)
	}
}

func clone(uri string) (Git, error) {
	tmpDir, err := ioutil.TempDir("", "latestversion")
	if err != nil {
		return "", err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	pats := strings.Split(u.Path, "/")
	u.Path = strings.Join(pats[:3], "/")
	gitclone(u.String(), tmpDir)
	return Git(tmpDir), nil
}

func gitclone(url string, to string) {
	exec.Command("git", "clone", url, to, "--depth", "1").Run()
}

func (g Git) Log() string {
	cmd := exec.Command("git", "log", "-1", "--date", "iso")
	cmd.Dir = string(g)
	out, _ := cmd.Output()
	return string(out)
}

func (g Git) Time() (time.Time, error) {
	content := g.Log()
	allIndexes := reDate.FindAllSubmatchIndex([]byte(content), -1)
	if len(allIndexes) == 0 {
		return time.Time{}, errors.New("invalid time format")
	}
	loc := allIndexes[0]
	return time.Parse("2006-01-02 15:04:05 -0700", content[loc[2]:loc[3]])
}

func (g Git) Commit() (string, error) {
	content := g.Log()
	allIndexes := reCommit.FindAllSubmatchIndex([]byte(content), -1)
	if len(allIndexes) == 0 {
		return "", errors.New("invalid commit format")
	}
	loc := allIndexes[0]

	return content[loc[2]:loc[3]], nil
}

func (g Git) Close() {
	os.RemoveAll(string(g))
}
