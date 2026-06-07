package markup

import "errors"

type MockConverter struct {
	ShouldFail bool
}

func (m *MockConverter) Convert(content []byte, ctx RenderContext) ([]byte, error) {
	if m.ShouldFail {
		return nil, errors.New("fail to convert: forced test error")
	}
	return []byte("<p>mock converted HTML</p>"), nil
}
