package runner

import (
	"time"
)

// SetZEVMContractsV2 set contracts for the ZEVM
func (r *E2ERunner) SetZEVMContractsV2() {
	r.Logger.Print("⚙️ deploying system contracts and ZRC20s on ZEVM")
	startTime := time.Now()
	defer func() {
		r.Logger.Info("System contract deployments took %s\n", time.Since(startTime))
	}()
}
