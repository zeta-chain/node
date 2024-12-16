// Set the global default cosmos sdk config on import
//
// this should ONLY be imported in test and internal packages
// as we do not want to set the defaults for external importers
package sdkconfigdefault

import "github.com/zeta-chain/node/pkg/sdkconfig"

func init() {
	sdkconfig.SetDefault(true)
}
