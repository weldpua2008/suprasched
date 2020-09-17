package job

// JobTypeSpec is a constructor and a usage description for each JobsFetcher type.
type JobTypeSpec struct {
	constructor func(string) (JobsFetcher, error)
	instance    func() JobsFetcher
	Summary     string
	Description string
	Beta        bool
	Deprecated  bool
}

// Constructors is a map of all JobFetcher types with their specs.
var Constructors = map[string]JobTypeSpec{}

// String constants representing each communicator type.
const (
	// ConstructorsJobsFetcherRest represents Rest fetcher.
	ConstructorsJobsFetcherRest = "HTTP"
)
