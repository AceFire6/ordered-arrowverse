package build

// Environment enum is used to distinguish between different deployment targets like [Production] or [Development]
type Environment string

const (
	Unset Environment = ""
	// Development enum is the name given to the environment when running in local development
	Development Environment = "development"
	// Production enum is the name given to the environment when running in the remote production environment
	Production Environment = "production"
)

func (e Environment) String() string {
	return string(e)
}
