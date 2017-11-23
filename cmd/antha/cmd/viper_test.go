package cmd

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestGetStringSliceWorkaroundNeeded(t *testing.T) {
	t.Skip("workaround not needed")

	f := pflag.NewFlagSet("", 0)
	f.StringSlice("empty", nil, "")
	v := viper.New()
	if err := v.BindPFlags(f); err != nil {
		t.Fatal(err)
	}
	if s := v.GetStringSlice("empty"); len(s) == 0 {
		t.Errorf("cmd.GetStringSlice() may not be needed: %q", s)
	} else if len(s) != 1 || s[0] != "[]" {
		t.Errorf("cmd.GetStringSlice() needs to be improved: %q", s)
	}
}

func TestGetStringSlice(t *testing.T) {
	f := pflag.NewFlagSet("", 0)
	f.StringSlice("empty", nil, "")
	f.StringSlice("one", []string{"one"}, "")
	v := viper.New()
	if err := v.BindPFlags(f); err != nil {
		t.Fatal(err)
	}

	if s := v.GetStringSlice("empty"); len(s) != 0 {
		t.Errorf("cmd.GetStringSlice() workaround may be needed: %q", s)
	} else if s := v.GetStringSlice("one"); len(s) != 1 {
		t.Errorf("cmd.GetStringSlice() workaround may be needed: %q", s)
	} else if len(s) != 1 || s[0] != "one" {
		t.Errorf("cmd.GetStringSlice() workaround may be needed: %q", s)
	}
}
