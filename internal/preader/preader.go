package preader

import "fmt"

type PReader struct {
	data map[string]any
}

func New(data map[string]any) *PReader {
	return &PReader{data}
}

func (p *PReader) GetStringRequired(k string) (string, error) {
	v, ok := p.data[k]
	if !ok {
		return "", fmt.Errorf("parameter %s is required", k)
	}

	if str, ok := v.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("")
}
