package communicator

// TypeSpec is a constructor and a usage description for each Communicator type.
type TypeSpec struct {
	constructor func(string) (Communicator, error)
	instance    func() Communicator
	Summary     string
	Description string
	Beta        bool
	Deprecated  bool
}

// Constructors is a map of all Communicator types with their specs.
var Constructors = map[string]TypeSpec{}

// String constants representing each communicator type.
const (
	// ConstructorsTypeRest represents HTTP communicator
	ConstructorsTypeRest = "HTTP"
)
