package handlers

import (
	"io/fs"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"

	"github.com/AceFire6/ordered-arrowverse/assets"
	"github.com/AceFire6/ordered-arrowverse/internal/customctx"
)

var (
	jinjaDirective = regexp.MustCompile(`\{%[\s\S]*?%\}`)
	jinjaExpr      = regexp.MustCompile(`\{\{[\s\S]*?\}\}`)
)

// legalFS is the read-only handle to the embedded templates directory.
// Routes call Handlers on the sub-directory containing the Termly-exported
// policy HTML files.
var legalFS = mustSubFS(assets.AssetFiles, "templates")

// mustSubFS panics if the embedded sub-filesystem can't be opened. This
// only fails at startup, so a panic is appropriate.
func mustSubFS(root fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(root, dir)
	if err != nil {
		panic("assets embed: missing directory: " + dir + ": " + err.Error())
	}

	return sub
}

// LegalDocument reads a Termly-exported HTML file from the embedded
// assets, strips the Jinja directive blocks that survived the Quart → Go
// rewrite, and serves it as a plain HTML response. The shared layout
// (navbar, footer buttons) is intentionally not applied to keep these
// pages compatible with the structure Termly issued.
//
// `name` is the file stem without the `_policy.html` suffix.
func LegalDocument(name string) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*customctx.Context)

		fileName := name + "_policy.html"
		raw, err := fs.ReadFile(legalFS, fileName)
		if err != nil {
			cc.Log.Err(err).Str("file_name", fileName).Msg("could not read legal asset")
			return echo.NewHTTPError(http.StatusNotFound, "legal document unavailable")
		}

		body := stripJinja(string(raw))

		return c.Blob(http.StatusOK, "text/html; charset=utf-8", []byte(body))
	}
}

// stripJinja removes every Jinja2 directive block (`{% ... %}`) and
// expression (`{{ ... }}`) from the input. Block bodies *and* the
// directives themselves are dropped, which is what the legal pages need:
// the body text was already rendered by Termly; the Jinja scaffolding
// only set up template inheritance that no longer applies.
func stripJinja(in string) string {
	cleaned := jinjaDirective.ReplaceAllString(in, "")
	cleaned = jinjaExpr.ReplaceAllString(cleaned, "")

	return cleaned
}
