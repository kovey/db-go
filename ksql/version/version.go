package version

import (
	"fmt"
	"strconv"
	"strings"
)

type VersionType int

const (
	NONE  VersionType = -1
	MAJOR VersionType = 3
	MINOR VersionType = 0
	BUILD VersionType = 7
)

func Version() string {
	return fmt.Sprintf("%d.%d.%d", MAJOR, MINOR, BUILD)
}

func Parse(v string) (VersionType, VersionType, VersionType) {
	info := strings.Split(v, ".")
	if len(info) != 3 {
		return NONE, NONE, NONE
	}

	verInfo := make([]VersionType, 3)
	for i, vi := range info {
		if val, err := strconv.Atoi(vi); err != nil {
			return NONE, NONE, NONE
		} else {
			verInfo[i] = VersionType(val)
		}
	}

	return verInfo[0], verInfo[1], verInfo[2]
}

func Compare(v string) int {
	major, minor, build := Parse(v)
	if major == NONE {
		return -1
	}

	if major < MAJOR {
		return -1
	} else if major > MAJOR {
		return 1
	}

	if minor < MINOR {
		return -1
	} else if minor > MINOR {
		return 1
	}

	if build < BUILD {
		return -1
	} else if build > BUILD {
		return 1
	}

	return 0
}
