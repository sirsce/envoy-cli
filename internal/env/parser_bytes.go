package env

import (
	"bytes"
)

// ParseBytes parses a byte slice of env content and returns a key-value map.
// It reuses the existing Parser logic by wrapping the bytes in a reader.
func (p *Parser) ParseBytes(data []byte) (map[string]string, error) {
	return p.ParseReader(bytes.NewReader(data))
}

// ToBytes serialises a key-value map to env file bytes using the Writer.
func (w *Writer) ToBytes(entries map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	if err := w.WriteMap(&buf, entries); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
