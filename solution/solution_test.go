package solution

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func tmpSolutionWithContent(t *testing.T, content string) string {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__xamarin-builder-test__")
	require.NoError(t, err)
	pth := filepath.Join(tmpDir, "solution.sln")
	require.NoError(t, fileutil.WriteStringToFile(pth, content))
	return pth
}

func TestAnalyzeSolution(t *testing.T) {
	t.Log("ios test")
	{
		pth := tmpSolutionWithContent(t, iosTestSolutionContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		solution, err := analyzeSolution(pth, false)
		require.NoError(t, err)
		require.Equal(t, pth, solution.Pth)
		require.Equal(t, "FAE04EC0-301F-11D3-BF4B-00C04F79EFBC", solution.ID)

		// ConfigMap
		desiredConfigs := []string{
			"Debug|iPhoneSimulator",
			"Release|iPhone",
			"Release|iPhoneSimulator",
			"Debug|iPhone",
			"Debug|Any CPU",
			"Release|Any CPU",
		}
		desiredMappedConfigs := []string{
			"Debug|iPhoneSimulator",
			"Release|iPhone",
			"Release|iPhoneSimulator",
			"Debug|iPhone",
			"Debug|Any CPU",
			"Release|Any CPU",
		}
		for i := 0; i < len(desiredConfigs); i++ {
			value, ok := solution.ConfigMap[desiredConfigs[i]]
			require.Equal(t, true, ok)
			require.Equal(t, desiredMappedConfigs[i], value)
		}

		// ProjectMap
		desiredProjectIDs := []string{
			"90F3C584-FD69-4926-9903-6B9771847782",
			"BA48743D-06F3-4D2D-ACFD-EE2642CE155A",
			"99A825A6-6F99-4B94-9F65-E908A6347F1E",
			"ED150913-76EB-446F-8B78-DC77E5795703",
		}
		for i := 0; i < len(desiredProjectIDs); i++ {
			_, ok := solution.ProjectMap[desiredProjectIDs[i]]
			require.Equal(t, true, ok)
		}

		project := solution.ProjectMap["90F3C584-FD69-4926-9903-6B9771847782"]
		require.Equal(t, "90F3C584-FD69-4926-9903-6B9771847782", project.ID)
		require.Equal(t, "CreditCardValidator.iOS", project.Name)
		require.Equal(t, filepath.Join(dir, "CreditCardValidator.iOS/CreditCardValidator.iOS.csproj"), project.Pth)

		// Project Config mapping
		desiredProjectConfigs := []string{
			"Debug|Any CPU",
			"Debug|iPhone",
			"Debug|iPhoneSimulator",
			"Release|Any CPU",
			"Release|iPhone",
			"Release|iPhoneSimulator",
		}
		desiredMappedProjectConfigs := []string{
			"Debug|iPhoneSimulator",
			"Debug|iPhone",
			"Debug|iPhoneSimulator",
			"Release|iPhone",
			"Release|iPhone",
			"Release|iPhoneSimulator",
		}
		for i := 0; i < len(desiredConfigs); i++ {
			value, ok := project.ConfigMap[desiredProjectConfigs[i]]
			require.Equal(t, true, ok)
			require.Equal(t, desiredMappedProjectConfigs[i], value)
		}

		project = solution.ProjectMap["BA48743D-06F3-4D2D-ACFD-EE2642CE155A"]
		require.Equal(t, "BA48743D-06F3-4D2D-ACFD-EE2642CE155A", project.ID)
		require.Equal(t, "CreditCardValidator.iOS.UITests", project.Name)
		require.Equal(t, filepath.Join(dir, "CreditCardValidator.iOS.UITests/CreditCardValidator.iOS.UITests.csproj"), project.Pth)

		project = solution.ProjectMap["99A825A6-6F99-4B94-9F65-E908A6347F1E"]
		require.Equal(t, "99A825A6-6F99-4B94-9F65-E908A6347F1E", project.ID)
		require.Equal(t, "CreditCardValidator", project.Name)
		require.Equal(t, filepath.Join(dir, "CreditCardValidator/CreditCardValidator.csproj"), project.Pth)

		project = solution.ProjectMap["ED150913-76EB-446F-8B78-DC77E5795703"]
		require.Equal(t, "ED150913-76EB-446F-8B78-DC77E5795703", project.ID)
		require.Equal(t, "CreditCardValidator.iOS.NunitTests", project.Name)
		require.Equal(t, filepath.Join(dir, "CreditCardValidator.iOS.NunitTests/CreditCardValidator.iOS.NunitTests.csproj"), project.Pth)
	}

	t.Log("android test")
	{
		pth := tmpSolutionWithContent(t, androidTestSolutionContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		solution, err := analyzeSolution(pth, false)
		require.NoError(t, err)
		require.Equal(t, pth, solution.Pth)
		require.Equal(t, "FAE04EC0-301F-11D3-BF4B-00C04F79EFBC", solution.ID)

		// ConfigMap
		desiredConfigs := []string{
			"Debug|Any CPU",
			"Release|Any CPU",
		}
		desiredMappedConfigs := []string{
			"Debug|Any CPU",
			"Release|Any CPU",
		}
		for i := 0; i < len(desiredConfigs); i++ {
			value, ok := solution.ConfigMap[desiredConfigs[i]]
			require.Equal(t, true, ok)
			require.Equal(t, desiredMappedConfigs[i], value)
		}

		// ProjectMap
		desiredProjectIDs := []string{
			"9D1D32A3-D13F-4F23-B7D4-EF9D52B06E60",
			"048C57FD-A3A8-41E5-94B6-C41C3B4F5D95",
			"99A825A6-6F99-4B94-9F65-E908A6347F1E",
			"EF586485-1B11-4873-9D60-FFDBCBFE7E99",
		}
		for i := 0; i < len(desiredProjectIDs); i++ {
			_, ok := solution.ProjectMap[desiredProjectIDs[i]]
			require.Equal(t, true, ok)
		}

		project := solution.ProjectMap["9D1D32A3-D13F-4F23-B7D4-EF9D52B06E60"]
		require.Equal(t, "9D1D32A3-D13F-4F23-B7D4-EF9D52B06E60", project.ID)
		require.Equal(t, "CreditCardValidator.Droid", project.Name)
		require.Equal(t, filepath.Join(dir, "CreditCardValidator.Droid/CreditCardValidator.Droid.csproj"), project.Pth)

		// Project Config mapping
		desiredProjectConfigs := []string{
			"Debug|Any CPU",
			"Release|Any CPU",
		}
		desiredMappedProjectConfigs := []string{
			"Debug|AnyCPU",
			"Release|AnyCPU",
		}

		for i := 0; i < len(desiredConfigs); i++ {
			value, ok := project.ConfigMap[desiredProjectConfigs[i]]
			require.Equal(t, true, ok)
			require.Equal(t, desiredMappedProjectConfigs[i], value)
		}

		project = solution.ProjectMap["048C57FD-A3A8-41E5-94B6-C41C3B4F5D95"]
		require.Equal(t, "048C57FD-A3A8-41E5-94B6-C41C3B4F5D95", project.ID)
		require.Equal(t, "CreditCardValidator.Droid.UITests", project.Name)
		require.Equal(t, filepath.Join(dir, "CreditCardValidator.Droid.UITests/CreditCardValidator.Droid.UITests.csproj"), project.Pth)

		project = solution.ProjectMap["99A825A6-6F99-4B94-9F65-E908A6347F1E"]
		require.Equal(t, "99A825A6-6F99-4B94-9F65-E908A6347F1E", project.ID)
		require.Equal(t, "CreditCardValidator", project.Name)
		require.Equal(t, filepath.Join(dir, "CreditCardValidator/CreditCardValidator.csproj"), project.Pth)

		project = solution.ProjectMap["EF586485-1B11-4873-9D60-FFDBCBFE7E99"]
		require.Equal(t, "EF586485-1B11-4873-9D60-FFDBCBFE7E99", project.ID)
		require.Equal(t, "CreditCardValidator.Droid.NunitTests", project.Name)
		require.Equal(t, filepath.Join(dir, "CreditCardValidator.Droid.NunitTests/CreditCardValidator.Droid.NunitTests.csproj"), project.Pth)
	}

	t.Log("mac test")
	{
		pth := tmpSolutionWithContent(t, macTestSolutionContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		solution, err := analyzeSolution(pth, false)
		require.NoError(t, err)
		require.Equal(t, pth, solution.Pth)
		require.Equal(t, "FAE04EC0-301F-11D3-BF4B-00C04F79EFBC", solution.ID)

		// ConfigMap
		desiredConfigs := []string{
			"Debug|Any CPU",
			"Release|Any CPU",
		}
		desiredMappedConfigs := []string{
			"Debug|Any CPU",
			"Release|Any CPU",
		}
		for i := 0; i < len(desiredConfigs); i++ {
			value, ok := solution.ConfigMap[desiredConfigs[i]]
			require.Equal(t, true, ok)
			require.Equal(t, desiredMappedConfigs[i], value)
		}

		// ProjectMap
		desiredProjectIDs := []string{
			"4DA5EAC6-6F80-4FEC-AF81-194210F10B51",
		}
		for i := 0; i < len(desiredProjectIDs); i++ {
			_, ok := solution.ProjectMap[desiredProjectIDs[i]]
			require.Equal(t, true, ok)
		}

		project := solution.ProjectMap["4DA5EAC6-6F80-4FEC-AF81-194210F10B51"]
		require.Equal(t, "4DA5EAC6-6F80-4FEC-AF81-194210F10B51", project.ID)
		require.Equal(t, "Hello_Mac", project.Name)
		require.Equal(t, filepath.Join(dir, "Hello_Mac/Hello_Mac.csproj"), project.Pth)

		// Project Config mapping
		desiredProjectConfigs := []string{
			"Debug|Any CPU",
			"Release|Any CPU",
		}
		desiredMappedProjectConfigs := []string{
			"Debug|AnyCPU",
			"Release|AnyCPU",
		}
		for i := 0; i < len(desiredConfigs); i++ {
			value, ok := project.ConfigMap[desiredProjectConfigs[i]]
			require.Equal(t, true, ok)
			require.Equal(t, desiredMappedProjectConfigs[i], value)
		}
	}

	t.Log("tv test")
	{
		pth := tmpSolutionWithContent(t, tvTestSolutionContent)
		defer func() {
			require.NoError(t, os.Remove(pth))
		}()
		dir := filepath.Dir(pth)

		solution, err := analyzeSolution(pth, false)
		require.NoError(t, err)
		require.Equal(t, pth, solution.Pth)
		require.Equal(t, "FAE04EC0-301F-11D3-BF4B-00C04F79EFBC", solution.ID)

		// ConfigMap
		desiredConfigs := []string{
			"Debug|iPhoneSimulator",
			"Release|iPhone",
			"Release|iPhoneSimulator",
			"Debug|iPhone",
		}
		desiredMappedConfigs := []string{
			"Debug|iPhoneSimulator",
			"Release|iPhone",
			"Release|iPhoneSimulator",
			"Debug|iPhone",
		}

		for i := 0; i < len(desiredConfigs); i++ {
			value, ok := solution.ConfigMap[desiredConfigs[i]]
			require.Equal(t, true, ok)
			require.Equal(t, desiredMappedConfigs[i], value)
		}

		// ProjectMap
		desiredProjectIDs := []string{
			"51D9C362-2997-4029-B38F-06C36F17056E",
		}
		for i := 0; i < len(desiredProjectIDs); i++ {
			_, ok := solution.ProjectMap[desiredProjectIDs[i]]
			require.Equal(t, true, ok)
		}

		project := solution.ProjectMap["51D9C362-2997-4029-B38F-06C36F17056E"]
		require.Equal(t, "51D9C362-2997-4029-B38F-06C36F17056E", project.ID)
		require.Equal(t, "tvos", project.Name)
		require.Equal(t, filepath.Join(dir, "tvos/tvos.csproj"), project.Pth)

		// Project Config mapping
		desiredProjectConfigs := []string{
			"Debug|iPhoneSimulator",
			"Release|iPhone",
			"Release|iPhoneSimulator",
			"Debug|iPhone",
		}
		desiredMappedProjectConfigs := []string{
			"Debug|iPhoneSimulator",
			"Release|iPhone",
			"Release|iPhoneSimulator",
			"Debug|iPhone",
		}
		for i := 0; i < len(desiredConfigs); i++ {
			value, ok := project.ConfigMap[desiredProjectConfigs[i]]
			require.Equal(t, true, ok)
			require.Equal(t, desiredMappedProjectConfigs[i], value)
		}
	}
}
