package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/solution"
	"github.com/bitrise-tools/go-xamarin/utility"
	"github.com/stretchr/testify/require"
)

func TestValidateSolutionPth(t *testing.T) {
	t.Log("it validates solution path")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		solutionPth := filepath.Join(tmpDir, "solution.sln")
		require.NoError(t, fileutil.WriteStringToFile(solutionPth, "solution"))
		require.NoError(t, validateSolutionPth(solutionPth))
	}

	t.Log("it fails if file not exist")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		solutionPth := filepath.Join(tmpDir, "solution.sln")
		require.Error(t, validateSolutionPth(solutionPth))
	}

	t.Log("it fails if path is not solution path")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		projectPth := filepath.Join(tmpDir, "project.csproj")
		require.Error(t, validateSolutionPth(projectPth))
	}
}

func TestValidateSolutionConfig(t *testing.T) {
	t.Log("it validates if solution config exist")
	{
		configuration := "Release"
		platform := "iPhone"
		config := utility.ToConfig(configuration, platform)

		solution := solution.Model{
			ConfigMap: map[string]string{
				config: config,
			},
		}

		require.NoError(t, validateSolutionConfig(solution, configuration, platform))
	}

	t.Log("it fails if solution config not exist")
	{
		configuration := "Release"
		platform := "iPhone"
		config := utility.ToConfig(configuration, platform)

		solution := solution.Model{
			ConfigMap: map[string]string{
				config: config,
			},
		}

		require.Error(t, validateSolutionConfig(solution, configuration, "Any CPU"))
	}
}

func TestIsProjectTypeAllowed(t *testing.T) {
	t.Log("empty whitelist means allow any project type")
	{
		whitelist := []constants.ProjectType{}
		require.Equal(t, true, isProjectTypeAllowed(constants.ProjectTypeIOS, whitelist...))
	}

	t.Log("it allows project type that exists in whitelist")
	{
		whitelist := []constants.ProjectType{constants.ProjectTypeIOS}
		require.Equal(t, true, isProjectTypeAllowed(constants.ProjectTypeIOS, whitelist...))
	}

	t.Log("it allows project type that exists in whitelist")
	{
		whitelist := []constants.ProjectType{constants.ProjectTypeAndroid, constants.ProjectTypeIOS}
		require.Equal(t, true, isProjectTypeAllowed(constants.ProjectTypeIOS, whitelist...))
	}

	t.Log("it allows project type that exists in whitelist")
	{
		whitelist := []constants.ProjectType{constants.ProjectTypeAndroid, constants.ProjectTypeIOS}
		require.Equal(t, true, isProjectTypeAllowed(constants.ProjectTypeAndroid, whitelist...))
	}

	t.Log("it does not allows project type that does not exists in whitelist")
	{
		whitelist := []constants.ProjectType{constants.ProjectTypeIOS}
		require.Equal(t, false, isProjectTypeAllowed(constants.ProjectTypeAndroid, whitelist...))
	}
}

func TestIsArchitectureArchiveablet(t *testing.T) {
	t.Log("default architectures is armv7")
	{
		require.Equal(t, true, isArchitectureArchiveable())
	}

	t.Log("arm architectures are archivables")
	{
		require.Equal(t, true, isArchitectureArchiveable("armv7"))
	}

	t.Log("it is case insensitive")
	{
		require.Equal(t, true, isArchitectureArchiveable("ARM7"))
	}

	t.Log("x86 architectures are not archivables")
	{
		require.Equal(t, false, isArchitectureArchiveable("x86"))
	}
}

func TestIsPlatformAnyCPU(t *testing.T) {
	t.Log("true for Any CPU")
	{
		require.Equal(t, true, isPlatformAnyCPU("Any CPU"))
	}

	t.Log("true for AnyCPU")
	{
		require.Equal(t, true, isPlatformAnyCPU("AnyCPU"))
	}

	t.Log("false for other platforms")
	{
		require.Equal(t, false, isPlatformAnyCPU("iPhone"))
	}
}

func TestAndroidPackageNameFromManifestContent(t *testing.T) {
	t.Log("it finds package name in manifest")
	{
		packageName, err := androidPackageNameFromManifestContent(manifestFileContent)
		require.NoError(t, err)
		require.Equal(t, "hu.bitrise.test", packageName)
	}
}

func createTestFile(t *testing.T, tmpDir, relPth string) {
	pth := filepath.Join(tmpDir, relPth)
	dirPth := filepath.Dir(pth)
	require.NoError(t, os.MkdirAll(dirPth, 0777))
	require.NoError(t, fileutil.WriteStringToFile(pth, "test"))
}

func TestExportApk(t *testing.T) {
	t.Log("it retruns empty path if no apk found")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		output, err := exportApk(tmpDir, "com.bitrise.xamarin.sampleapp")
		require.NoError(t, err)
		require.Equal(t, "", output)
	}

	t.Log("it finds apk")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"com.bitrise.xamarin.sampleapp.apk",
			"FormsViewGroup.dll",
			"Java.Interop.dll.mdb",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportApk(tmpDir, "com.bitrise.xamarin.sampleapp")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "com.bitrise.xamarin.sampleapp.apk"), output)
	}

	t.Log("it finds apk - even if assembly name empty")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"com.bitrise.xamarin.sampleapp.apk",
			"FormsViewGroup.dll",
			"Java.Interop.dll.mdb",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportApk(tmpDir, "")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "com.bitrise.xamarin.sampleapp.apk"), output)
	}

	t.Log("it prefers signed apk")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"com.bitrise.xamarin.sampleapp.apk",
			"com.bitrise.xamarin.sampleapp-Signed.apk",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportApk(tmpDir, "")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "com.bitrise.xamarin.sampleapp-Signed.apk"), output)
	}
}

func TestExportLatestXCArchiveFromXcodeArchives(t *testing.T) {
	t.Log("it retruns empty path if no xcarchive found")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		output, err := exportLatestXCArchive(tmpDir, "XamarinSampleApp.iOS")
		require.NoError(t, err)
		require.Equal(t, "", output)
	}

	t.Log("it sorts by filename")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"2016-07-10/XamarinSampleApp.iOS 10-07-16 4.41 PM.xcarchive", // latest
			"2016-07-10/XamarinSampleApp.iOS 10-07-16 3.41 PM 2.xcarchive",
			"2016-07-10/XamarinSampleApp.iOS 10-07-16 2.41 PM.xcarchive",
			"2016-07-10/XamarinSampleApp.iOS 10-07-16 1.41 PM.xcarchive",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportLatestXCArchive(tmpDir, "XamarinSampleApp.iOS")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "2016-07-10/XamarinSampleApp.iOS 10-07-16 4.41 PM.xcarchive"), output)
	}

	t.Log("it sorts by filename - even if count number in pth")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"2016-07-10/XamarinSampleApp.iOS 10-07-16 4.41 PM 2.xcarchive", // latest
			"2016-07-10/XamarinSampleApp.iOS 10-07-16 4.41 PM.xcarchive",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportLatestXCArchive(tmpDir, "XamarinSampleApp.iOS")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "2016-07-10/XamarinSampleApp.iOS 10-07-16 4.41 PM 2.xcarchive"), output)
	}

	t.Log("it sorts by dirname")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"2016-07-10/XamarinSampleApp.iOS 10-07-16 4.41 PM.xcarchive", // latest
			"2016-06-10/XamarinSampleApp.iOS 10-07-16 4.41 PM.xcarchive",
			"2016-05-10/XamarinSampleApp.iOS 10-07-16 4.41 PM.xcarchive",
			"2016-04-10/XamarinSampleApp.iOS 10-07-16 4.41 PM.xcarchive",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportLatestXCArchive(tmpDir, "XamarinSampleApp.iOS")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "2016-07-10/XamarinSampleApp.iOS 10-07-16 4.41 PM.xcarchive"), output)
	}

	t.Log("it retruns latest xcarchive if assembly name empty")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"2016-07-10/a 10-07-16 3.45 PM.xcarchive", // latest
			"2016-07-10/b 10-07-16 1.41 PM.xcarchive",
			"2016-07-10/c 10-07-16 2.41 PM.xcarchive",
			"2016-07-10/d 10-07-16 3.41 PM.xcarchive",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportLatestXCArchive(tmpDir, "")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "2016-07-10/a 10-07-16 3.45 PM.xcarchive"), output)
	}
}

func TestExportLatestIpa(t *testing.T) {
	t.Log("it retruns empty path if no ipa found")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		output, err := exportLatestIpa(tmpDir, "XamarinSampleApp.iOS")
		require.NoError(t, err)
		require.Equal(t, "", output)
	}

	t.Log("it sorts by dirname")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"Multiplatform.iOS 2016-10-06 11-45-23/Multiplatform.iOS.ipa", // latest
			"Multiplatform.iOS 2016-09-06 11-45-23 2/Multiplatform.iOS.ipa",
			"Multiplatform.iOS 2016-08-06 11-45-23/Multiplatform.iOS.ipa",
			"Multiplatform.iOS 2016-07-06 17-45-23/Multiplatform.iOS.ipa",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportLatestIpa(tmpDir, "Multiplatform.iOS")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "Multiplatform.iOS 2016-10-06 11-45-23/Multiplatform.iOS.ipa"), output)
	}

	t.Log("it sorts by dirname - even if count number in pth")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"Multiplatform.iOS 2016-10-06 11-45-23 2/Multiplatform.iOS.ipa", // latest
			"Multiplatform.iOS 2016-10-06 11-45-23/Multiplatform.iOS.ipa",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportLatestIpa(tmpDir, "Multiplatform.iOS")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "Multiplatform.iOS 2016-10-06 11-45-23 2/Multiplatform.iOS.ipa"), output)
	}

	t.Log("it retruns latest ipa if assembly name empty")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"a 2016-10-06 11-45-25/Multiplatform.iOS.ipa", // latest
			"b 2016-10-06 11-45-23/Multiplatform.iOS.ipa",
			"c 2016-10-06 11-44-23/Multiplatform.iOS.ipa",
			"d 2016-10-06 11-43-23/Multiplatform.iOS.ipa",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportLatestIpa(tmpDir, "")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "a 2016-10-06 11-45-25/Multiplatform.iOS.ipa"), output)
	}
}

func TestExportAppDsym(t *testing.T) {
	t.Log("it retruns empty path if no dSYM found")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		output, err := exportAppDSYM(tmpDir, "Multiplatform.iOS")
		require.NoError(t, err)
		require.Equal(t, "", output)
	}

	t.Log("it finds dSYM")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"Multiplatform.iOS.app.dSYM",
			"TTTAttributedLabel.framework.dSYM",
			"a.app.dSYM",
			"Multiplatform.iOS.app",
			"Multiplatform.iOS.exe",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportAppDSYM(tmpDir, "Multiplatform.iOS")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "Multiplatform.iOS.app.dSYM"), output)
	}

	t.Log("it finds dSYM - even if assembly name empty")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"Multiplatform.iOS.app.dSYM",
			"TTTAttributedLabel.framework.dSYM",
			"Multiplatform.iOS.app",
			"Multiplatform.iOS.exe",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportAppDSYM(tmpDir, "")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "Multiplatform.iOS.app.dSYM"), output)
	}
}

func TestExportFrameworkDsyms(t *testing.T) {
	t.Log("it retruns empty slice if no dSYM found")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		output, err := exportFrameworkDSYMs(tmpDir)
		require.NoError(t, err)
		require.Equal(t, 0, len(output))
	}

	t.Log("it finds dSYMs")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"Multiplatform.iOS.app.dSYM",
			"TTTAttributedLabel.framework.dSYM",
			"a.app.dSYM",
			"Multiplatform.iOS.app",
			"Multiplatform.iOS.exe",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportFrameworkDSYMs(tmpDir)
		require.NoError(t, err)
		require.Equal(t, []string{filepath.Join(tmpDir, "TTTAttributedLabel.framework.dSYM")}, output)
	}
}

func TestExportPKG(t *testing.T) {
	t.Log("it retruns empty path if no pkg found")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		output, err := exportPKG(tmpDir, "Multiplatform.Mac")
		require.NoError(t, err)
		require.Equal(t, "", output)
	}

	t.Log("it finds pkg")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"Multiplatform.Mac-1.0.pkg",
			"Multiplatform.Mac.app",
			"Multiplatform.Mac.exe",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportPKG(tmpDir, "Multiplatform.Mac")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "Multiplatform.Mac-1.0.pkg"), output)
	}

	t.Log("it finds pkg - even if assembly name empty")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"Multiplatform.Mac-1.0.pkg",
			"Multiplatform.Mac.app",
			"Multiplatform.Mac.exe",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportPKG(tmpDir, "Multiplatform.Mac")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "Multiplatform.Mac-1.0.pkg"), output)
	}
}

func TestExportApp(t *testing.T) {
	t.Log("it retruns empty path if no pkg found")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		output, err := exportApp(tmpDir, "Multiplatform.Mac")
		require.NoError(t, err)
		require.Equal(t, "", output)
	}

	t.Log("it finds pkg")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"Multiplatform.Mac-1.0.pkg",
			"Multiplatform.Mac.app",
			"Multiplatform.Mac.exe",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportApp(tmpDir, "Multiplatform.Mac")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "Multiplatform.Mac.app"), output)
	}

	t.Log("it finds pkg - even if assembly name empty")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("utility_test")
		require.NoError(t, err)

		archives := []string{
			"Multiplatform.Mac-1.0.pkg",
			"Multiplatform.Mac.app",
			"Multiplatform.Mac.exe",
		}

		for _, archive := range archives {
			createTestFile(t, tmpDir, archive)
		}

		output, err := exportApp(tmpDir, "Multiplatform.Mac")
		require.NoError(t, err)
		require.Equal(t, filepath.Join(tmpDir, "Multiplatform.Mac.app"), output)
	}
}