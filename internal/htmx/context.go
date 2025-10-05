package htmx

import "github.com/labstack/echo/v4"

// Context encodes the values and existence of the headers from HTMX as
// described here: https://htmx.org/docs/#request-headers
type Context struct {
	// CurrentURL the current URL of the browser
	CurrentURL string
	// Prompt the user response to an hx-prompt
	Prompt string
	// Target the id of the target element if it exists
	Target string
	// TriggerName the name of the triggered element if it exists
	TriggerName string
	// Trigger the id of the triggered element if it exists
	Trigger string
	// Boosted indicates that the request is via an element using hx-boost
	Boosted bool
	// HistoryRestoreRequest “true” if the request is for history restoration after a miss in the local history cache
	HistoryRestoreRequest bool
	// Request always “true” during HTMX requests
	Request bool
}

func ContextFromRequestContext(ctx echo.Context) Context {
	reqHeaders := ctx.Request().Header

	hxBoosted := reqHeaders.Get(ReqHXBoosted) != ""
	hxCurrentURL := reqHeaders.Get(ReqHXCurrentURL)
	hxHistoryRestoreRequest := reqHeaders.Get(ReqHXHistoryRestoreRequest) == "true"
	hxPrompt := reqHeaders.Get(ReqHXPrompt)
	hxRequest := reqHeaders.Get(ReqHXRequest) == "true"
	hxTarget := reqHeaders.Get(ReqHXTarget)
	hxTriggerName := reqHeaders.Get(ReqHXTriggerName)
	hxTrigger := reqHeaders.Get(ReqHXTrigger)

	return Context{
		CurrentURL:            hxCurrentURL,
		Prompt:                hxPrompt,
		Target:                hxTarget,
		TriggerName:           hxTriggerName,
		Trigger:               hxTrigger,
		Boosted:               hxBoosted,
		HistoryRestoreRequest: hxHistoryRestoreRequest,
		Request:               hxRequest,
	}
}
