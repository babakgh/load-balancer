package engine_test

// MockRequester is a simple implementation of Requester
type MockRequest struct {
	IDVal  string
	KeyVal string
}

func (m *MockRequest) ID() string {
	return m.IDVal
}

func (m *MockRequest) Key() string {
	return m.KeyVal
}
