package res

type Builder interface {
	Build() (any, error)
}
