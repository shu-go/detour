package main

import (
	"regexp"

	"github.com/shu-go/nmfmt"
)

type rule struct {
	Name string `json:"name,omitempty" yaml:",omitempty"`
	Old  string `json:"old,omitempty"`
	New  string `json:"new,omitempty" yaml:",omitempty"`

	Regexp bool `json:"regexp,omitempty" yaml:",omitempty"`
}

func (r rule) String() string {
	if r.Name != "" {
		return r.Name
	}
	re := ""
	if r.Regexp {
		re = "(regexp) "
	}
	return nmfmt.Sprintf("$Regexp$Old -> $New", "Old", r.Old, "New", r.New, "Regexp", re)
}

func (r rule) Apply(s string) (string, bool) {
	var re *regexp.Regexp
	if r.Regexp {
		re = regexp.MustCompile("(?i)" + r.Old)
	} else {
		re = regexp.MustCompile("(?i)" + regexp.QuoteMeta(r.Old))
	}

	after := re.ReplaceAllString(s, r.New)
	return after, s != after
}
