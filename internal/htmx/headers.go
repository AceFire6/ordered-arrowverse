package htmx

const (
	// ReqHXBoosted indicates that the request is via an element using hx-boost
	ReqHXBoosted = "HX-Boosted"
	// ReqHXCurrentURL the current URL of the browser
	ReqHXCurrentURL = "HX-Current-URL"
	// ReqHXHistoryRestoreRequest “true” if the request is for history restoration after a miss in the local history cache
	ReqHXHistoryRestoreRequest = "HX-History-Restore-Request"
	// ReqHXPrompt the user response to an hx-prompt
	ReqHXPrompt = "HX-Prompt"
	// ReqHXRequest always “true”
	ReqHXRequest = "HX-Request"
	// ReqHXTarget the id of the target element if it exists
	ReqHXTarget = "HX-Target"
	// ReqHXTriggerName the name of the triggered element if it exists
	ReqHXTriggerName = "HX-Trigger-Name"
	// ReqHXTrigger the id of the triggered element if it exists
	ReqHXTrigger = "HX-Trigger"

	// RespHXLocation allows you to do a client-side redirect that does not do a full page reload
	RespHXLocation = "HX-Location"
	// RespHXPushURL pushes a new url into the history stack
	RespHXPushURL = "HX-Push-Url"
	// RespHXRedirect can be used to do a client-side redirect to a new location
	RespHXRedirect = "HX-Redirect"
	// RespHXRefresh if set to “true” the client-side will do a full refresh of the page
	RespHXRefresh = "HX-Refresh"
	// RespHXReplaceURL replaces the current URL in the location bar
	RespHXReplaceURL = "HX-Replace-Url"
	// RespHXReswap allows you to specify how the response will be swapped. See hx-swap for possible values
	RespHXReswap = "HX-Reswap"
	// RespHXRetarget a CSS selector that updates the target of the content update to a different element on the page
	RespHXRetarget = "HX-Retarget"
	// RespHXReselect a CSS selector that allows you to choose which part of the response is used to be swapped in. Overrides an existing hx-select on the triggering element = "hx-select on the triggering element"
	RespHXReselect = "HX-Reselect"
	// RespHXTrigger allows you to trigger client-side events
	RespHXTrigger = "HX-Trigger"
	// RespHXTriggerAfterSettle allows you to trigger client-side events after the settle step
	RespHXTriggerAfterSettle = "HX-Trigger-After-Settle"
	// RespHXTriggerAfterSwap allows you to trigger client-side events after the swap step
	RespHXTriggerAfterSwap = "HX-Trigger-After-Swap"
)
