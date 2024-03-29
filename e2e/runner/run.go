package runner

import (
	"fmt"
	"runtime"
	"time"
)

// RunE2ETests runs a list of e2e tests
func (runner *E2ERunner) RunE2ETests(e2eTests []E2ETest) (err error) {
	for _, e2eTest := range e2eTests {
		if err := runner.RunE2ETest(e2eTest, true); err != nil {
			return err
		}
	}
	return nil
}

// RunE2ETest runs a e2e test
func (runner *E2ERunner) RunE2ETest(e2eTest E2ETest, checkAccounting bool) (err error) {
	// return an error on panic
	// https://github.com/zeta-chain/node/issues/1500
	defer func() {
		if r := recover(); r != nil {
			// print stack trace
			stack := make([]byte, 4096)
			n := runtime.Stack(stack, false)
			err = fmt.Errorf("%s failed: %v, stack trace %s", e2eTest.Name, r, stack[:n])
		}
	}()

	startTime := time.Now()
	runner.Logger.Print("⏳running - %s", e2eTest.Description)

	// run e2e test, if args are not provided, use default args
	args := e2eTest.Args
	if len(args) == 0 {
		args = e2eTest.DefaultArgs()
	}
	e2eTest.E2ETest(runner, args)

	//check supplies
	if checkAccounting {
		if err := runner.CheckZRC20ReserveAndSupply(); err != nil {
			return err
		}
	}

	runner.Logger.Print("✅ completed in %s - %s", time.Since(startTime), e2eTest.Description)

	return err
}

// RunE2ETestsIntoReport runs a list of e2e tests by name in a list of e2e tests and returns a report
// The function doesn't return an error, it returns a report with the error
func (runner *E2ERunner) RunE2ETestsIntoReport(e2eTests []E2ETest) (TestReports, error) {
	// go through all tests
	reports := make(TestReports, 0, len(e2eTests))
	for _, test := range e2eTests {
		// get info before test
		balancesBefore, err := runner.GetAccountBalances(true)
		if err != nil {
			return nil, err
		}
		timeBefore := time.Now()

		// run test
		testErr := runner.RunE2ETest(test, false)
		if testErr != nil {
			runner.Logger.Print("test %s failed: %s", test.Name, testErr.Error())
		}

		// wait 5 sec to make sure we get updated balances
		time.Sleep(5 * time.Second)

		// get info after test
		balancesAfter, err := runner.GetAccountBalances(true)
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
