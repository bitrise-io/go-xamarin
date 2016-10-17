package xbuild

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
		xbuild := New("solution.sln")
		require.NotNil(t, xbuild)

		require.Equal(t, constants.XbuildPath, xbuild.buildTool)
		require.Equal(t, "solution.sln", xbuild.solutionPth)
		require.Equal(t, "", xbuild.configuration)
		require.Equal(t, "", xbuild.platform)
		require.Equal(t, "", xbuild.target)

		require.Equal(t, false, xbuild.buildIpa)
		require.Equal(t, false, xbuild.archiveOnBuild)

		require.Equal(t, 0, len(xbuild.customOptions))
	}
}

func TestSetProperties(t *testing.T) {
	t.Log("it sets target")
	{
		xbuild := New("solution.sln")
		require.NotNil(t, xbuild)
		require.Equal(t, "", xbuild.target)

		xbuild.SetTarget("Build")
		require.Equal(t, "Build", xbuild.target)
	}

	t.Log("it sets configuration")
	{
		xbuild := New("solution.sln")
		require.NotNil(t, xbuild)
		require.Equal(t, "", xbuild.configuration)

		xbuild.SetConfiguration("Release")
		require.Equal(t, "Release", xbuild.configuration)
	}

	t.Log("it sets platform")
	{
		xbuild := New("solution.sln")
		require.NotNil(t, xbuild)
		require.Equal(t, "", xbuild.platform)

		xbuild.SetPlatform("iPhone")
		require.Equal(t, "iPhone", xbuild.platform)
	}

	t.Log("it sets build ipa")
	{
		xbuild := New("solution.sln")
		require.NotNil(t, xbuild)
		require.Equal(t, false, xbuild.buildIpa)

		xbuild.SetBuildIpa(true)
		require.Equal(t, true, xbuild.buildIpa)
	}

	t.Log("it sets archive on build")
	{
		xbuild := New("solution.sln")
		require.NotNil(t, xbuild)
		require.Equal(t, false, xbuild.archiveOnBuild)

		xbuild.SetArchiveOnBuild(true)
		require.Equal(t, true, xbuild.archiveOnBuild)
	}

	t.Log("it appends custom options")
	{
		xbuild := New("solution.sln")
		require.NotNil(t, xbuild)
		require.Equal(t, 0, len(xbuild.customOptions))

		customOptions := []string{"/verbosity:minimal", "/nologo"}
		xbuild.SetCustomOptions(customOptions...)
		testutil.EqualSlicesWithoutOrder(t, customOptions, xbuild.customOptions)
	}
}

func TestBuildCommandSlice(t *testing.T) {
	t.Log("it build command slice from model")
	{
		xbuild := New("solution.sln")
		desired := []string{constants.XbuildPath, "solution.sln", "/p:SolutionDir=."}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetTarget("Build")
		desired = []string{constants.XbuildPath, "solution.sln", "/target:Build", "/p:SolutionDir=."}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetConfiguration("Release")
		desired = []string{constants.XbuildPath, "solution.sln", "/target:Build", "/p:SolutionDir=.", "/p:Configuration=Release"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetPlatform("iPhone")
		desired = []string{constants.XbuildPath, "solution.sln", "/target:Build", "/p:SolutionDir=.", "/p:Configuration=Release", "/p:Platform=iPhone"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetArchiveOnBuild(true)
		desired = []string{constants.XbuildPath, "solution.sln", "/target:Build", "/p:SolutionDir=.", "/p:Configuration=Release", "/p:Platform=iPhone", "/p:ArchiveOnBuild=true"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetBuildIpa(true)
		desired = []string{constants.XbuildPath, "solution.sln", "/target:Build", "/p:SolutionDir=.", "/p:Configuration=Release", "/p:Platform=iPhone", "/p:ArchiveOnBuild=true", "/p:BuildIpa=true"}
		require.Equal(t, desired, xbuild.buildCommandSlice())

		xbuild.SetCustomOptions("/nologo")
		desired = []string{constants.XbuildPath, "solution.sln", "/target:Build", "/p:SolutionDir=.", "/p:Configuration=Release", "/p:Platform=iPhone", "/p:ArchiveOnBuild=true", "/p:BuildIpa=true", "/nologo"}
		require.Equal(t, desired, xbuild.buildCommandSlice())
	}
}

func TestPrintableCommand(t *testing.T) {
	t.Log("it creates printable command")
	{
		xbuild := New("solution.sln")
		desired := fmt.Sprintf(`"%s" "solution.sln" "/p:SolutionDir=."`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetTarget("Build")
		desired = fmt.Sprintf(`"%s" "solution.sln" "/target:Build" "/p:SolutionDir=."`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetConfiguration("Release")
		desired = fmt.Sprintf(`"%s" "solution.sln" "/target:Build" "/p:SolutionDir=." "/p:Configuration=Release"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetPlatform("iPhone")
		desired = fmt.Sprintf(`"%s" "solution.sln" "/target:Build" "/p:SolutionDir=." "/p:Configuration=Release" "/p:Platform=iPhone"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetArchiveOnBuild(true)
		desired = fmt.Sprintf(`"%s" "solution.sln" "/target:Build" "/p:SolutionDir=." "/p:Configuration=Release" "/p:Platform=iPhone" "/p:ArchiveOnBuild=true"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetBuildIpa(true)
		desired = fmt.Sprintf(`"%s" "solution.sln" "/target:Build" "/p:SolutionDir=." "/p:Configuration=Release" "/p:Platform=iPhone" "/p:ArchiveOnBuild=true" "/p:BuildIpa=true"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())

		xbuild.SetCustomOptions("/nologo")
		desired = fmt.Sprintf(`"%s" "solution.sln" "/target:Build" "/p:SolutionDir=." "/p:Configuration=Release" "/p:Platform=iPhone" "/p:ArchiveOnBuild=true" "/p:BuildIpa=true" "/nologo"`, constants.XbuildPath)
		require.Equal(t, desired, xbuild.PrintableCommand())
	}
}
