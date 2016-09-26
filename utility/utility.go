package utility

import (
	"fmt"
	"strings"

	"github.com/bitrise-tools/go-xamarin/constants"
)

// ToConfig ...
func ToConfig(configuration, platform string) string {
	return fmt.Sprintf("%s|%s", configuration, platform)
}

// FixWindowsPath ...
func FixWindowsPath(pth string) string {
	return strings.Replace(pth, `\`, "/", -1)
}

// SplitAndStripList ...
func SplitAndStripList(list, separator string) []string {
	split := strings.Split(list, separator)
	elements := []string{}
	for _, s := range split {
		elements = append(elements, strings.TrimSpace(s))
	}
	return elements
}

func isProjectType(projectType constants.XamarinProjectType, guid string) bool {
	typeGuids := constants.ProjectTypeGUIDMap[projectType]
	for _, typeGUID := range typeGuids {
		if typeGUID == guid {
			return true
		}
	}
	return false
}

// IdetifyProjectType ...
func IdetifyProjectType(typeGUIDList string) constants.XamarinProjectType {
	typeGUIDs := SplitAndStripList(typeGUIDList, ";")
	for _, typeGUID := range typeGUIDs {
		guid := strings.Trim(typeGUID, "{")
		guid = strings.Trim(guid, "}")

		if isProjectType(constants.XamarinAndroid, guid) {
			return constants.XamarinAndroid
		} else if isProjectType(constants.XamarinIos, guid) {
			return constants.XamarinIos
		} else if isProjectType(constants.MonoMac, guid) {
			return constants.MonoMac
		} else if isProjectType(constants.XamarinMac, guid) {
			return constants.XamarinMac
		} else if isProjectType(constants.XamarinTVOS, guid) {
			return constants.XamarinTVOS
		}
	}
	return constants.Unknown
}
