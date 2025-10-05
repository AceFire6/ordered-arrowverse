package echojson

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"

	"github.com/labstack/echo/v4"
)

// JSONSerializer This is an adaptation of the DefaultJSONSerializer from echo
type JSONSerializer struct{}

func (f JSONSerializer) Serialize(c echo.Context, input interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}

	return enc.Encode(input) //nolint:wrapcheck // no need to wrap this library error
}

func (f JSONSerializer) Deserialize(c echo.Context, input interface{}) error {
	err := json.NewDecoder(c.Request().Body).Decode(input)

	var ute *json.UnmarshalTypeError
	var se *json.SyntaxError

	if errors.As(err, &ute) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if errors.As(err, &se) {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}

	return err //nolint:wrapcheck // no need to wrap this library error
}
