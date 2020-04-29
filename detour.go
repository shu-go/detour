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
	Rules   Rules  `cli:"rule,r"`
	RuleSet string `cli:"rule-set=FILENAME"`

    Verbose bool `cli:"verbose,v"`
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

    if c.Verbose {
        fmt.Println("Rules:")
        for _, r := range c.Rules {
            fmt.Println("  "+r.Old+" => "+r.New)
        }
    }

    var files []string
    if len(args) == 0{
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

        if c.Verbose{
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

		s, err := shortcut.Open(f)
		if err != nil {
			continue // go to next
		}

        if c.Verbose{
            fmt.Println("  "+f)
        }

		for _, r := range c.Rules {
			re := regexp.MustCompile("(?i)" + r.Old)
			s.TargetPath = re.ReplaceAllString(s.TargetPath, r.New)
			s.WorkingDirectory = re.ReplaceAllString(s.WorkingDirectory, r.New)
		}

		s.Save(f)
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
	app.Desc = ""
	app.Version = Version
	app.Usage = ``
	app.Copyright = "(C) 2020 Shuhei Kubota"
	app.Run(os.Args)
}
