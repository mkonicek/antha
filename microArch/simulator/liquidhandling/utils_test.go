package liquidhandling

import "testing"

func TestSummariseChannels(t *testing.T) {

	tests := [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{0, 1, 2, 4, 5, 6, 7},
		{2, 3, 5, 6, 8},
		{0, 2, 4, 6},
		{3},
	}

	expected := []string{
		"channels 0-7",
		"channels 0-2,4-7",
		"channels 2-3,5-6,8",
		"channels 0,2,4,6",
		"channel 3",
	}

	for i := 0; i < len(tests); i++ {
		if g, e := summariseChannels(tests[i]), expected[i]; g != e {
			t.Errorf("test %d: got \"%s\", expected \"%s\"", i, g, e)
		}
	}
}
