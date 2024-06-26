package runner

import "fmt"

// E2ETestFunc is a function representing a E2E test
// It takes a E2ERunner as an argument
type E2ETestFunc func(*E2ERunner, []string)

// E2ETest represents a E2E test with a name, args, description and test func
type E2ETest struct {
	Name           string
	Description    string
	Args           []string
	ArgsDefinition []ArgDefinition
	E2ETest        E2ETestFunc
}

// NewE2ETest creates a new instance of E2ETest with specified parameters.
func NewE2ETest(name, description string, argsDefinition []ArgDefinition, e2eTestFunc E2ETestFunc) E2ETest {
	return E2ETest{
		Name:           name,
		Description:    description,
		ArgsDefinition: argsDefinition,
		E2ETest:        e2eTestFunc,
		Args:           []string{},
	}
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
	tests := []E2ETestRunConfig{}
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
		e2eTestToRun := NewE2ETest(
			e2eTest.Name,
			e2eTest.Description,
			e2eTest.ArgsDefinition,
			e2eTest.E2ETest,
		)
		// update e2e test args
		e2eTestToRun.Args = testSpec.Args
		tests = append(tests, e2eTestToRun)
	}

	return tests, nil
}

// findE2ETest finds a e2e test by name
func findE2ETestByName(e2eTests []E2ETest, e2eTestName string) (E2ETest, bool) {
	for _, test := range e2eTests {
		if test.Name == e2eTestName {
			return test, true
		}
	}
	return E2ETest{}, false
}
