package wizard

// Step constants
const (
	StepProjectName = 0
	StepModulePath  = 1
	StepPreset      = 2
	StepModules     = 3
	StepModuleOpts  = 4
	StepOutputDepth = 5
	StepConfirm     = 6
	TotalSteps      = 7
)

// OutputDepth describes how much boilerplate to generate
type OutputDepth int

const (
	DepthMinimal OutputDepth = iota
	DepthBoilerplate
	DepthFull
)

func (d OutputDepth) String() string {
	switch d {
	case DepthMinimal:
		return "Minimal"
	case DepthBoilerplate:
		return "Boilerplate"
	case DepthFull:
		return "Full example"
	default:
		return "Unknown"
	}
}

var outputDepths = []string{"Minimal", "Boilerplate", "Full example"}
