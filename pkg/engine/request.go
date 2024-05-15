package engine

type Requester interface {
	ID() string
	Key() string
}
