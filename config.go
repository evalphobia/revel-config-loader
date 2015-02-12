/*
	this package is a config loader for revel
	see: https://github.com/evalphobia/revel-config-loader
*/

package loader

import (
	"github.com/revel/revel"
	"os"
	"strings"
)

var (
	configs     map[string]*revel.MergedConfig
	debugParams map[string]string // this allows to override parameters when the debug build
)

// load config file from the name.
func LoadConfig(filename string) *revel.MergedConfig {
	separator := "/"
	if os.IsPathSeparator('\\') {
		separator = "\\"
	} else {
		separator = "/"
	}

	// add .conf extension
	if !strings.HasSuffix(filename, ".conf") {
		filename += ".conf"
	}
	c, err := revel.LoadConfig(revel.RunMode + separator + filename)
	if err != nil {
		// fallback default setting
		c, err = revel.LoadConfig(filename)
	}
	if err != nil {
		revel.ERROR.Printf("error load config: %s, %s", filename, err)
	}
	return c
}

// get parameter from the config file
func GetConfigValueDefault(conf, section, key, df string) string {
	if debugParams != nil {
		param, ok := debugParams[strings.Join([]string{conf, section, key}, "_")]
		if ok {
			return param
		}
	}
	if configs == nil {
		configs = make(map[string]*revel.MergedConfig)
	}

	if configs[conf] == nil {
		configs[conf] = LoadConfig(conf)
	}
	configs[conf].SetSection(section)
	return configs[conf].StringDefault(key, df)
}
