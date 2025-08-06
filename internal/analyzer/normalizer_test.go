package analyzer

import "testing"

func TestNameNormalizer_NormalizeName(t *testing.T) {
	normalizer := NewNameNormalizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple ascii name",
			input:    "John Doe",
			expected: "john doe",
		},
		{
			name:     "name with diacritics",
			input:    "Martin Pražák",
			expected: "martin prazak",
		},
		{
			name:     "name with punctuation",
			input:    "martin.prazak",
			expected: "martin prazak",
		},
		{
			name:     "name with mixed case and punctuation",
			input:    "Martin.Prazak",
			expected: "martin prazak",
		},
		{
			name:     "name with multiple spaces",
			input:    "Martin   Prazak",
			expected: "martin prazak",
		},
		{
			name:     "name with leading/trailing spaces",
			input:    "  Martin Prazak  ",
			expected: "martin prazak",
		},
		{
			name:     "name with underscores",
			input:    "martin_prazak",
			expected: "martin prazak",
		},
		{
			name:     "complex case with everything",
			input:    "  Martin.Pražák_123  ",
			expected: "martin prazak 123",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "mixed diacritics",
			input:    "José María",
			expected: "jose maria",
		},
		{
			name:     "german umlaut",
			input:    "Björn Müller",
			expected: "bjorn muller",
		},
		{
			name:     "french accents",
			input:    "François Léon",
			expected: "francois leon",
		},
		{
			name:     "numbers preserved",
			input:    "User123",
			expected: "user123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizer.NormalizeName(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNameNormalizer_IdenticalNormalization(t *testing.T) {
	normalizer := NewNameNormalizer()

	// Test that different variations of the same name normalize to the same result
	variations := []string{
		"Martin Prazak",
		"martin prazak",
		"MARTIN PRAZAK",
		"Martin.Prazak",
		"martin.prazak",
		"Martin_Prazak",
		"Martin  Prazak",
		"  Martin Prazak  ",
		"Martin Pražák",
		"martin.pražák",
	}

	expected := "martin prazak"
	for _, variation := range variations {
		result := normalizer.NormalizeName(variation)
		if result != expected {
			t.Errorf("NormalizeName(%q) = %q, want %q", variation, result, expected)
		}
	}
}