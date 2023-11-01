package internal

import (
	"strings"
)

func AssetIndex(version string) (index string) {
	if version == "1.7.10" {
		return version
	}

	l := strings.Split(version, ".")

	return strings.Join(l[:len(l)-1], ".")
}
