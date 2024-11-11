package driver

import "strings"

var supports = map[string]bool{
	"mysql":      true,
	"postgresql": true,
	"oracle":     true,
	"sqlserver":  true,
	"db2":        true,
	"sqlite":     true,
	"firebird":   true,
	"h2":         true,
}

func SupportSavePoint(driverName string) bool {
	return supports[strings.ToLower(driverName)]
}
