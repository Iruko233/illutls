package profiles

// This package registers all built-in browser profiles via init() functions.
// Import it with a blank identifier to load all profiles:
//
//	import _ "github.com/Iruko233/illutls/profiles"
import (
	"github.com/Iruko233/illutls"
)

// Re-export for use by individual profile files.
var register = illutls.RegisterProfile
