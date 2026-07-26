package handlers

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/AceFire6/ordered-arrowverse/assets"
	"github.com/AceFire6/ordered-arrowverse/internal/customctx"
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

// stripJinja removes Jinja2 `{% ... %}` directive blocks from the input.
// It does not resolve Jinja expressions; the pre-rewrite HTML only
// references `{{ static_url('favicon.png') }}` for an icon that is
// already linked from the application layout, so we drop those too.
func stripJinja(in string) string {
	var out strings.Builder
	out.Grow(len(in))

	i := 0
	for i < len(in) {
		if i+1 < len(in) && in[i] == '{' && in[i+1] == '%' {
			end := strings.Index(in[i:], "%}")
			if end == -1 {
				// Malformed - emit remainder as-is.
				out.WriteString(in[i:])

				return out.String()
			}
			i += end + 2
			continue
		}
		out.WriteByte(in[i])
		i++
	}

	outStr := out.String()
	outStr = strings.ReplaceAll(outStr, "{{ static_url('favicon.png') }}", "")
	outStr = strings.ReplaceAll(outStr, `{{ static_url("favicon.png") }}`, "")

	return outStr
}
