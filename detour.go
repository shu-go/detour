package main

import (
	"errors"
	"fmt"
	"slices"

	//"path/filepath"
	"os"
	"strings"

	"encoding/json"

	"github.com/mattn/go-zglob"

	"github.com/shu-go/gli/v2"
	"github.com/shu-go/shortcut"
)

type globalCmd struct {
	RuleSet string `cli:"rule-set=JSON_FILENAME"`

	Verbose bool `cli:"verbose,v"`
	DryRun  bool `cli:"dry-run,n"`

	GenCmd genCmd `cli:"generate,gen" help:"generate a example file"`
}

func (c globalCmd) Run(args []string) error {
	if c.RuleSet == "" {
		return errors.New("option rule-set is required")
	}

	rules := []rule{}

	if c.RuleSet != "" {
		data, err := os.ReadFile(c.RuleSet)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &rules)
		if err != nil {
			cover := struct {
				Rules []rule `json:"rules"`
			}{}
			err = json.Unmarshal(data, &cover)
			if err != nil {
				return err
			}
			rules = cover.Rules
		}
	}

	rules = slices.DeleteFunc(rules, func(r rule) bool {
		return strings.TrimSpace(r.Old) == ""
	})

	if c.DryRun {
		fmt.Println("")
		fmt.Println("DRY RUN MODE")
		fmt.Println("")
	}

	if c.Verbose {
		fmt.Println("Rules:")
		for _, r := range rules {
			fmt.Println("  " + r.String())
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

		for _, r := range rules {
			after, tmpchanged := r.Apply(s.TargetPath)
			if tmpchanged {
				s.TargetPath = after
				changed = true
			}

			after, tmpchanged = r.Apply(s.WorkingDirectory)
			if tmpchanged {
				s.WorkingDirectory = after
				changed = true
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

type genCmd struct {
	_ any `usage:"detour generate {JSON_FILENAME}"`
}

func (c genCmd) Run(args []string) error {
	if len(args) != 1 {
		return errors.New("JSON_FILENAME is required")
	}

	cover := struct {
		Rules []rule
	}{
		Rules: []rule{
			{Name: "C: -> D:", Old: "C:", New: "D:"},
			{Name: "detour -> shortcut", Old: "detour", New: "shortcut"},
			{Name: "", Old: `\bg.`, New: "go", Regexp: true},
		},
	}

	f, err := os.Create(args[0])
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(cover)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// Version is app version
var Version string

func main() {
	app := gli.NewWith(&globalCmd{})
	app.Name = "detour"
	app.Desc = "Windows shortcut replacer tool"
	app.Version = Version
	app.Usage = `detour -v --rule-set myrules.json ./subdir/**/*`
	app.Copyright = "(C) 2020 Shuhei Kubota"
	app.Run(os.Args)

}
