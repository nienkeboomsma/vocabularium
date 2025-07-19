package collatinus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkBySentence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "one sentence",
			input:    "Arma virumque cano.",
			expected: []string{"Arma virumque cano."},
		},
		{
			name:     "one sentence without a full stop",
			input:    "Arma virumque cano",
			expected: []string{"Arma virumque cano."},
		},
		{
			name:  "three sentences",
			input: "Arma virumque cano. Troiae qui primus ab oris Italiam, fato profugus, Laviniaque venit litora. Multum ille et terris iactatus et alto vi superum saevae memorem Iunonis ob iram.",
			expected: []string{
				"Arma virumque cano.",
				"Troiae qui primus ab oris Italiam, fato profugus, Laviniaque venit litora.",
				"Multum ille et terris iactatus et alto vi superum saevae memorem Iunonis ob iram.",
			},
		},
		{
			name:  "three sentences with excess spacing",
			input: "   Arma virumque cano.       Troiae qui primus ab oris Italiam, fato profugus, Laviniaque venit litora.    Multum ille et terris iactatus et alto vi superum saevae memorem Iunonis ob iram.",
			expected: []string{
				"Arma virumque cano.",
				"Troiae qui primus ab oris Italiam, fato profugus, Laviniaque venit litora.",
				"Multum ille et terris iactatus et alto vi superum saevae memorem Iunonis ob iram.",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, _ := chunkBySentence(test.input)
			assert.Equal(t, test.expected, output)
		})
	}
}
