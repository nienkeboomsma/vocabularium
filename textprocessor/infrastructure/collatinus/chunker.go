package collatinus

import (
	"strings"
)

func chunkBySentence(data string) []string {
	parts := strings.Split(data, ".")
	chunks := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if len(part) > 0 {
			chunk := part + "."
			chunks = append(chunks, chunk)
		}
	}

	return chunks
}
