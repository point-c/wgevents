package app

import (
	"os"
	"strings"
)

func OutputFilename(original string) string {
	return strings.TrimSuffix(original, ".go") + "_generated.go"
}

const (
	EnvDefaultPackageName = "wgevents"
	EnvDefaultGoFile      = "events.go"
)

const (
	EnvKeyPackageName = "GOPACKAGE"
	EnvKeyGoFile      = "GOFILE"
)

var (
	gofile    = EnvOrDefault(EnvKeyGoFile, EnvDefaultGoFile)
	gopackage = EnvOrDefault(EnvKeyPackageName, EnvDefaultPackageName)
)

func GoFile() string    { return gofile }
func GoPackage() string { return gopackage }

func EnvOrDefault(key, def string) string {
	return IfStringEmptyUseDefault(os.Getenv(key), def)
}

func IfStringEmptyUseDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
