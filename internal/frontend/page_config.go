package frontend

import (
	"slices"

	"github.com/a-h/templ"
)

// PageConfig The settings used to render a given page
type PageConfig struct {
	// Title is the title of the page in the browser
	Title string
	// Heading is the text displayed next to the logo in the navbar
	Heading     string
	Reverser    URLReverser
	Contents    templ.Component
	JsScripts   []string
	Stylesheets []string
	// Site data
	// This is the list shown in the filters and used to generate the acknowledgements on the layouts page
	ShowList []ShowData
}

type URLReverser interface {
	Reverse(name string, params ...interface{}) string
}

// NewPage Create a new page - only accept the required arguments as inputs
func NewPage(contents templ.Component) *PageConfig {
	// set up the defaults first
	pageConfig := &PageConfig{
		Title:       "",
		Heading:     "",
		Contents:    contents,
		Reverser:    nil,
		JsScripts:   []string{},
		Stylesheets: []string{},
	}

	return pageConfig
}

func (page *PageConfig) Reverse(name string, params ...interface{}) string {
	return page.Reverser.Reverse(name, params...)
}

func (page *PageConfig) WithTitle(title string) *PageConfig {
	page.Title = title

	return page
}

func (page *PageConfig) WithHeading(heading string) *PageConfig {
	page.Heading = heading

	return page
}

func (page *PageConfig) WithContents(contents templ.Component) *PageConfig {
	page.Contents = contents

	return page
}

func (page *PageConfig) WithJsScripts(jsScripts []string) *PageConfig {
	page.JsScripts = slices.Concat(page.JsScripts, jsScripts)

	return page
}

func (page *PageConfig) WithStylesheets(stylesheets []string) *PageConfig {
	page.Stylesheets = slices.Concat(page.Stylesheets, stylesheets)

	return page
}

func (page *PageConfig) WithReverser(reverser URLReverser) *PageConfig {
	page.Reverser = reverser

	return page
}

func (page *PageConfig) Copy() *PageConfig {
	pageCopy := &PageConfig{
		Title:       page.Title,
		Heading:     page.Heading,
		Reverser:    nil,
		Contents:    nil,
		JsScripts:   nil,
		Stylesheets: nil,
	}

	if page.Reverser != nil {
		reverseFuncCopy := &page.Reverser
		pageCopy.Reverser = *reverseFuncCopy
	}

	if page.Contents != nil {
		contentsCopy := &page.Contents
		pageCopy.Contents = *contentsCopy
	}

	pageCopy.JsScripts = append(pageCopy.JsScripts, page.JsScripts...)
	pageCopy.Stylesheets = append(pageCopy.Stylesheets, page.Stylesheets...)

	return pageCopy
}

type MergeOption = int

const (
	MergeJS MergeOption = 1 << iota
	MergeCSS
	OverwriteJS
	OverwriteCSS
	OverwriteAll = OverwriteJS | OverwriteCSS
	MergeAll     = MergeJS | MergeCSS
)

// MergePageConfigs takes in page configs A and B and creates a new PageConfig that combines
// A and B using the strategy set by option the MergeOption input
func MergePageConfigs(pageConfigA, pageConfigB *PageConfig, option MergeOption) *PageConfig {
	newPageConfig := pageConfigA.Copy()

	if pageConfigB.Title != "" {
		newPageConfig.Title = pageConfigB.Title
	}

	if pageConfigB.Heading != "" {
		newPageConfig.Heading = pageConfigB.Heading
	}

	newPageConfig.Contents = pageConfigB.Contents

	if pageConfigB.Reverser != nil {
		newPageConfig.Reverser = pageConfigB.Reverser
	}

	shouldMergeJS := option&MergeJS != 0
	shouldMergeCSS := option&MergeCSS != 0

	if shouldMergeJS {
		newPageConfig.JsScripts = append(newPageConfig.JsScripts, pageConfigB.JsScripts...)
	} else {
		newPageConfig.JsScripts = make([]string, len(pageConfigB.JsScripts))
		copy(newPageConfig.JsScripts, pageConfigB.JsScripts)
	}

	if shouldMergeCSS {
		newPageConfig.Stylesheets = append(newPageConfig.Stylesheets, pageConfigB.Stylesheets...)
	} else {
		newPageConfig.Stylesheets = make([]string, len(pageConfigB.Stylesheets))
		copy(newPageConfig.Stylesheets, pageConfigB.Stylesheets)
	}

	return newPageConfig
}
