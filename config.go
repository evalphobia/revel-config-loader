/*
	this package is a config loader for revel
	see: https://github.com/evalphobia/revel-config-loader
*/

package loader

import (
	"flag"
	"go/build"
	"os"
	"strings"

	"github.com/revel/revel"
)

const (
	overrideDir = "override"
)

var (
	configs     map[string]*revel.MergedConfig
	debugParams map[string]string // this allows to override parameters when the debug build
	searchPath  = flag.String("revelconf", "", "")
)

// load config file from the name.
func LoadConfig(filename string) *revel.MergedConfig {
	if *searchPath != "" {
		AddSearchPath(*searchPath)
		*searchPath = ""
	}

	c, err := loadPriorConfig(filename)
	if err != nil {
		revel.ERROR.Printf("error load config: %s", err.Error())
		revel.ERROR.Printf("filename: %s, runmode: %s, path: %s", filename, revel.RunMode, revel.ConfPaths)
	}
	return c
}

func loadPriorConfig(filename string) (*revel.MergedConfig, error) {
	filename = addSuffix(filename)

	// load overrided config
	if revel.RunMode != "test" {
		c, err := loadOverrideConfig(filename)
		if err == nil {
			return c, err
		}
	}

	// load enviromental config
	c, err := revel.LoadConfig(revel.RunMode + getSeparator() + filename)
	if err == nil {
		return c, err
	}

	// load root config
	return revel.LoadConfig(filename)
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
		if configs[conf] == nil {
			return ""
		}
	}
	configs[conf].SetSection(section)
	return configs[conf].StringDefault(key, df)
}

// add an optional search path
func AddSearchPath(path string) {
	gopaths := strings.Split(build.Default.GOPATH, ":")
	for _, gopath := range gopaths {
		paths := strings.Split(path, ":")
		for _, p := range paths {
			revel.ConfPaths = append(revel.ConfPaths, gopath+"/src/"+p)
		}
	}
}

func getSeparator() string {
	if os.IsPathSeparator('\\') {
		return "\\"
	}
	return "/"
}

func addSuffix(file string) string {
	if !strings.HasSuffix(file, ".conf") {
		file += ".conf"
	}
	return file
}

func loadOverrideConfig(filename string) (*revel.MergedConfig, error) {
	return revel.LoadConfig(overrideDir + getSeparator() + filename)
}
