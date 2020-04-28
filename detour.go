package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/shu-go/gli"
	"github.com/shu-go/shortcut"
)

type globalCmd struct {
	Rules   Rules  `cli:"rule,r=OLD:NEW"`
	RuleSet string `cli:"rule-set=FILENAME"`
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
			var r Rule
			if err := r.Scan(line); err != nil {
				return err
			}

			c.Rules = append(c.Rules, r)
		}
	}

	for _, arg := range args {
		s, err := shortcut.Open(arg)
		if err != nil {
			return err
		}

		for _, r := range c.Rules {
			log.Print(r.Old, "=>", r.New)
			re := regexp.MustCompile("(?i)" + r.Old)
			s.TargetPath = re.ReplaceAllString(s.TargetPath, r.New)
			s.WorkingDirectory = re.ReplaceAllString(s.WorkingDirectory, r.New)
		}

		s.Save(arg)
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
