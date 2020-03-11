package models

import "github.com/pkg/errors"

type Snippet struct {
	Name        string
	Category    SnippetCategory
	UsageID     string
	UsageName   string
	Description string
	Code        string
}

type SnippetResource struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	UsageID     string `json:"usage_id"`
	UsageName   string `json:"usage_name"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

type SnippetListResource []SnippetResource

type SnippetCategory int

const (
	TagSnippet SnippetCategory = iota + 1
	AdapterSnippet
)

func (s SnippetResource) ToSnippet() (Snippet, error) {
	c, ok := StringToSnippetCategory(s.Category)
	if !ok {
		return Snippet{}, errors.New("Can't convert snippet resource category\n")
	}

	return Snippet{
		Name:        s.Name,
		Category:    c,
		UsageID:     s.UsageID,
		UsageName:   s.UsageName,
		Description: s.Description,
		Code:        s.Code,
	}, nil
}

func StringToSnippetCategory(s string) (SnippetCategory, bool) {
	switch s {
	case TagSnippet.String():
		return TagSnippet, true
	case AdapterSnippet.String():
		return AdapterSnippet, true
	}
	return 0, false
}

func (snippetCategory SnippetCategory) String() string {
	names := [...]string{
		"unknown",
		"tag",
		"adapter"}

	if snippetCategory < TagSnippet || snippetCategory > AdapterSnippet {
		return names[0]
	}
	return names[snippetCategory]
}
