package frontend

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Render This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
//
// From: https://github.com/a-h/templ/blob/38d9ecaa826269d8227bf7b7effe0b9703ecdb9a/examples/integration-echo/main.go#L16-L26
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return fmt.Errorf("templ render: %w", err)
	}

	err := ctx.HTMLBlob(statusCode, buf.Bytes())
	if err != nil {
		return fmt.Errorf("rendered templ HTML response: %w", err)
	}

	return nil
}
