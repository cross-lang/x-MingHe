package health

import (
	"github.com/google/wire"
)

// ProviderSet Health handler provider set
var ProviderSet = wire.NewSet(NewHandler)
