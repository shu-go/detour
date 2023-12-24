package main

import "github.com/shu-go/nmfmt"

type rule struct {
	Name string `json:"name,omitempty"`
	Old  string `json:"old,omitempty"`
	New  string `json:"new,omitempty"`
}

func (r rule) String() string {
	if r.Name != "" {
		return r.Name
	}
	return nmfmt.Sprintf("$Old->$New", "Old", r.Old, "New", r.New)
}
