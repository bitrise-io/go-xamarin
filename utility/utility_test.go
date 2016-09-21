package utility

import (
	"strings"
	"testing"

	"github.com/bitrise-tools/xamarin-builder/constants"
	"github.com/stretchr/testify/require"
)

func TestToConfig(t *testing.T) {
	t.Log("creates config from configuration and platform")
	{
		configuration := "Release"
		platform := "iPhone"
		config := ToConfig(configuration, platform)
		require.Equal(t, "Release|iPhone", config)
	}

	t.Log("creates config from configuration")
	{
		configuration := "Release"
		platform := ""
		config := ToConfig(configuration, platform)
		require.Equal(t, "Release|", config)
	}

	t.Log("creates config from platform")
	{
		configuration := ""
		platform := "iPhone"
		config := ToConfig(configuration, platform)
		require.Equal(t, "|iPhone", config)
	}

	t.Log("creates empty config")
	{
		configuration := ""
		platform := ""
		config := ToConfig(configuration, platform)
		require.Equal(t, "|", config)
	}
}

func TestFixWindowsPath(t *testing.T) {
	t.Log("fixes absolute windows path")
	{
		pth := `\bin\iPhoneSimulator\Debug`
		require.Equal(t, "/bin/iPhoneSimulator/Debug", FixWindowsPath(pth))
	}

	t.Log("fixes relative windows path")
	{
		pth := `bin\iPhoneSimulator\Debug`
		require.Equal(t, "bin/iPhoneSimulator/Debug", FixWindowsPath(pth))
	}

	t.Log("fixes relative windows path")
	{
		pth := `..\CreditCardValidator\CreditCardValidator.csproj`
		require.Equal(t, "../CreditCardValidator/CreditCardValidator.csproj", FixWindowsPath(pth))
	}

	t.Log("do not modify absolute unix path")
	{
		pth := "/bin/iPhoneSimulator/Debug"
		require.Equal(t, "/bin/iPhoneSimulator/Debug", FixWindowsPath(pth))
	}

	t.Log("do not modify relative unix path")
	{
		pth := "bin/iPhoneSimulator/Debug"
		require.Equal(t, "bin/iPhoneSimulator/Debug", FixWindowsPath(pth))
	}

	t.Log("do not modify relative unix path")
	{
		pth := "../CreditCardValidator/CreditCardValidator.csproj"
		require.Equal(t, "../CreditCardValidator/CreditCardValidator.csproj", FixWindowsPath(pth))
	}
}

func TestSplitAndStripList(t *testing.T) {
	t.Log("splits string list")
	{
		list := "{FEACFBD2-3405-455C-9665-78FE426C6842};{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}"
		split := SplitAndStripList(list, ";")
		require.Equal(t, 2, len(split))
		require.Equal(t, "{FEACFBD2-3405-455C-9665-78FE426C6842}", split[0])
		require.Equal(t, "{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}", split[1])
	}

	t.Log("splits string list")
	{
		list := "ARMv7, ARM64"
		split := SplitAndStripList(list, ",")
		require.Equal(t, 2, len(split))
		require.Equal(t, "ARMv7", split[0])
		require.Equal(t, "ARM64", split[1])
	}

	t.Log("do not split unless proper separator")
	{
		list := "ARMv7, ARM64"
		split := SplitAndStripList(list, ";")
		require.Equal(t, 1, len(split))
		require.Equal(t, "ARMv7, ARM64", split[0])
	}
}

func TestIsProjectType(t *testing.T) {
	t.Log("guid matches with project type")
	{
		guid := "EFBA0AD7-5A72-4C68-AF49-83D382785DCF"
		require.Equal(t, true, isProjectType(constants.XamarinAndroid, guid))
	}

	t.Log("guid does not match with project type")
	{
		guid := "E613F3A2-FE9C-494F-B74E-F63BCB86FEA6"
		require.Equal(t, false, isProjectType(constants.XamarinAndroid, guid))
	}
}

func TestIdetifyProjectType(t *testing.T) {
	t.Log("identifies XamarinAndroid")
	{
		guidList := "EFBA0AD7-5A72-4C68-AF49-83D382785DCF"
		require.Equal(t, constants.XamarinAndroid, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinAndroid")
	{
		guidList := "10368E6C-D01B-4462-8E8B-01FC667A7035"
		require.Equal(t, constants.XamarinAndroid, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinIos")
	{
		guidList := "E613F3A2-FE9C-494F-B74E-F63BCB86FEA6"
		require.Equal(t, constants.XamarinIos, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinIos")
	{
		guidList := "6BC8ED88-2882-458C-8E55-DFD12B67127B"
		require.Equal(t, constants.XamarinIos, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinIos")
	{
		guidList := "F5B4F3BC-B597-4E2B-B552-EF5D8A32436F"
		require.Equal(t, constants.XamarinIos, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinIos")
	{
		guidList := "FEACFBD2-3405-455C-9665-78FE426C6842"
		require.Equal(t, constants.XamarinIos, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinIos")
	{
		guidList := "8FFB629D-F513-41CE-95D2-7ECE97B6EEEC"
		require.Equal(t, constants.XamarinIos, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinIos")
	{
		guidList := "EE2C853D-36AF-4FDB-B1AD-8E90477E2198"
		require.Equal(t, constants.XamarinIos, IdetifyProjectType(guidList))
	}

	t.Log("identifies MonoMac")
	{
		guidList := "1C533B1C-72DD-4CB1-9F6B-BF11D93BCFBE"
		require.Equal(t, constants.MonoMac, IdetifyProjectType(guidList))
	}

	t.Log("identifies MonoMac")
	{
		guidList := "948B3504-5B70-4649-8FE4-BDE1FB46EC69"
		require.Equal(t, constants.MonoMac, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinMac")
	{
		guidList := "42C0BBD9-55CE-4FC1-8D90-A7348ABAFB23"
		require.Equal(t, constants.XamarinMac, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinMac")
	{
		guidList := "A3F8F2AB-B479-4A4A-A458-A89E7DC349F1"
		require.Equal(t, constants.XamarinMac, IdetifyProjectType(guidList))
	}

	t.Log("identifies XamarinTVOS")
	{
		guidList := "06FA79CB-D6CD-4721-BB4B-1BD202089C55"
		require.Equal(t, constants.XamarinTVOS, IdetifyProjectType(guidList))
	}

	t.Log("catches the first project type")
	{
		guidList := strings.Join([]string{
			"{10368E6C-D01B-4462-8E8B-01FC667A7035}", // XamarinAndroid
			"{EE2C853D-36AF-4FDB-B1AD-8E90477E2198}", // XamarinIos
			"{948B3504-5B70-4649-8FE4-BDE1FB46EC69}", // MonoMac
			"{A3F8F2AB-B479-4A4A-A458-A89E7DC349F1}", // XamarinMac
			"{06FA79CB-D6CD-4721-BB4B-1BD202089C55}", // XamarinTVOS
		}, ";")
		require.Equal(t, constants.XamarinAndroid, IdetifyProjectType(guidList))
	}
}
