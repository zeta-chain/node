package runner

import (
	"fmt"
	"runtime"
	"time"

	"github.com/zeta-chain/node/cmd/zetacored/config"
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
		if err := r.runTestWithProtocolBalanceCheck(e2eTest); err != nil {
			return err
		}
	}
	return nil
}

func (r *E2ERunner) RunE2ETestsNoError(e2eTests []E2ETest) (err error) {
	zetacoredVersion := r.GetZetacoredVersion()
	for _, e2eTest := range e2eTests {
		if !utils.MinimumVersionCheck(e2eTest.MinimumVersion, zetacoredVersion) {
			r.Logger.Print("⚠️ skipping test - %s (minimum version %s)", e2eTest.Name, e2eTest.MinimumVersion)
			continue
		}
		if err := r.Ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled: %w", err)
		}
		if err := r.RunE2ETest(e2eTest); err != nil {
			r.Logger.Print("test %s failed: %s", e2eTest.Name, err.Error())
		}
	}
	return nil
}

func (r *E2ERunner) runTestWithProtocolBalanceCheck(e2eTest E2ETest) error {
	balancesBefore := r.checkProtocolAddressBalance(config.BaseDenom)

	if err := r.RunE2ETest(e2eTest); err != nil {
		return err
	}

	balancesAfter := r.checkProtocolAddressBalance(config.BaseDenom)
	if !balancesAfter.Equal(balancesBefore) {
		r.Logger.Print("⚠️ protocol address balance changed during test %s: before %s, after %s",
			e2eTest.Name, balancesBefore.String(), balancesAfter.String())
	}

	return nil
}

// RunE2ETest runs a e2e test
func (r *E2ERunner) RunE2ETest(e2eTest E2ETest) error {
	// wait for all dependencies to complete
	// this is only used by Bitcoin RBF test at the moment
	if len(e2eTest.Dependencies) > 0 {
		r.Logger.Print("⏳ waiting   - %s", e2eTest.Name)
		for _, dependency := range e2eTest.Dependencies {
			dependency.Wait()
		}
	}

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
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				r.Logger.Print("panic: %+v. Stack: %s", recoverVal, string(buf[:n]))

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
		testErr := r.RunE2ETest(test)
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
