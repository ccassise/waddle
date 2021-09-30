package mock

type MockWriter struct {
	Wrote []byte
}

func (w *MockWriter) Write(b []byte) (int, error) {
	w.Wrote = append(w.Wrote, b...)
	return len(b), nil
}
