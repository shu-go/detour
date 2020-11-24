package main

import (
	"bufio"
	"fmt"

	//"path/filepath"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mattn/go-zglob"

	"github.com/shu-go/gli"
	"github.com/shu-go/shortcut"
)

type globalCmd struct {
	Rules   rules  `cli:"rule,r"`
	RuleSet string `cli:"rule-set=FILENAME"`

	Verbose bool `cli:"verbose,v"`
	DryRun  bool `cli:"dry-run,n"`
}

func (c globalCmd) Run(args []string) error {
	if c.RuleSet != "" {
		f, err := os.Open(c.RuleSet)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if err := c.Rules.Parse(line); err != nil {
				return nil
			}
		}
	}

	if c.DryRun {
		fmt.Println("")
		fmt.Println("DRY RUN MODE")
		fmt.Println("")
	}

	if c.Verbose {
		fmt.Println("Rules:")
		for _, r := range c.Rules {
			fmt.Println("  " + r.Old + " => " + r.New)
		}
	}

	var files []string
	if len(args) == 0 {
		args = append(args, "*.lnk")
	}
	for _, arg := range args {
		if !strings.HasSuffix(arg, ".lnk") {
			arg += ".lnk"
		}
		//ff, err := filepath.Glob(arg)
		ff, err := zglob.Glob(arg)
		if err != nil {
			return err
		}

		files = append(files, ff...)
	}

	if c.Verbose {
		fmt.Println("Files:")
	}
	for _, f := range files {
		fi, err := os.Lstat(f)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			continue
		}

		changed := false

		s, err := shortcut.Open(f)
		if err != nil {
			continue // go to next
		}

		for _, r := range c.Rules {
			re := regexp.MustCompile("(?i)" + r.Old)

			after := re.ReplaceAllString(s.TargetPath, r.New)
			if s.TargetPath != after {
				changed = true
				s.TargetPath = after
			}
			after = re.ReplaceAllString(s.WorkingDirectory, r.New)
			if s.WorkingDirectory != after {
				changed = true
				s.WorkingDirectory = after
			}
		}

		if c.Verbose {
			if changed {
				fmt.Println("* " + f)
			} else {
				fmt.Println("  " + f)
			}
		}

		if changed && !c.DryRun {
			s.Save(f)
		}
	}

	return nil
}

// Version is app version
var Version string

func init() {
	if Version == "" {
		Version = "dev-" + time.Now().Format("20060102")
	}
}

func main() {
	app := gli.NewWith(&globalCmd{})
	app.Name = "detour"
	app.Desc = "Windows shortcut replacer tool"
	app.Version = Version
	app.Usage = `detour -r old1:new1 -r old2:new2
detour -r old1:new1 -r old2:new2  ./subdir/*
detour --rule-set your_rules.txt`
	app.Copyright = "(C) 2020 Shuhei Kubota"
	app.Run(os.Args)
}
