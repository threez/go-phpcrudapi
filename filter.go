package phpcrudapi

import (
	"net/url"
)

type Filter interface {
	Query() string
}

type NoFilter struct{}

func (*NoFilter) Query() string { return "" }

type StringFilter struct {
	data url.Values
}

func NewStringFilter() *StringFilter {
	return &StringFilter{data: make(url.Values)}
}

func (f *StringFilter) Query() string {
	return f.data.Encode()
}

func IncludesFilter(names ...string) *StringFilter {
	f := NewStringFilter()
	f.data["include"] = names
	return f
}
