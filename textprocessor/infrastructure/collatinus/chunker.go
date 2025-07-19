package collatinus

import (
	"errors"
	"strings"
)

func chunkBySentence(data string) ([]string, error) {
	parts := strings.Split(data, ".")
	chunks := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if len(part) > 4000 {
			return []string{}, errors.New("maximum supported sentence length is 4,000 characters; sentences that exceed this must be broken up by a full stop")
		}

		if len(part) > 0 {
			chunk := part + "."
			chunks = append(chunks, chunk)
		}
	}

	return chunks, nil
}
