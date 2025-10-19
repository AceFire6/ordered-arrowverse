package frontend

import (
	"errors"
	"fmt"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

var ErrNoLayoutSet = errors.New("no layout set on Frontend")

// Frontend is the default settings for all frontend pages and helps with rendering specific pages
// and components.
type Frontend struct {
	DefaultPageConfig *PageConfig
	// Layout root layout component that the other components are rendered into
	Layout  LayoutFunc
	Title   string
	Heading string
}

type (
	LayoutFunc = func(config *PageConfig) templ.Component
)

type NewParams struct {
	// Reverser is the function we call to derive the URLs for named URLs by their names
	Reverser URLReverser
	// Heading is the text displayed next to the logo in the navbar
	Heading  string
	Title    string
	ShowList []ShowData
}

func New(frontendParams NewParams) *Frontend {
	return &Frontend{
		DefaultPageConfig: &PageConfig{
			Title:       frontendParams.Title,
			Heading:     frontendParams.Heading,
			Reverser:    frontendParams.Reverser,
			Contents:    nil,
			JsScripts:   []string{},
			Stylesheets: []string{},
			ShowList:    frontendParams.ShowList,
		},
		Title:   frontendParams.Title,
		Heading: frontendParams.Heading,
		Layout:  nil,
	}
}

func (frontend *Frontend) WithTitle(title string) *Frontend {
	frontend.DefaultPageConfig.Title = title

	return frontend
}

func (frontend *Frontend) WithLayout(layout LayoutFunc) *Frontend {
	frontend.Layout = layout

	return frontend
}

func (frontend *Frontend) RenderPage(ctx echo.Context, statusCode int, pageConfig *PageConfig) error {
	if frontend.Layout == nil {
		return ErrNoLayoutSet
	}

	renderPageConfig := MergePageConfigs(frontend.DefaultPageConfig, pageConfig, MergeAll)
	layoutWithComponent := frontend.Layout(renderPageConfig)

	err := Render(ctx, statusCode, layoutWithComponent)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	return nil
}
