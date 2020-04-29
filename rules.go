package main

import (
	"errors"
	"strings"

    "github.com/dlclark/regexp2"
)

const Separator = `:`

type Rule struct {
	Old, New string
}

func (r *Rule) Scan(s string) error {
    re := regexp2.MustCompile(`(?<!\\)`+Separator, 0)
    m, err := re.FindStringMatch(s)
    if err != nil {
        return err
    }

	if m == nil {
		return errors.New(s +" ... `OLD" + Separator + "NEW` required")
    }

    pos := m.Capture.Index
	if pos == 0 {
		return errors.New("OLD must not be empty")
	}

    re2 := regexp2.MustCompile(`(?<!\\)\\`+Separator, 0)
    t, err := re2.Replace(s[0:pos], Separator, 0, -1)
    if err != nil {
        return err
    }
	r.Old = t

    t, err = re2.Replace(s[pos+len(Separator):], Separator, 0, -1)
    if err != nil {
        return err
    }
	r.New = t

	return nil
}

type Rules []Rule

func (rr *Rules) Parse(s string) error {
    // comment
    if strings.HasPrefix(s, "#") {
        return nil
    }

	var r Rule
	if err := r.Scan(s); err != nil {
		return err
	}

	*rr = append(*rr, r)

	return nil
}
