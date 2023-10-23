package version

import (
	"fmt"
)

var Version string = "v1.0.1"

func GetVersion() (version string) {
	version = Version
	return version
}

func GetMsgVersion() (msgVersion string) {
	msgVersion = fmt.Sprintf(
		"version = %s",
		Version)
	return msgVersion
}
