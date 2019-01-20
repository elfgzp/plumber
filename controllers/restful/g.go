package restful

var (
	DefaultPageLimit int
	MaxPageLimit     int
)

func init() {
	DefaultPageLimit = 10
	MaxPageLimit = 100
}
