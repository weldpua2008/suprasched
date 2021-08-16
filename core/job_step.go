package core

// JobStep is a description of a job step.
type JobStep struct {
	// command
	CMD string
	// the lower priority is - the sooner it will executed
	// equal priority steps are executed ib random order
	StepPriority int
	// CmdEnv stores environment variables
	CmdEnv []string
	// RunAs defines user
	RunAs string
	// UseSHELL defines if we should wrap the command with shell
	UseSHELL bool
	// TTR defines Time-to-run in Millisecond
	TTR uint64
	// IgnoreExitCode if true
	IgnoreExitCode bool
}

type JobStepsSorter []JobStep

func (s JobStepsSorter) Len() int           { return len(s) }
func (s JobStepsSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s JobStepsSorter) Less(i, j int) bool { return s[i].StepPriority < s[j].StepPriority }
