package builder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/xamarin-builder/constants"
	"github.com/bitrise-tools/xamarin-builder/solution"
	"github.com/bitrise-tools/xamarin-builder/utility"
)

const (
	packageNamePattern = `(?i)<manifest.*package="(?P<package_name>.*)">`
)

func validateSolutionPth(pth string) error {
	ext := filepath.Ext(pth)
	if ext != constants.SolutionExt {
		return fmt.Errorf("path is not a solution file path: %s", pth)
	}
	if exist, err := pathutil.IsPathExists(pth); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("solution not exist at: %s", pth)
	}
	return nil
}

func validateSolutionConfig(solution solution.Model, configuration, platform string) error {
	config := utility.ToConfig(configuration, platform)
	if _, ok := solution.ConfigMap[config]; !ok {
		return fmt.Errorf("invalid solution config, available: %v", solution.ConfigList())
	}
	return nil
}

func isFilterForProjectType(projectType constants.XamarinProjectType, projectTypeFilter ...constants.ProjectType) bool {
	for _, filter := range projectTypeFilter {
		switch filter {
		case constants.Android:
			if projectType == constants.XamarinAndroid {
				return true
			}
		case constants.Ios:
			if projectType == constants.XamarinIos {
				return true
			}
		case constants.Mac:
			if projectType == constants.XamarinMac || projectType == constants.MonoMac {
				return true
			}
		case constants.TVos:
			if projectType == constants.XamarinTVOS {
				return true
			}
		}
	}
	return false
}

func archiveableArchitecture(architectures []string) bool {
	// default is armv7
	if len(architectures) == 0 {
		return true
	}

	for _, arch := range architectures {
		arch = strings.ToLower(arch)
		if !strings.HasPrefix(arch, "arm") {
			return false
		}
	}

	return true
}

func androidPackageName(manifestPth string) (string, error) {
	packageName := ""

	content, err := fileutil.ReadStringFromFile(manifestPth)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := regexp.MustCompile(packageNamePattern).FindStringSubmatch(line); len(matches) == 2 {
			packageName = matches[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return packageName, nil
}

func exportApk(outputDir, manifestPth string, isSigned bool) (string, error) {
	apkPth := ""

	packageName, err := androidPackageName(manifestPth)
	if err != nil {
		return "", err
	}

	apks, err := filepath.Glob(filepath.Join(outputDir, "*.apk"))
	if err != nil {
		return "", err
	}

	if isSigned {
		pattern := fmt.Sprintf(`(?i)%s.*signed\.apk`, packageName)

		for _, apk := range apks {
			if match := regexp.MustCompile(pattern).FindString(apk); match != "" {
				apkPth = apk
			}
		}
	}

	if apkPth == "" {
		pattern := fmt.Sprintf(`(?i)%s.*\.apk`, packageName)

		for _, apk := range apks {
			if match := regexp.MustCompile(pattern).FindString(apk); match != "" {
				apkPth = apk
			}
		}
	}

	if apkPth == "" {
		pattern := fmt.Sprintf(`(?i).*\.apk`)

		for _, apk := range apks {
			if match := regexp.MustCompile(pattern).FindString(apk); match != "" {
				apkPth = apk
			}
		}
	}

	return apkPth, nil
}

func exportLatestXCArchiveFromXcodeArchives(projectName string) (string, error) {
	userHomeDir := os.Getenv("HOME")
	if userHomeDir == "" {
		return "", fmt.Errorf("failed to get user home dir")
	}
	xcodeArchivesDir := filepath.Join(userHomeDir, "Library/Developer/Xcode/Archives")
	if exist, err := pathutil.IsDirExists(xcodeArchivesDir); err != nil {
		return "", err
	} else if !exist {
		return "", fmt.Errorf("no default Xcode archive path found at: %s", xcodeArchivesDir)
	}

	latestArchive := ""
	latestArchiveDate := time.Time{}

	pattern := filepath.Join(xcodeArchivesDir, "*", projectName+" *.xcarchive")

	archives, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	for _, archive := range archives {
		base := filepath.Base(archive)
		ext := filepath.Ext(archive)
		dateStr := strings.TrimPrefix(base, projectName)
		dateStr = strings.TrimSuffix(dateStr, ext)
		dateStr = strings.TrimSpace(dateStr)

		archiveDate, err := time.Parse("1-2-06 3.04 PM", dateStr)
		if err != nil {
			return "", err
		}

		if latestArchive == "" || archiveDate.After(latestArchiveDate) {
			latestArchive = archive
			latestArchiveDate = archiveDate
		}
	}

	return latestArchive, nil
}

func exportIpa(outputDir, assemblyName string) (string, error) {
	pattern := filepath.Join(outputDir, assemblyName+"*", assemblyName+".ipa")
	ipas, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(ipas) == 0 {
		return "", nil
	}

	latestIpaPth := ""
	latestIpaDate := time.Time{}

	for _, ipa := range ipas {
		ipaDir := filepath.Dir(ipa)
		ipaDirBase := filepath.Base(ipaDir)
		dateStr := strings.TrimPrefix(ipaDirBase, assemblyName)
		dateStr = strings.TrimSpace(dateStr)

		ipaDate, err := time.Parse("2006-01-02 15-04-05", dateStr)
		if err != nil {
			return "", err
		}

		if latestIpaPth == "" || ipaDate.After(latestIpaDate) {
			latestIpaPth = ipa
			latestIpaDate = ipaDate
		}
	}

	return latestIpaPth, nil
}

func exportPkg(outputDir, assemblyName string) (string, error) {
	pattern := filepath.Join(outputDir, assemblyName+"*.pkg")
	pkgs, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(pkgs) == 0 {
		return "", nil
	}
	return pkgs[0], nil
}

func exportApp(outputDir, assemblyName string) (string, error) {
	pattern := filepath.Join(outputDir, assemblyName+"*.app")
	apps, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(apps) == 0 {
		return "", nil
	}
	return apps[0], nil
}

// TimeoutHandlerModel ....
type TimeoutHandlerModel struct {
	timer   *time.Timer
	timeout time.Duration

	running bool

	onTimeout func()
}

// Start ...
func (handler *TimeoutHandlerModel) Start() {
	if &handler.timeout != nil {
		handler.timer = time.NewTimer(handler.timeout)
		handler.running = true

		go func() {
			for _ = range handler.timer.C {
				if handler.onTimeout != nil {
					handler.onTimeout()
				}
			}
		}()
	}
}

// Stop ...
func (handler *TimeoutHandlerModel) Stop() {
	if handler.running {
		handler.timer.Stop()
		handler.running = false
	}
}

// Running ...
func (handler TimeoutHandlerModel) Running() bool {
	return handler.running
}

// NewTimeoutHandler ...
func NewTimeoutHandler(timeout time.Duration, onTimeout func()) TimeoutHandlerModel {
	return TimeoutHandlerModel{
		timeout:   timeout,
		onTimeout: onTimeout,
	}
}

func runCommandInDiagnosticMode(command cmdex.CommandModel, checkPattern string, retryOnHang bool) error {
	log.Warn("Run in diagnostic mode")

	cmd := command.GetCmd()
	timeout := false

	// Create a timer that will FORCE kill the process if normal kill does not work
	var forceKillError error
	forceKillTimeoutHandler := NewTimeoutHandler(1*time.Second, func() {
		log.Warn("Timeout")
		timeout = true
		forceKillError = cmd.Process.Signal(syscall.SIGKILL)
	})
	// ----

	// Create a timer that will kill the process
	var killError error
	killTimeoutHandler := NewTimeoutHandler(1*time.Second, func() {
		log.Warn("Timeout")
		timeout = true
		forceKillTimeoutHandler.Start()
		killError = cmd.Process.Signal(syscall.SIGQUIT)
	})
	// ----

	// Redirect output
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdoutReader)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)

			killTimeoutHandler.Stop()
			if strings.Contains(strings.TrimSpace(line), checkPattern) {
				killTimeoutHandler.Start()
			}
		}
	}()
	if err := scanner.Err(); err != nil {
		return err
	}
	// ----

	if err := cmd.Start(); err != nil {
		return err
	}

	// Only proceed once the process has finished
	cmdErr := cmd.Wait()

	killTimeoutHandler.Stop()
	forceKillTimeoutHandler.Stop()

	if timeout {
		return fmt.Errorf("timed out")
	}
	if cmdErr != nil {
		return cmdErr
	}
	if killError != nil {
		return killError
	}
	if forceKillError != nil {
		return forceKillError
	}

	return nil
	// ----
}
