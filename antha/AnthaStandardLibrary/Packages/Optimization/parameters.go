package Optimization

type genericMapParameters map[string]interface{}

func (gmp genericMapParameters) GetInt(s string) (int, bool) {
	in, ok := gmp[s]

	if !ok {
		return 0, false
	}

	i, ok := in.(int)

	if !ok {
		return 0, false
	}
	return i, true
}

func (gmp genericMapParameters) GetFloat(s string) (float64, bool) {
	in, ok := gmp[s]

	if !ok {
		return 0.0, false
	}

	f, ok := in.(float64)

	if !ok {
		return 0.0, false
	}

	return f, true
}
func (gmp genericMapParameters) GetString(s string) (string, bool) {
	in, ok := gmp[s]

	if !ok {
		return "", false
	}

	s2, ok := in.(string)

	if !ok {
		return "", false
	}

	return s2, true
}
func (gmp genericMapParameters) GetBool(s string) (bool, bool) {
	in, ok := gmp[s]

	if !ok {
		return false, false
	}

	b, ok := in.(bool)

	if !ok {
		return false, false
	}
	return b, true
}

func (gmp genericMapParameters) Set(s string, v interface{}) {
	gmp[s] = v
}

type GAParameters struct {
	genericMapParameters
}

func NewGAParameters() GAParameters {
	m := make(map[string]interface{})
	return GAParameters{genericMapParameters: m}
}

func (gap GAParameters) MethodName() string {
	return "GA"
}
func (gap GAParameters) MaxIterations() int {
	it, _ := gap.GetInt("max_iterations")
	return it
}
