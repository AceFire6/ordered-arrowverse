package customctx

import (
	"encoding/json"
	"net/http"

	"github.com/AceFire6/ordered-arrowverse/internal/htmx"
)

// Trigger fires a client-side event (the HX-Trigger response header)
// by name, optionally with JSON-serialised detail. Useful after a
// successful post to refresh a stripe of the UI.
func (cc *Context) Trigger(eventName string, detail any) error {
	if cc == nil {
		return nil
	}

	if detail == nil {
		cc.Response().Header().Set(htmx.RespHXTrigger, eventName)
		return nil
	}

	encoded, err := json.Marshal(detail)
	if err != nil {
		return err //nolint:wrapcheck // handler will translate
	}

	event := eventName
	if len(encoded) > 2 { // ignore "{}"
		event = event + ":" + string(encoded)
	}
	cc.Response().Header().Set(htmx.RespHXTrigger, event)

	return nil
}

// TriggerAfterSettle fires a client-side event after the settle phase.
// Same shape as Trigger; see htmx.org/docs/#response-headers.
func (cc *Context) TriggerAfterSettle(eventName string, detail any) error {
	if cc == nil {
		return nil
	}

	if detail == nil {
		cc.Response().Header().Set(htmx.RespHXTriggerAfterSettle, eventName)

		return nil
	}

	encoded, err := json.Marshal(detail)
	if err != nil {
		return err //nolint:wrapcheck
	}

	event := eventName
	if len(encoded) > 2 {
		event = event + ":" + string(encoded)
	}
	cc.Response().Header().Set(htmx.RespHXTriggerAfterSettle, event)

	return nil
}

// HXRedirect instructs htmx to perform a client-side redirect to the
// supplied path. Falls back to a 302 echo redirect when the request was
// not made via htmx, preserving compatibility with plain browser
// navigations.
//
// Method name HXRedirect avoids shadowing the embedded
// echo.Context.Redirect(int, string) from the parent context.
func (cc *Context) HXRedirect(target string) error {
	if cc == nil {
		return nil
	}

	if cc.HTMX.Request {
		cc.Response().Header().Set(htmx.RespHXRedirect, target)

		return cc.NoContent(http.StatusNoContent)
	}

	return cc.Redirect(http.StatusFound, target) //nolint:wrapcheck
}

// PushURL pushes the supplied path into the browser history stack.
func (cc *Context) PushURL(path string) {
	if cc == nil {
		return
	}

	cc.Response().Header().Set(htmx.RespHXPushURL, path)
}

// ReplaceURL replaces the current history entry with path (instead of
// pushing a new one).
func (cc *Context) ReplaceURL(path string) {
	if cc == nil {
		return
	}

	cc.Response().Header().Set(htmx.RespHXReplaceURL, path)
}

// Location emits an HX-Location header that performs an explicit
// client-side location change without a full reload.
func (cc *Context) Location(path string) {
	if cc == nil {
		return
	}

	cc.Response().Header().Set(htmx.RespHXLocation, path)
}
