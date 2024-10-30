package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"

	//"path/filepath"
	"os"
	"strings"

	"encoding/json"

	"github.com/mattn/go-zglob"
	"gopkg.in/yaml.v2"

	"github.com/shu-go/gli/v2"
	"github.com/shu-go/shortcut"
)

type globalCmd struct {
	RuleSet string `cli:"rule-set=JSON_OR_YAML_FILENAME"`

	Verbose bool `cli:"verbose,v"`
	DryRun  bool `cli:"dry-run,n"`

	GenCmd genCmd `cli:"generate,gen" help:"generate a example file"`
}

func (c globalCmd) Run(args []string) error {
	if c.RuleSet == "" {
		return errors.New("option --rule-set is required")
	}

	rules := []rule{}

	if c.RuleSet != "" {
		isYAML := false
		ext := strings.ToLower(filepath.Ext(c.RuleSet))
		if in(ext, ".yaml", ".yml") {
			isYAML = true
		}

		data, err := os.ReadFile(c.RuleSet)
		if err != nil {
			return err
		}

		cover := struct {
			Rules []rule `json:"rules"`
		}{}

		if isYAML {
			err = yaml.Unmarshal(data, &rules)
			if err != nil {
				err = yaml.Unmarshal(data, &cover)
				if err != nil {
					return err
				}
				rules = cover.Rules
			}
		} else {
			err = json.Unmarshal(data, &rules)
			if err != nil {
				err = json.Unmarshal(data, &cover)
				if err != nil {
					return err
				}
				rules = cover.Rules
			}
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
	_ any `usage:"detour generate {JSON_OR_YAML_FILENAME}"`
}

func (c genCmd) Run(args []string) error {
	if len(args) != 1 {
		return errors.New("FILENAME is required")
	}

	filename := args[0]

	isYAML := false
	ext := strings.ToLower(filepath.Ext(filename))
	if in(ext, ".yaml", ".yml") {
		isYAML = true
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

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	if isYAML {
		data, err := yaml.Marshal(cover)
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		if err != nil {
			return err
		}
	} else {
		enc := json.NewEncoder(f)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		err = enc.Encode(cover)
		if err != nil {
			return err
		}
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func in(s string, elems ...string) bool {
	result := false

	for _, e := range elems {
		if strings.EqualFold(s, e) {
			result = true
			break
		}
	}

	return result
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
