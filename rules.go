package main

import (
	"errors"
	"strings"

	"github.com/dlclark/regexp2"
)

const separator = `:`

type rule struct {
	Old, New string
}

func (r *rule) scan(s string) error {
	re := regexp2.MustCompile(`(?<!\\)`+separator, 0)
	m, err := re.FindStringMatch(s)
	if err != nil {
		return err
	}

	if m == nil {
		return errors.New(s + " ... `OLD" + separator + "NEW` required")
	}

	pos := m.Capture.Index
	if pos == 0 {
		return errors.New("OLD must not be empty")
	}

	re2 := regexp2.MustCompile(`(?<!\\)\\`+separator, 0)
	t, err := re2.Replace(s[0:pos], separator, 0, -1)
	if err != nil {
		return err
	}
	r.Old = t

	t, err = re2.Replace(s[pos+len(separator):], separator, 0, -1)
	if err != nil {
		return err
	}
	r.New = t

	return nil
}

type rules []rule

func (rr *rules) Parse(s string) error {
	// comment
	if strings.HasPrefix(s, "#") {
		return nil
	}

	var r rule
	if err := r.scan(s); err != nil {
		return err
	}

	*rr = append(*rr, r)

	return nil
}
