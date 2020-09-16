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

// DescriberTypeSpec is a constructor and a usage description for each ClustersDescriber type.
type DescriberTypeSpec struct {
	constructor func(string) (ClustersDescriber, error)
	instance    func() ClustersDescriber
	Summary     string
	Description string
	Beta        bool
	Deprecated  bool
}

var (
	// DescriberConstructors is a map of all Communicator types with their specs.
	DescriberConstructors = map[string]DescriberTypeSpec{}

	// FetcherConstructors is a map of all Communicator types with their specs.
	FetcherConstructors = map[string]FetcherTypeSpec{}
)

// String constants representing each Fetcher type.
const (
	// ConstructorsFetcherTypeRest represents HTTP fetcher
	ConstructorsFetcherTypeRest     = "HTTP"
	ConstructorsDescriberTypeAwsEMR = "EMR"
	ConstructorsDescriberTypeRest   = "HTTP"
)
