package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/stretchr/testify/require"
)

func tmpProjectWithContent(t *testing.T, content string) string {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__xamarin-builder-test__")
	require.NoError(t, err)
	pth := filepath.Join(tmpDir, "project.csproj")
	require.NoError(t, fileutil.WriteStringToFile(pth, content))
	return pth
}

func stringSliceContainsOnly(slice []string, item ...string) bool {
	if len(slice) != len(item) {
		return false
	}

	testMap := map[string]bool{}
	for _, i := range slice {
		testMap[i] = false
	}
	for _, e := range item {
		_, ok := testMap[e]
		if !ok {
			return false
		}
		testMap[e] = true
	}
	for _, ok := range testMap {
		if !ok {
			return false
		}
	}
	return true
}

func frameworkSliceContainsOnly(slice []constants.TestFramework, item ...string) bool {
	stringSlice := []string{}
	for _, s := range slice {
		stringSlice = append(stringSlice, string(s))
	}
	return stringSliceContainsOnly(stringSlice, item...)
}

func TestAnalyzeProject(t *testing.T) {
	t.Log("ios test")
	{
		pth := tmpProjectWithContent(t, iosTestProjectContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		project, err := analyzeProject(pth)
		require.NoError(t, err)

		require.Equal(t, "90F3C584-FD69-4926-9903-6B9771847782", project.ID)
		require.Equal(t, constants.ProjectTypeIOS, project.ProjectType)
		require.Equal(t, "exe", project.OutputType)
		require.Equal(t, "CreditCardValidator.iOS", project.AssemblyName)
		require.Equal(t, 0, len(project.TestFrameworks))
		require.Equal(t, true, stringSliceContainsOnly(project.ReferredProjectIDs, "99A825A6-6F99-4B94-9F65-E908A6347F1E"))
		require.Equal(t, "", project.ManifestPth)
		require.Equal(t, false, project.AndroidApplication)

		// Configs
		config, ok := project.Configs["Debug|iPhoneSimulator"]
		require.Equal(t, true, ok)
		require.Equal(t, "Debug", config.Configuration)
		require.Equal(t, "iPhoneSimulator", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhoneSimulator/Debug"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "i386"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|iPhone"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "iPhone", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhone/Release"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "ARMv7", "ARM64"))
		require.Equal(t, true, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|iPhoneSimulator"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "iPhoneSimulator", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhoneSimulator/Release"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "i386"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Debug|iPhone"]
		require.Equal(t, true, ok)
		require.Equal(t, "Debug", config.Configuration)
		require.Equal(t, "iPhone", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhone/Debug"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "ARMv7", "ARM64"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)
	}

	t.Log("android test")
	{
		pth := tmpProjectWithContent(t, androidTestProjectContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		project, err := analyzeProject(pth)
		require.NoError(t, err)

		require.Equal(t, "9D1D32A3-D13F-4F23-B7D4-EF9D52B06E60", project.ID)
		require.Equal(t, constants.ProjectTypeAndroid, project.ProjectType)
		require.Equal(t, "library", project.OutputType)
		require.Equal(t, "CreditCardValidator.Droid", project.AssemblyName)
		require.Equal(t, 0, len(project.TestFrameworks))
		require.Equal(t, true, stringSliceContainsOnly(project.ReferredProjectIDs, "99A825A6-6F99-4B94-9F65-E908A6347F1E"))
		require.Equal(t, filepath.Join(dir, "Properties/AndroidManifest.xml"), project.ManifestPth)
		require.Equal(t, true, project.AndroidApplication)

		// Configs
		config, ok := project.Configs["Debug|AnyCPU"]
		require.Equal(t, true, ok)
		require.Equal(t, "Debug", config.Configuration)
		require.Equal(t, "AnyCPU", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/Debug"), config.OutputDir)
		require.Equal(t, 0, len(config.MtouchArchs))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|AnyCPU"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "AnyCPU", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/Release"), config.OutputDir)
		require.Equal(t, 0, len(config.MtouchArchs))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, true, config.SignAndroid)
	}

	t.Log("mac test")
	{
		pth := tmpProjectWithContent(t, macTestProjectContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		project, err := analyzeProject(pth)
		require.NoError(t, err)

		require.Equal(t, "4DA5EAC6-6F80-4FEC-AF81-194210F10B51", project.ID)
		require.Equal(t, constants.ProjectTypeMacOS, project.ProjectType)
		require.Equal(t, "exe", project.OutputType)
		require.Equal(t, "Hello_Mac", project.AssemblyName)
		require.Equal(t, 0, len(project.TestFrameworks))
		require.Equal(t, 0, len(project.ReferredProjectIDs))
		require.Equal(t, "", project.ManifestPth)
		require.Equal(t, false, project.AndroidApplication)

		// Configs
		config, ok := project.Configs["Debug|AnyCPU"]
		require.Equal(t, true, ok)
		require.Equal(t, "Debug", config.Configuration)
		require.Equal(t, "AnyCPU", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/Debug"), config.OutputDir)
		require.Equal(t, 0, len(config.MtouchArchs))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|AnyCPU"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "AnyCPU", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/Release"), config.OutputDir)
		require.Equal(t, 0, len(config.MtouchArchs))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)
	}

	t.Log("tv os test")
	{
		pth := tmpProjectWithContent(t, tvTestProjectContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		project, err := analyzeProject(pth)
		require.NoError(t, err)

		require.Equal(t, "51D9C362-2997-4029-B38F-06C36F17056E", project.ID)
		require.Equal(t, constants.ProjectTypeTvOS, project.ProjectType)
		require.Equal(t, "exe", project.OutputType)
		require.Equal(t, "tvos", project.AssemblyName)
		require.Equal(t, 0, len(project.TestFrameworks))
		require.Equal(t, 0, len(project.ReferredProjectIDs))
		require.Equal(t, "", project.ManifestPth)
		require.Equal(t, false, project.AndroidApplication)

		// Configs
		config, ok := project.Configs["Debug|iPhoneSimulator"]
		require.Equal(t, true, ok)
		require.Equal(t, "Debug", config.Configuration)
		require.Equal(t, "iPhoneSimulator", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhoneSimulator/Debug"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "x86_64"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|iPhone"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "iPhone", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhone/Release"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "ARM64"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|iPhoneSimulator"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "iPhoneSimulator", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhoneSimulator/Release"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "x86_64"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Debug|iPhone"]
		require.Equal(t, true, ok)
		require.Equal(t, "Debug", config.Configuration)
		require.Equal(t, "iPhone", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhone/Debug"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "ARM64"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)
	}

	t.Log("xamarin uitest test")
	{
		pth := tmpProjectWithContent(t, xamarinUITestProjectContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		project, err := analyzeProject(pth)
		require.NoError(t, err)

		require.Equal(t, "BA48743D-06F3-4D2D-ACFD-EE2642CE155A", project.ID)
		require.Equal(t, "", string(project.ProjectType))
		require.Equal(t, "library", project.OutputType)
		require.Equal(t, "CreditCardValidator.iOS.UITests", project.AssemblyName)
		require.Equal(t, true, frameworkSliceContainsOnly(project.TestFrameworks, "Xamarin.UITest", "nunit.framework"))
		require.Equal(t, true, stringSliceContainsOnly(project.ReferredProjectIDs, "90F3C584-FD69-4926-9903-6B9771847782"))
		require.Equal(t, "", project.ManifestPth)
		require.Equal(t, false, project.AndroidApplication)

		// Configs
		config, ok := project.Configs["Debug|AnyCPU"]
		require.Equal(t, true, ok)
		require.Equal(t, "Debug", config.Configuration)
		require.Equal(t, "AnyCPU", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/Debug"), config.OutputDir)
		require.Equal(t, 0, len(config.MtouchArchs))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|AnyCPU"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "AnyCPU", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/Release"), config.OutputDir)
		require.Equal(t, 0, len(config.MtouchArchs))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)
	}

	t.Log("nunit test")
	{
		pth := tmpProjectWithContent(t, nunitTestProjectContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		project, err := analyzeProject(pth)
		require.NoError(t, err)

		require.Equal(t, "ED150913-76EB-446F-8B78-DC77E5795703", project.ID)
		require.Equal(t, "", string(project.ProjectType))
		require.Equal(t, "library", project.OutputType)
		require.Equal(t, "CreditCardValidator.iOS.NunitTests", project.AssemblyName)
		require.Equal(t, true, frameworkSliceContainsOnly(project.TestFrameworks, "nunit.framework"))
		require.Equal(t, 0, len(project.ReferredProjectIDs))
		require.Equal(t, "", project.ManifestPth)
		require.Equal(t, false, project.AndroidApplication)

		// Configs
		config, ok := project.Configs["Debug|AnyCPU"]
		require.Equal(t, true, ok)
		require.Equal(t, "Debug", config.Configuration)
		require.Equal(t, "AnyCPU", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/Debug"), config.OutputDir)
		require.Equal(t, 0, len(config.MtouchArchs))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|AnyCPU"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "AnyCPU", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/Release"), config.OutputDir)
		require.Equal(t, 0, len(config.MtouchArchs))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)
	}

	t.Log("nunit lite test")
	{
		pth := tmpProjectWithContent(t, nunitLiteTestProjectContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		project, err := analyzeProject(pth)
		require.NoError(t, err)

		require.Equal(t, "95615CA5-0D75-4389-A6E0-78309A686712", project.ID)
		require.Equal(t, constants.ProjectTypeIOS, project.ProjectType)
		require.Equal(t, "exe", project.OutputType)
		require.Equal(t, "CreditCardValidator.iOS.NunitLiteTests", project.AssemblyName)
		require.Equal(t, true, frameworkSliceContainsOnly(project.TestFrameworks, "MonoTouch.NUnitLite"))
		require.Equal(t, 0, len(project.ReferredProjectIDs))
		require.Equal(t, "", project.ManifestPth)
		require.Equal(t, false, project.AndroidApplication)

		// Configs
		config, ok := project.Configs["Debug|iPhoneSimulator"]
		require.Equal(t, true, ok)
		require.Equal(t, "Debug", config.Configuration)
		require.Equal(t, "iPhoneSimulator", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhoneSimulator/Debug"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "i386"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|iPhone"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "iPhone", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhone/Release"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "ARMv7", "ARM64"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)

		config, ok = project.Configs["Release|iPhoneSimulator"]
		require.Equal(t, true, ok)
		require.Equal(t, "Release", config.Configuration)
		require.Equal(t, "iPhoneSimulator", config.Platform)
		require.Equal(t, filepath.Join(dir, "bin/iPhoneSimulator/Release"), config.OutputDir)
		require.Equal(t, true, stringSliceContainsOnly(config.MtouchArchs, "i386"))
		require.Equal(t, false, config.BuildIpa)
		require.Equal(t, false, config.SignAndroid)
	}
}
