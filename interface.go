package res

type Builder interface {
	Label() string
	Build() (any, error)
}
