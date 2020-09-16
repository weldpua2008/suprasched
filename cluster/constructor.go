package cluster

// FetcherTypeSpec is a constructor and a usage description for each ClustersFetcher type.
type FetcherTypeSpec struct {
	constructor func(string) (ClustersFetcher, error)
	instance    func() ClustersFetcher
	Summary     string
	Description string
	Beta        bool
	Deprecated  bool
}

// FetcherConstructors is a map of all Communicator types with their specs.
var FetcherConstructors = map[string]FetcherTypeSpec{}

// String constants representing each Fetcher type.
const (
	// ConstructorsFetcherTypeRest represents HTTP fetcher
	ConstructorsFetcherTypeRest = "HTTP"
)
