package engine

type Backender interface {
	ID() string
	Key() string
	Process(request Requester) error
}
