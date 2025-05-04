package runner

import (
	"fmt"
	"runtime"
	"time"

	"github.com/zeta-chain/node/e2e/utils"
)

// RunE2ETests runs a list of e2e tests
func (r *E2ERunner) RunE2ETests(e2eTests []E2ETest) (err error) {
	zetacoredVersion := r.GetZetacoredVersion()
	for _, e2eTest := range e2eTests {
		if !utils.MinimumVersionCheck(e2eTest.MinimumVersion, zetacoredVersion) {
			r.Logger.Print("⚠️ skipping test - %s (minimum version %s)", e2eTest.Name, e2eTest.MinimumVersion)
			continue
		}
		if err := r.Ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled: %w", err)
		}
		if err := r.RunE2ETest(e2eTest, true); err != nil {
			return err
		}
	}
	return nil
}

// RunE2ETest runs a e2e test
func (r *E2ERunner) RunE2ETest(e2eTest E2ETest, checkAccounting bool) error {
	startTime := time.Now()
	// note: spacing is padded to width of completed message
	r.Logger.Print("⏳ running   - %s", e2eTest.Name)

	// run e2e test, if args are not provided, use default args
	args := e2eTest.Args
	if len(args) == 0 {
		args = e2eTest.DefaultArgs()
	}
	errChan := make(chan error)
	go func() {
		defer func() {
			if recoverVal := recover(); recoverVal != nil {
				switch recoverVal.(type) {
				case runtime.Error:
					// this is a probably a nil dereference or divide by zero which we would want to log
					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					r.Logger.Info("panic: %s", string(buf[:n]))
				}
				errChan <- fmt.Errorf("panic: %v", recoverVal)
			}
			close(errChan)
		}()
		e2eTest.E2ETest(r, args)
		errChan <- nil
	}()

	select {
	case <-r.Ctx.Done():
		return fmt.Errorf("context cancelled in %s after %s", e2eTest.Name, time.Since(startTime))
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("%s failed (duration %s): %w", e2eTest.Name, time.Since(startTime), err)
		}
	}

	// check zrc20 balance vs. supply
	if checkAccounting {
		r.CheckZRC20BalanceAndSupply()
	}

	r.Logger.Print("✅ completed - %s (%s)", e2eTest.Name, time.Since(startTime))

	return nil
}

// RunE2ETestsIntoReport runs a list of e2e tests by name in a list of e2e tests and returns a report
// The function doesn't return an error, it returns a report with the error
func (r *E2ERunner) RunE2ETestsIntoReport(e2eTests []E2ETest) (TestReports, error) {
	// go through all tests
	reports := make(TestReports, 0, len(e2eTests))
	for _, test := range e2eTests {
		// get info before test
		balancesBefore, err := r.GetAccountBalances(true)
		if err != nil {
			return nil, err
		}
		timeBefore := time.Now()

		// run test
		testErr := r.RunE2ETest(test, false)
		if testErr != nil {
			r.Logger.Print("test %s failed: %s", test.Name, testErr.Error())
		}

		// wait 5 sec to make sure we get updated balances
		time.Sleep(5 * time.Second)

		// get info after test
		balancesAfter, err := r.GetAccountBalances(true)
		if err != nil {
			return nil, err
		}
		timeAfter := time.Now()

		// create report
		report := TestReport{
			Name:     test.Name,
			Success:  testErr == nil,
			Time:     timeAfter.Sub(timeBefore),
			GasSpent: GetAccountBalancesDiff(balancesBefore, balancesAfter),
		}
		reports = append(reports, report)
	}

	return reports, nil
}
