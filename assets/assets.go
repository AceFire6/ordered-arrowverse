package assets

import "embed"

//go:embed css js templates favicon.png
var AssetFiles embed.FS
