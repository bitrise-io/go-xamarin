package builder

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/analyzers/solution"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/utility"
)

// FileInfoWithPath ...
type FileInfoWithPath struct {
	modTime time.Time
	path    string
}

func fileInfos(dir string) ([]FileInfoWithPath, error) {
	fileInfos := []FileInfoWithPath{}

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileInfos = append(fileInfos, FileInfoWithPath{info.ModTime(), path})

		return nil
	}); err != nil {
		return []FileInfoWithPath{}, err
	}

	return fileInfos, nil
}

func isInTimeInterval(modTime, startTime, endTime time.Time) bool {
	return (modTime.After(startTime) || modTime.Equal(startTime)) && (modTime.Before(endTime) || modTime.Equal(endTime))
}

func filterFilesInfosForTimeWindow(fileInfos []FileInfoWithPath, startTime, endTime time.Time) []FileInfoWithPath {
	if startTime.IsZero() || endTime.IsZero() || startTime.Equal(endTime) {
		return []FileInfoWithPath{}
	}

	filteredFileInfos := []FileInfoWithPath{}

	for _, fileInfo := range fileInfos {
		if isInTimeInterval(fileInfo.modTime, startTime, endTime) {
			filteredFileInfos = append(filteredFileInfos, fileInfo)
		}
	}

	return filteredFileInfos
}

func findLastModifiedWithRegexps(fileInfos []FileInfoWithPath, regexps ...*regexp.Regexp) *FileInfoWithPath {
	var lastModTime time.Time
	latestFileInfo := new(FileInfoWithPath)

	for _, re := range regexps {
		for _, fileInfo := range fileInfos {
			fileName := filepath.Base(fileInfo.path)
			if match := re.FindString(fileName); match == fileName {
				if lastModTime.IsZero() || fileInfo.modTime.After(lastModTime) {
					lastModTime = fileInfo.modTime
					*latestFileInfo = fileInfo
				}
			}
		}
	}

	return latestFileInfo
}

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

func whitelistAllows(projectType constants.SDK, projectTypeWhiteList ...constants.SDK) bool {
	if len(projectTypeWhiteList) == 0 {
		return true
	}

	for _, filter := range projectTypeWhiteList {
		switch filter {
		case constants.SDKIOS:
			if projectType == constants.SDKIOS {
				return true
			}
		case constants.SDKTvOS:
			if projectType == constants.SDKTvOS {
				return true
			}
		case constants.SDKMacOS:
			if projectType == constants.SDKMacOS {
				return true
			}
		case constants.SDKAndroid:
			if projectType == constants.SDKAndroid {
				return true
			}
		}
	}

	return false
}

func isArchitectureArchiveable(architectures ...string) bool {
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

func isPlatformAnyCPU(platform string) bool {
	return (platform == "Any CPU" || platform == "AnyCPU")
}

func androidPackageName(manifestPth string) (string, error) {
	content, err := fileutil.ReadStringFromFile(manifestPth)
	if err != nil {
		return "", err
	}

	return androidPackageNameFromManifestContent(content)
}

func androidPackageNameFromManifestContent(manifestContent string) (string, error) {
	// package is attribute of the rott xml element
	manifestContent = "<a>" + manifestContent + "</a>"

	type Manifest struct {
		Package string `xml:"package,attr"`
	}

	type Result struct {
		Manifest Manifest `xml:"manifest"`
	}

	var result Result
	if err := xml.Unmarshal([]byte(manifestContent), &result); err != nil {
		return "", err
	}

	return result.Manifest.Package, nil
}

func exportLastModifiedArtifact(dir string, startTime, endTime time.Time, patterns ...string) (string, error) {
	regexps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		regexps[i] = regexp.MustCompile(pattern)
	}

	fileInfos, err := fileInfos(dir)
	if err != nil {
		return "", err
	}

	var fileInfo *FileInfoWithPath

	if fileInfosInTimeWindow := filterFilesInfosForTimeWindow(fileInfos, startTime, endTime); len(fileInfosInTimeWindow) > 0 {
		fileInfo = findLastModifiedWithRegexps(fileInfosInTimeWindow, regexps...)
	}

	if fileInfo == nil {
		fileInfo = findLastModifiedWithRegexps(fileInfos, regexps...)
	}

	if fileInfo != nil {
		return fileInfo.path, nil
	}
	return "", nil
}

func exportApk(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := exportLastModifiedArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i)%s.*signed\.apk$`, assemblyName),
		fmt.Sprintf(`(?i)%s\.apk$`, assemblyName),
	)
	if err != nil {
		return "", err
	}

	fmt.Printf("artificatPth: %s\n", artificatPth)

	if len(artificatPth) > 0 {
		return artificatPth, nil
	}

	log.Printf("")
	log.Warnf("Switching to legacy exporter")
	log.Printf("")

	apks, err := filepath.Glob(filepath.Join(outputDir, "*.apk"))
	if err != nil {
		return "", fmt.Errorf("failed to find apk, error: %s", err)
	}

	rePattern := fmt.Sprintf(`(?i)%s.*signed.apk`, assemblyName)
	re := regexp.MustCompile(rePattern)

	filteredApks := []string{}
	for _, apk := range apks {
		if match := re.FindString(apk); match != "" {
			filteredApks = append(filteredApks, apk)
		}
	}

	if len(filteredApks) == 0 {
		rePattern := fmt.Sprintf(`%s.apk`, assemblyName)
		re := regexp.MustCompile(rePattern)

		for _, apk := range apks {
			if match := re.FindString(apk); match != "" {
				filteredApks = append(filteredApks, apk)
			}
		}

		if len(filteredApks) == 0 {
			filteredApks = apks
		}
	}

	if len(filteredApks) == 0 {
		log.Errorf("Legacy exporter failed to find apk in (%s)", outputDir)
		return "", nil
	}

	return filteredApks[0], nil
}

func exportLatestIpa(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := exportLastModifiedArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.ipa$`, assemblyName),
		`(?i)\.ipa$`,
	)
	if err != nil {
		return "", err
	}

	if len(artificatPth) > 0 {
		return artificatPth, nil
	}

	return "", nil
}

func exportLatestXCArchive(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := exportLastModifiedArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.xcarchive$`, assemblyName),
		fmt.Sprintf(`(?i)\.xcarchive$`),
	)
	if err != nil {
		return "", err
	}

	if len(artificatPth) > 0 {
		return artificatPth, nil
	}

	return "", nil
}

func exportLatestXCArchiveFromXcodeArchives(assemblyName string, startTime, endTime time.Time) (string, error) {
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

	return exportLatestXCArchive(xcodeArchivesDir, assemblyName, startTime, endTime)
}

func exportAppDSYM(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := exportLastModifiedArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i)%s\.app\.dSYM$`, assemblyName),
		`(?i)\.app\.dSYM$`,
	)
	if err != nil {
		return "", err
	}

	if len(artificatPth) > 0 {
		return artificatPth, nil
	}

	log.Printf("")
	log.Warnf("Switching to legacy exporter")
	log.Printf("")

	pattern := filepath.Join(outputDir, "*.app.dSYM")
	dSYMs, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to find dsym with pattern (%s), error: %s", pattern, err)
	}
	if len(dSYMs) == 0 {
		return "", nil
	}

	rePattern := fmt.Sprintf("%s.app.dSYM", assemblyName)
	re := regexp.MustCompile(rePattern)

	filteredDsyms := []string{}
	for _, dSYM := range dSYMs {
		if match := re.FindString(dSYM); match != "" {
			filteredDsyms = append(filteredDsyms, dSYM)
		}
	}

	if len(filteredDsyms) == 0 {
		filteredDsyms = dSYMs
	}

	if len(filteredDsyms) == 0 {
		log.Errorf("Legacy exporter failed to find app.dSYM in (%s)", outputDir)
		return "", nil
	}

	return filteredDsyms[0], nil
}

func exportFrameworkDSYMs(outputDir string) ([]string, error) {
	// Multiplatform/iOS/bin/iPhone/Release/TTTAttributedLabel.framework.dSYM
	pattern := filepath.Join(outputDir, "*.framework.dSYM")
	dSYMs, err := filepath.Glob(pattern)
	if err != nil {
		return []string{}, fmt.Errorf("failed to find dsym with pattern (%s), error: %s", pattern, err)
	}
	return dSYMs, nil
}

func exportPKG(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := exportLastModifiedArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i)%s\.pkg$`, assemblyName),
		`(?i)\.pkg$`,
	)
	if err != nil {
		return "", err
	}

	if len(artificatPth) > 0 {
		return artificatPth, nil
	}

	log.Printf("")
	log.Warnf("Switching to legacy exporter")
	log.Printf("")

	pattern := filepath.Join(outputDir, "*.pkg")
	pkgs, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to find pkg with pattern (%s), error: %s", pattern, err)
	}
	if len(pkgs) == 0 {
		return "", nil
	}

	rePattern := fmt.Sprintf("%s.*.pkg", assemblyName)
	re := regexp.MustCompile(rePattern)

	filteredPKGs := []string{}
	for _, pkg := range pkgs {
		if match := re.FindString(pkg); match != "" {
			filteredPKGs = append(filteredPKGs, pkg)
		}
	}

	if len(filteredPKGs) == 0 {
		filteredPKGs = pkgs
	}

	if len(filteredPKGs) == 0 {
		log.Errorf("Legacy exporter failed to find pkg in (%s)", outputDir)
		return "", nil
	}

	return filteredPKGs[0], nil
}

func exportApp(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := exportLastModifiedArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i)%s\.app$`, assemblyName),
		`(?i)\.app$`,
	)
	if err != nil {
		return "", err
	}

	if len(artificatPth) > 0 {
		return artificatPth, nil
	}

	log.Printf("")
	log.Warnf("Switching to legacy exporter")
	log.Printf("")

	pattern := filepath.Join(outputDir, "*.app")
	apps, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to find app with pattern (%s), error: %s", pattern, err)
	}
	if len(apps) == 0 {
		return "", nil
	}

	rePattern := fmt.Sprintf("%s.app", assemblyName)
	re := regexp.MustCompile(rePattern)

	filteredAPPs := []string{}
	for _, app := range apps {
		if match := re.FindString(app); match != "" {
			filteredAPPs = append(filteredAPPs, app)
		}
	}

	if len(filteredAPPs) == 0 {
		filteredAPPs = apps
	}

	if len(filteredAPPs) == 0 {
		log.Errorf("Legacy exporter failed to find app in (%s)", outputDir)
		return "", nil
	}

	return filteredAPPs[0], nil
}

func exportDLL(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := exportLastModifiedArtifact(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i)%s\.dll$`, assemblyName),
		`(?i)\.dll$`,
	)
	if err != nil {
		return "", err
	}

	if len(artificatPth) > 0 {
		return artificatPth, nil
	}

	log.Printf("")
	log.Warnf("Switching to legacy exporter")
	log.Printf("")

	pattern := filepath.Join(outputDir, "*.dll")
	dlls, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to find dll with pattern (%s), error: %s", pattern, err)
	}
	if len(dlls) == 0 {
		return "", nil
	}

	rePattern := fmt.Sprintf("%s.dll", assemblyName)
	re := regexp.MustCompile(rePattern)

	filteredDLLs := []string{}
	for _, dll := range dlls {
		if match := re.FindString(dll); match != "" {
			filteredDLLs = append(filteredDLLs, dll)
		}
	}

	if len(filteredDLLs) == 0 {
		filteredDLLs = dlls
	}

	if len(filteredDLLs) == 0 {
		log.Errorf("Legacy exporter failed to find DLL in (%s)", outputDir)
		return "", nil
	}

	return filteredDLLs[0], nil
}
