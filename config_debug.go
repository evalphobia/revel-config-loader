// +build debug tests

package loader

import (
	"strings"
)

// allows to override parameters
func SetDebugParameter(conf, section, key, value string) {
	if debugParams == nil {
		debugParams = make(map[string]string)
	}
	debugParams[strings.Join([]string{conf, section, key}, "_")] = value
}
