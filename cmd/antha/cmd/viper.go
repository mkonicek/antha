package cmd

import (
	"strings"

	"github.com/spf13/viper"
)

func getStringSliceWorkaround(key string) []string {
	l := viper.GetStringSlice(key)
	if len(l) != 1 {
		panic("cmd.GetStringSlice: unexpected length")
	}
	e := l[0]
	if len(e) < 2 || e[0] != '[' || e[len(e)-1] != ']' {
		panic("cmd.GetStringSlice: unexpected length")
	} else if len(e) == 2 {
		return nil
	}

	return strings.Split(e[1:len(e)-1], ",")
}

// GetStringSlice gets a slice of values for a key.
//
// It is workaround for a bug in viper: string slice flags like:
//  ["a", "b"]
// are flattened to a string slice like this:
//  ["[a,b]"]
func GetStringSlice(key string) []string {
	if true {
		return viper.GetStringSlice(key)
	}
	return getStringSliceWorkaround(key)
}
