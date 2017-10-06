package cmd

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestGetStringSliceWorkaroundNeeded(t *testing.T) {
	f := pflag.NewFlagSet("", 0)
	f.StringSlice("empty", nil, "")
	v := viper.New()
	v.BindPFlags(f)
	if s := v.GetStringSlice("empty"); len(s) == 0 {
		t.Errorf("cmd.GetStringSlice() may not be needed: %q", s)
	} else if len(s) != 1 || s[0] != "[]" {
		t.Errorf("cmd.GetStringSlice() needs to be improved: %q", s)
	}
}
