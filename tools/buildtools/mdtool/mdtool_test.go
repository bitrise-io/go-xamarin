package mdtool

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/go-utils/testutil"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Log("it create new xbuild model")
	{
		mdtool := New("solution.sln")
		require.NotNil(t, mdtool)

		require.Equal(t, constants.MDToolPath, mdtool.buildTool)
		require.Equal(t, "solution.sln", mdtool.solutionPth)
		require.Equal(t, "", mdtool.projectName)
		require.Equal(t, "", mdtool.configuration)
		require.Equal(t, "", mdtool.platform)
		require.Equal(t, "", mdtool.target)

		require.Equal(t, 0, len(mdtool.customOptions))
	}
}

func TestSetProperties(t *testing.T) {
	t.Log("it sets target")
	{
		mdtool := New("solution.sln")
		require.NotNil(t, mdtool)
		require.Equal(t, "", mdtool.target)

		mdtool.SetTarget("build")
		require.Equal(t, "build", mdtool.target)
	}

	t.Log("it sets configuration")
	{
		mdtool := New("solution.sln")
		require.NotNil(t, mdtool)
		require.Equal(t, "", mdtool.configuration)

		mdtool.SetConfiguration("Release")
		require.Equal(t, "Release", mdtool.configuration)
	}

	t.Log("it sets platform")
	{
		mdtool := New("solution.sln")
		require.NotNil(t, mdtool)
		require.Equal(t, "", mdtool.platform)

		mdtool.SetPlatform("iPhone")
		require.Equal(t, "iPhone", mdtool.platform)
	}

	t.Log("it sets project name")
	{
		mdtool := New("solution.sln")
		require.NotNil(t, mdtool)
		require.Equal(t, "", mdtool.projectName)

		mdtool.SetProjectName("project.csproj")
		require.Equal(t, "project.csproj", mdtool.projectName)
	}

	t.Log("it appends custom options")
	{
		mdtool := New("solution.sln")
		require.NotNil(t, mdtool)
		require.Equal(t, 0, len(mdtool.customOptions))

		customOptions := []string{"/verbosity:minimal", "/nologo"}
		mdtool.SetCustomOptions(customOptions...)
		testutil.EqualSlicesWithoutOrder(t, customOptions, mdtool.customOptions)
	}
}

func TestBuildCommandSlice(t *testing.T) {
	t.Log("it build command slice from model")
	{
		mdtool := New("solution.sln")
		desired := []string{constants.MDToolPath, "solution.sln"}
		require.Equal(t, desired, mdtool.buildCommandSlice())

		mdtool.SetTarget("build")
		desired = []string{constants.MDToolPath, "build", "solution.sln"}
		require.Equal(t, desired, mdtool.buildCommandSlice())

		mdtool.SetConfiguration("Release")
		desired = []string{constants.MDToolPath, "build", "solution.sln", "-c:Release"}
		require.Equal(t, desired, mdtool.buildCommandSlice())

		mdtool.SetPlatform("iPhone")
		desired = []string{constants.MDToolPath, "build", "solution.sln", "-c:Release|iPhone"}
		require.Equal(t, desired, mdtool.buildCommandSlice())

		mdtool.SetPlatform("AnyCPU")
		desired = []string{constants.MDToolPath, "build", "solution.sln", "-c:Release"}
		require.Equal(t, desired, mdtool.buildCommandSlice())

		mdtool.SetProjectName("project.csproj")
		desired = []string{constants.MDToolPath, "build", "solution.sln", "-c:Release", "-p:project.csproj"}
		require.Equal(t, desired, mdtool.buildCommandSlice())

		mdtool.SetCustomOptions("-r:PREFIX")
		desired = []string{constants.MDToolPath, "build", "solution.sln", "-c:Release", "-p:project.csproj", "-r:PREFIX"}
		require.Equal(t, desired, mdtool.buildCommandSlice())
	}
}

func TestPrintableCommand(t *testing.T) {
	t.Log("it creates printable command")
	{
		mdtool := New("solution.sln")
		desired := fmt.Sprintf(`"%s" "solution.sln"`, constants.MDToolPath)
		require.Equal(t, desired, mdtool.PrintableCommand())

		mdtool.SetTarget("build")
		desired = fmt.Sprintf(`"%s" "build" "solution.sln"`, constants.MDToolPath)
		require.Equal(t, desired, mdtool.PrintableCommand())

		mdtool.SetConfiguration("Release")
		desired = fmt.Sprintf(`"%s" "build" "solution.sln" "-c:Release"`, constants.MDToolPath)
		require.Equal(t, desired, mdtool.PrintableCommand())

		mdtool.SetPlatform("iPhone")
		desired = fmt.Sprintf(`"%s" "build" "solution.sln" "-c:Release|iPhone"`, constants.MDToolPath)
		require.Equal(t, desired, mdtool.PrintableCommand())

		mdtool.SetPlatform("AnyCPU")
		desired = fmt.Sprintf(`"%s" "build" "solution.sln" "-c:Release"`, constants.MDToolPath)
		require.Equal(t, desired, mdtool.PrintableCommand())

		mdtool.SetProjectName("project.csproj")
		desired = fmt.Sprintf(`"%s" "build" "solution.sln" "-c:Release" "-p:project.csproj"`, constants.MDToolPath)
		require.Equal(t, desired, mdtool.PrintableCommand())

		mdtool.SetCustomOptions("-r:PREFIX")
		desired = fmt.Sprintf(`"%s" "build" "solution.sln" "-c:Release" "-p:project.csproj" "-r:PREFIX"`, constants.MDToolPath)
		require.Equal(t, desired, mdtool.PrintableCommand())
	}
}
