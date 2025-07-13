package collatinus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitise(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "sentence with linebreaks",
			input:    []byte("Arma virumque cano, Troiae qui primus ab oris\nItaliam, fato profugus, Laviniaque venit\rlitora, multum ille et terris iactatus et alto\nvi superum saevae memorem Iunonis ob iram."),
			expected: "Arma virumque cano, Troiae qui primus ab oris Italiam, fato profugus, Laviniaque venit litora, multum ille et terris iactatus et alto vi superum saevae memorem Iunonis ob iram.",
		},
		{
			name:     "sentence with null-byte",
			input:    []byte("Arma virumque cano, Troiae q\x00ui primus ab oris Italiam, fato profugus, Laviniaque venit litora, multum ille et terris iactatus et alto vi superum saevae memorem Iunonis ob iram."),
			expected: "Arma virumque cano, Troiae qui primus ab oris Italiam, fato profugus, Laviniaque venit litora, multum ille et terris iactatus et alto vi superum saevae memorem Iunonis ob iram.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := sanitise(test.input)
			assert.Equal(t, test.expected, output)
		})
	}
}
