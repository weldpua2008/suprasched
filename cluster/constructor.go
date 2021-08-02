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

// TerminatorTypeSpec is a constructor and a usage description for each ClustersTerminator type.
type TerminatorTypeSpec struct {
	constructor func(string) (ClustersTerminator, error)
	instance    func() ClustersTerminator
	Summary     string
	Description string
	Beta        bool
	Deprecated  bool
}

var (
	// DescriberConstructors is a map of all ClustersDescriber types with their specs.
	DescriberConstructors = map[string]DescriberTypeSpec{}
	// FetcherConstructors is a map of all ClustersFetcher types with their specs.
	FetcherConstructors = map[string]FetcherTypeSpec{}
	// TerminatorConstructors is a map of all ClustersTerminator types with their specs.
	TerminatorConstructors = map[string]TerminatorTypeSpec{}
)

// String constants representing each Fetcher type.
const (
	// ConstructorsFetcherTypeRest represents HTTP fetcher
	ConstructorsFetcherTypeRest     = "HTTP"
	ConstructorsDescriberTypeAwsEMR = "EMR"
	ConstructorsDescriberTypeRest   = "HTTP"
	ConstructorsTerminatorTypeEMR   = "EMR"
)
