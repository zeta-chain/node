package runner

import (
	"fmt"
	"sync"
)

// E2ETestFunc is a function representing a E2E test
// It takes a E2ERunner as an argument
type E2ETestFunc func(*E2ERunner, []string)

type E2ETestOpt func(*E2ETest)

// WithMinimumVersion sets a minimum zetacored version that is required to run the test.
// The test will be skipped if the minimum version is not satisfied.
func WithMinimumVersion(version string) E2ETestOpt {
	return func(t *E2ETest) {
		t.MinimumVersion = version
	}
}

// WithDependencies sets dependencies to the E2ETest to wait for completion
func WithDependencies(dependencies ...E2EDependency) E2ETestOpt {
	return func(t *E2ETest) {
		t.Dependencies = dependencies
	}
}

// E2ETest represents a E2E test with a name, args, description and test func
type E2ETest struct {
	Name           string
	Description    string
	Args           []string
	ArgsDefinition []ArgDefinition
	Dependencies   []E2EDependency
	E2ETest        E2ETestFunc
	MinimumVersion string
}

// NewE2ETest creates a new instance of E2ETest with specified parameters.
func NewE2ETest(
	name, description string,
	argsDefinition []ArgDefinition,
	e2eTestFunc E2ETestFunc,
	opts ...E2ETestOpt,
) E2ETest {
	test := E2ETest{
		Name:           name,
		Description:    description,
		ArgsDefinition: argsDefinition,
		E2ETest:        e2eTestFunc,
		Args:           []string{},
	}
	for _, opt := range opts {
		opt(&test)
	}
	return test
}

// E2EDependency defines a structure that holds a E2E test dependency
type E2EDependency struct {
	name      string
	waitGroup *sync.WaitGroup
}

// NewE2EDependency creates a new instance of E2Edependency with specified parameters.
func NewE2EDependency(name string) E2EDependency {
	var wg sync.WaitGroup
	wg.Add(1)

	return E2EDependency{
		name:      name,
		waitGroup: &wg,
	}
}

// Wait waits for the E2EDependency to complete
func (d *E2EDependency) Wait() {
	d.waitGroup.Wait()
}

// Done marks the E2EDependency as done
func (d *E2EDependency) Done() {
	d.waitGroup.Done()
}

// ArgDefinition defines a structure for holding an argument's description along with it's default value.
type ArgDefinition struct {
	Description  string
	DefaultValue string
}

// DefaultArgs extracts and returns array of default arguments from the ArgsDefinition.
func (e E2ETest) DefaultArgs() []string {
	defaultArgs := make([]string, len(e.ArgsDefinition))
	for i, spec := range e.ArgsDefinition {
		defaultArgs[i] = spec.DefaultValue
	}
	return defaultArgs
}

// ArgsDescription returns a string representing the arguments description in a readable format.
func (e E2ETest) ArgsDescription() string {
	argsDescription := ""
	for _, def := range e.ArgsDefinition {
		argDesc := fmt.Sprintf("%s (%s)", def.Description, def.DefaultValue)
		if argsDescription != "" {
			argsDescription += ", "
		}
		argsDescription += argDesc
	}
	return argsDescription
}

// E2ETestRunConfig defines the basic configuration for initiating an E2E test, including its name and optional runtime arguments.
type E2ETestRunConfig struct {
	Name string
	Args []string
}

// GetE2ETestsToRunByName prepares a list of E2ETests to run based on given test names without arguments
func (r *E2ERunner) GetE2ETestsToRunByName(availableTests []E2ETest, testNames ...string) ([]E2ETest, error) {
	tests := make([]E2ETestRunConfig, 0, len(testNames))
	for _, testName := range testNames {
		tests = append(tests, E2ETestRunConfig{
			Name: testName,
			Args: []string{},
		})
	}
	return r.GetE2ETestsToRunByConfig(availableTests, tests)
}

// GetE2ETestsToRunByConfig prepares a list of E2ETests to run based on provided test names and their corresponding arguments
func (r *E2ERunner) GetE2ETestsToRunByConfig(
	availableTests []E2ETest,
	testConfigs []E2ETestRunConfig,
) ([]E2ETest, error) {
	tests := []E2ETest{}
	for _, testSpec := range testConfigs {
		e2eTest, found := findE2ETestByName(availableTests, testSpec.Name)
		if !found {
			return nil, fmt.Errorf("e2e test %s not found", testSpec.Name)
		}
		if r.TestFilter != nil && !r.TestFilter.MatchString(e2eTest.Name) {
			continue
		}
		e2eTestToRun := E2ETest{
			Name:           e2eTest.Name,
			Description:    e2eTest.Description,
			ArgsDefinition: e2eTest.ArgsDefinition,
			Dependencies:   e2eTest.Dependencies,
			E2ETest:        e2eTest.E2ETest,
			MinimumVersion: e2eTest.MinimumVersion,
		}
		// update e2e test args
		e2eTestToRun.Args = testSpec.Args
		tests = append(tests, e2eTestToRun)
	}

	return tests, nil
}

// findE2ETestByName finds a e2e test by name
func findE2ETestByName(e2eTests []E2ETest, e2eTestName string) (E2ETest, bool) {
	for _, test := range e2eTests {
		if test.Name == e2eTestName {
			return test, true
		}
	}
	return E2ETest{}, false
}
