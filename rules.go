package main

import (
	"errors"
	"strings"
)

const Separator = `=>`

type Rule struct {
	Old, New string
}

func (r *Rule) Scan(s string) error {
	pos := strings.Index(s, Separator)
	if pos == -1 {
		return errors.New("`OLD" + Separator + "NEW` required")
	} else if pos == 0 {
		return errors.New("OLD must not be empty")
	}

	r.Old = s[0:pos]
	r.New = s[pos+len(Separator):]

	return nil
}

type Rules []Rule

func (rr *Rules) Parse(s string) error {
	var r Rule
	if err := r.Scan(s); err != nil {
		return err
	}

	*rr = append(*rr, r)

	return nil
}
