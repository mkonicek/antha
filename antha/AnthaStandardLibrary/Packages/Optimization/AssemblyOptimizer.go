package Optimization

type AssemblyOptimizer interface {
	OptimizeAssembly(sequences []string, environment AssemblyOptimizerEnvironment, parameters AssemblyOptimizerParameters) AssemblyStrategy
}

type AssemblyOptimizerEnvironment struct {
	AssemblyMethods   []AssemblyMethod
	SequenceProviders []SequenceProvider
}

type SequenceProvider interface {
	Name() string
	MethodsImplemented() []string
	GetMethod(string) AssemblyMethod
}

type AssemblyOptimizerParameters interface {
	MethodName() string
	MaxIterations() int
	Set(string, interface{})
	GetInt(string) (int, bool)
	GetFloat(string) (float64, bool)
	GetString(string) (string, bool)
	GetBool(string) (bool, bool)
}

type AssemblyStep struct {
	Assemblies []Assembly
}

type Assembly struct {
	Method    AssemblyMethod
	Raw       []string
	Assembled []string
}

type AssemblyStrategy struct {
	Cost   int            // in pence
	Steps  []AssemblyStep // ordered
	Orders map[string]SequenceProvider
}
