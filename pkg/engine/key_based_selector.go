package engine

type KeyBasedSelector struct {
}

func NewKeyBasedSelector(total int64) *KeyBasedSelector {
	return &KeyBasedSelector{}
}

func (r *KeyBasedSelector) Select(request Requester) int64 {
	return 0
}
