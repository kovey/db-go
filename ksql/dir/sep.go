package dir

import "runtime"

func Sep() string {
	if runtime.GOOS != "windows" {
		return "/"
	}

	return "\\"
}
