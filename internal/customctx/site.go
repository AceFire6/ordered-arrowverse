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
