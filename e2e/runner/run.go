package runner

import (
	"time"
)

// RunE2ETests runs a list of e2e tests
func (r *E2ERunner) RunE2ETests(e2eTests []E2ETest) (err error) {
	for _, e2eTest := range e2eTests {
		if err := r.RunE2ETest(e2eTest, true); err != nil {
			return err
		}
	}
	return nil
}

// RunE2ETest runs a e2e test
func (r *E2ERunner) RunE2ETest(e2eTest E2ETest, checkAccounting bool) error {
	startTime := time.Now()
	r.Logger.Print("⏳running - %s", e2eTest.Description)

	// run e2e test, if args are not provided, use default args
	args := e2eTest.Args
	if len(args) == 0 {
		args = e2eTest.DefaultArgs()
	}
	e2eTest.E2ETest(r, args)

	//check supplies
	if checkAccounting {
		if err := r.CheckZRC20ReserveAndSupply(); err != nil {
			return err
		}
	}

	r.Logger.Print("✅ completed in %s - %s", time.Since(startTime), e2eTest.Description)

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
