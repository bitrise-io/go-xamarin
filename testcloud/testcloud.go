package testcloud

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-tools/go-xamarin/constants"
)

/*
   request = ['mono', "\"#{test_cloud}\"", 'submit', "\"#{apk_path}\"", options[:api_key]]
   request << options[:sign_parameters] if options[:sign_parameters]
   request << "--user #{options[:user]}"
   request << "--assembly-dir \"#{assembly_dir}\""
   request << "--devices #{options[:devices]}"
   request << '--async-json' if options[:async] == 'yes'
   request << "--series #{options[:series]}" if options[:series]
   request << "--nunit-xml #{@result_log_path}"
   request << '--fixture-chunk' if options[:parallelization] == 'by_test_fixture'
   request << '--test-chunk' if options[:parallelization] == 'by_test_chunk'
   request << options[:other_parameters]

   request = [
     "mono \"#{test_cloud}\"",
     "submit \"#{ipa_path}\"",
     options[:api_key],
     "--assembly-dir \"#{assembly_dir}\"",
     "--nunit-xml \"#{@result_log_path}\"",
     "--user #{options[:user]}",
     "--devices \"#{options[:devices]}\""
   ]
   request << '--async-json' if options[:async] == 'yes'
   request << "--dsym \"#{dsym_path}\"" if dsym_path
   request << "--series \"#{options[:series]}\"" if options[:series]
   request << '--fixture-chunk' if options[:parallelization] == 'by_test_fixture'
   request << '--test-chunk' if options[:parallelization] == 'by_test_chunk'
   request << options[:other_parameters].to_s if options[:other_parameters]
*/
type Parallelization string

const (
	// ParallelizationByTestFixture ...
	ParallelizationByTestFixture Parallelization = "fixture-chunk"
	// ParallelizationByTestChunk ...
	ParallelizationByTestChunk Parallelization = "test-chunk"
)

// Model ...
type Model struct {
	testCloudExePth string

	apkPth  string
	ipaPth  string
	dsymPth string

	apiKey          string
	user            string
	assemblyDir     string
	devices         string
	isAsyncJSON     bool
	series          string
	nunitXMLPth     string
	parallelization Parallelization

	signOptions   []string
	customOptions []string
}

// NewModel ...
func NewModel(testCloudExexPth string) Model {
	return Model{testCloudExePth: testCloudExexPth}
}

// SetAPKPth ...
func (testCloud *Model) SetAPKPth(apkPth string) *Model {
	testCloud.apkPth = apkPth
	return testCloud
}

// SetIPAPth ...
func (testCloud *Model) SetIPAPth(ipaPth string) *Model {
	testCloud.ipaPth = ipaPth
	return testCloud
}

// SetDSYMPth ...
func (testCloud *Model) SetDSYMPth(dsymPth string) *Model {
	testCloud.dsymPth = dsymPth
	return testCloud
}

// SetAPIKey ...
func (testCloud *Model) SetAPIKey(apiKey string) *Model {
	testCloud.apiKey = apiKey
	return testCloud
}

// SetUser ...
func (testCloud *Model) SetUser(user string) *Model {
	testCloud.user = user
	return testCloud
}

// SetAssemblyDir ...
func (testCloud *Model) SetAssemblyDir(assemblyDir string) *Model {
	testCloud.assemblyDir = assemblyDir
	return testCloud
}

// SetDevices ...
func (testCloud *Model) SetDevices(devices string) *Model {
	testCloud.devices = devices
	return testCloud
}

// SetIsAsyncJSON ...
func (testCloud *Model) SetIsAsyncJSON(isAsyncJSON bool) *Model {
	testCloud.isAsyncJSON = isAsyncJSON
	return testCloud
}

// SetSeries ...
func (testCloud *Model) SetSeries(series string) *Model {
	testCloud.series = series
	return testCloud
}

// SetNunitXMLPth ...
func (testCloud *Model) SetNunitXMLPth(nunitXMLPth string) *Model {
	testCloud.nunitXMLPth = nunitXMLPth
	return testCloud
}

// SetParallelization ...
func (testCloud *Model) SetParallelization(parallelization Parallelization) *Model {
	testCloud.parallelization = parallelization
	return testCloud
}

// SetSignOptions ...
func (testCloud *Model) SetSignOptions(options ...string) *Model {
	testCloud.signOptions = options
	return testCloud
}

// SetCustomOptions ...
func (testCloud *Model) SetCustomOptions(options ...string) *Model {
	testCloud.customOptions = options
	return testCloud
}

func (testCloud *Model) submitCommandSlice() []string {
	cmdSlice := []string{constants.MonoPath}
	cmdSlice = append(cmdSlice, testCloud.testCloudExePth)
	cmdSlice = append(cmdSlice, "submit")

	if testCloud.apkPth != "" {
		cmdSlice = append(cmdSlice, testCloud.apkPth)
	}

	if testCloud.ipaPth != "" {
		cmdSlice = append(cmdSlice, testCloud.ipaPth)
	}
	if testCloud.dsymPth != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("--dsym %s", testCloud.dsymPth))
	}

	cmdSlice = append(cmdSlice, testCloud.apiKey)

	for _, option := range testCloud.signOptions {
		cmdSlice = append(cmdSlice, option)
	}

	cmdSlice = append(cmdSlice, fmt.Sprintf("--user %s", testCloud.user))
	cmdSlice = append(cmdSlice, fmt.Sprintf("--assembly-dir %s", testCloud.assemblyDir))
	cmdSlice = append(cmdSlice, fmt.Sprintf("--devices %s", testCloud.devices))

	if testCloud.isAsyncJSON {
		cmdSlice = append(cmdSlice, "--async-json")
	}

	cmdSlice = append(cmdSlice, fmt.Sprintf("--series %s", testCloud.series))

	if testCloud.nunitXMLPth != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("--nunit-xml %s", testCloud.nunitXMLPth))
	}

	if testCloud.parallelization == ParallelizationByTestChunk {
		cmdSlice = append(cmdSlice, "--test-chunk")
	} else if testCloud.parallelization == ParallelizationByTestFixture {
		cmdSlice = append(cmdSlice, "--fixture-chunk")
	}

	for _, option := range testCloud.customOptions {
		cmdSlice = append(cmdSlice, option)
	}

	return cmdSlice
}

// PrintableCommand ...
func (testCloud Model) PrintableCommand() string {
	cmdSlice := testCloud.submitCommandSlice()

	return cmdex.PrintableCommandArgs(true, cmdSlice)
}

// Submit ...
func (testCloud Model) Submit() error {
	cmdSlice := testCloud.submitCommandSlice()

	command, err := cmdex.NewCommandFromSlice(cmdSlice)
	if err != nil {
		return err
	}

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
