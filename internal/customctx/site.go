package customctx

import "github.com/AceFire6/ordered-arrowverse/internal/frontend"

// ApplySiteConfig copies the legacy-site signal (host match) and
// promotional values onto the supplied PageConfig. Safe to call when
// cc.Site is nil.
//
// Returns the same pointer for chainable use:
//
//	pageConfig := customctx.ApplySiteConfig(cc, frontend.NewPage(...))
func ApplySiteConfig(cc *Context, pageConfig *frontend.PageConfig) *frontend.PageConfig {
	if pageConfig == nil {
		return nil
	}
	if cc == nil || cc.Site == nil {
		return pageConfig
	}

	pageConfig.OldSiteHost = cc.Site.OldSiteHost
	pageConfig.NewSiteURL = cc.Site.NewSiteURL
	pageConfig.UsingOldSite = cc.Site.OldSiteHost != "" &&
		cc.Request().Host == cc.Site.OldSiteHost

	return pageConfig
}

// CSRFCsrfToken returns the CSRF token from the Echo context, or an
// empty string when the middleware hasn't populated one for this
// request. The token comes from Echo's CSRF middleware which seeds it
// into the context under the key "csrf".
func CSRFCsrfToken(cc *Context) string {
	if cc == nil {
		return ""
	}
	raw := cc.Get("csrf")
	if token, ok := raw.(string); ok {
		return token
	}

	return ""
}

// ApplyCSRFToken copies Echo's CSRF middleware token onto the supplied
// PageConfig so the layout can render a <meta> tag for htmx.
func ApplyCSRFToken(cc *Context, pageConfig *frontend.PageConfig) *frontend.PageConfig {
	if pageConfig == nil {
		return nil
	}
	if cc == nil {
		return pageConfig
	}

	pageConfig.CSRFToken = CSRFCsrfToken(cc)

	return pageConfig
}
