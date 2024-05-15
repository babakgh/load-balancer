package engine

type Selector interface {
	Select(request Requester) int64
}
