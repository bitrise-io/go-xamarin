package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/bitrise-io/go-utils/pathutil"
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

		if info.IsDir() {
			return nil
		}

		fileInfos = append(fileInfos, FileInfoWithPath{info.ModTime(), path})

		return nil
	}); err != nil {
		return []FileInfoWithPath{}, err
	}

	return fileInfos, nil
}

func isInTimeInterval(t, startTime, endTime time.Time) bool {
	if startTime.After(endTime) {
		return false
	}
	return (t.After(startTime) || t.Equal(startTime)) && (t.Before(endTime) || t.Equal(endTime))
}

func filterFilesInfosForTimeWindow(fileInfos []FileInfoWithPath, startTime, endTime time.Time) []FileInfoWithPath {
	if startTime.IsZero() || endTime.IsZero() || startTime.Equal(endTime) {
		return []FileInfoWithPath{}
	}
	if startTime.After(endTime) {
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

// finds the last modified file matching to most strict regexp
// order of regexps should be: most strict -> less strict
func findLastModifiedWithFileNameRegexps(fileInfos []FileInfoWithPath, regexps ...*regexp.Regexp) *FileInfoWithPath {
	if len(fileInfos) == 0 {
		return nil
	}

	lastModifiedFileInfoPtr := new(FileInfoWithPath)

	if len(regexps) > 0 {
		for _, re := range regexps {
			for _, fileInfo := range fileInfos {
				fileName := filepath.Base(fileInfo.path)
				if match := re.FindString(fileName); match == fileName {
					if lastModifiedFileInfoPtr == nil || fileInfo.modTime.After(lastModifiedFileInfoPtr.modTime) {
						*lastModifiedFileInfoPtr = fileInfo
					}
				}
			}

			// return with the most strict match
			if len(lastModifiedFileInfoPtr.path) > 0 {
				return lastModifiedFileInfoPtr
			}
		}
	} else {
		for _, fileInfo := range fileInfos {
			if lastModifiedFileInfoPtr == nil || fileInfo.modTime.After(lastModifiedFileInfoPtr.modTime) {
				*lastModifiedFileInfoPtr = fileInfo
			}
		}
	}

	return lastModifiedFileInfoPtr
}

// exports the last modified file matching to most strict regexps within a time window
// order of regexps should be: most strict -> less strict
func findLastModifiedWithFileNameRegexpsInTimeWindow(dir string, startTime, endTime time.Time, patterns ...string) (string, error) {
	regexps := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		regexps[i] = regexp.MustCompile(pattern)
	}

	fileInfos, err := fileInfos(dir)
	if err != nil {
		return "", err
	}

	fileInfosInTimeWindow := filterFilesInfosForTimeWindow(fileInfos, startTime, endTime)
	fileInfo := findLastModifiedWithFileNameRegexps(fileInfosInTimeWindow, regexps...)

	if fileInfo != nil {
		return fileInfo.path, nil
	}
	return "", nil
}

func exportApk(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := findLastModifiedWithFileNameRegexpsInTimeWindow(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*signed.*\.apk$`, assemblyName),
		fmt.Sprintf(`(?i).*%s.*\.apk$`, assemblyName),
		`(?i).*signed.*\.apk$`,
		`(?i).*\.apk$`,
	)
	if err != nil {
		return "", err
	}
	return artificatPth, nil
}

func exportIpa(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := findLastModifiedWithFileNameRegexpsInTimeWindow(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.ipa$`, assemblyName),
		`(?i).*\.ipa$`,
	)
	if err != nil {
		return "", err
	}
	return artificatPth, nil
}

func exportXCArchive(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := findLastModifiedWithFileNameRegexpsInTimeWindow(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.xcarchive$`, assemblyName),
		fmt.Sprintf(`(?i).*\.xcarchive$`),
	)
	if err != nil {
		return "", err
	}
	return artificatPth, nil
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

	return exportXCArchive(xcodeArchivesDir, assemblyName, startTime, endTime)
}

func exportAppDSYM(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := findLastModifiedWithFileNameRegexpsInTimeWindow(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.app\.dSYM$`, assemblyName),
		`(?i).*\.app\.dSYM$`,
	)
	if err != nil {
		return "", err
	}
	return artificatPth, nil
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
	artificatPth, err := findLastModifiedWithFileNameRegexpsInTimeWindow(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.pkg$`, assemblyName),
		`(?i).*\.pkg$`,
	)
	if err != nil {
		return "", err
	}
	return artificatPth, nil
}

func exportApp(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := findLastModifiedWithFileNameRegexpsInTimeWindow(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.app$`, assemblyName),
		`(?i).*\.app$`,
	)
	if err != nil {
		return "", err
	}
	return artificatPth, nil
}

func exportDLL(outputDir, assemblyName string, startTime, endTime time.Time) (string, error) {
	artificatPth, err := findLastModifiedWithFileNameRegexpsInTimeWindow(outputDir, startTime, endTime,
		fmt.Sprintf(`(?i).*%s.*\.dll$`, assemblyName),
		`(?i).*\.dll$`,
	)
	if err != nil {
		return "", err
	}
	return artificatPth, nil
}
