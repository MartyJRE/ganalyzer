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
			expected: "johndoe",
		},
		{
			name:     "name with diacritics",
			input:    "Martin Pražák",
			expected: "martinprazak",
		},
		{
			name:     "name with punctuation",
			input:    "martin.prazak",
			expected: "martinprazak",
		},
		{
			name:     "name with mixed case and punctuation",
			input:    "Martin.Prazak",
			expected: "martinprazak",
		},
		{
			name:     "name with multiple spaces",
			input:    "Martin   Prazak",
			expected: "martinprazak",
		},
		{
			name:     "name with leading/trailing spaces",
			input:    "  Martin Prazak  ",
			expected: "martinprazak",
		},
		{
			name:     "name with underscores",
			input:    "martin_prazak",
			expected: "martinprazak",
		},
		{
			name:     "complex case with everything",
			input:    "  Martin.Pražák_123  ",
			expected: "martinprazak123",
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
			expected: "josemaria",
		},
		{
			name:     "german umlaut",
			input:    "Björn Müller",
			expected: "bjornmuller",
		},
		{
			name:     "french accents",
			input:    "François Léon",
			expected: "francoisleon",
		},
		{
			name:     "numbers preserved",
			input:    "User123",
			expected: "user123",
		},
		{
			name:     "czech diacritics with ě",
			input:    "Michal Pěkný",
			expected: "michalpekny",
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

	expected := "martinprazak"
	for _, variation := range variations {
		result := normalizer.NormalizeName(variation)
		if result != expected {
			t.Errorf("NormalizeName(%q) = %q, want %q", variation, result, expected)
		}
	}
}

func TestNameNormalizer_MichalPeknyVariants(t *testing.T) {
	normalizer := NewNameNormalizer()

	// Test that all Michal Pekny variations normalize to the same result
	variations := []string{
		"Michal Pěkný",
		"Michal Pekny",
		"michal.pekny",
		"michalpekny",
		"michal pekny",
		"MICHAL PEKNY",
	}

	expected := "michalpekny"
	for _, variation := range variations {
		result := normalizer.NormalizeName(variation)
		if result != expected {
			t.Errorf("NormalizeName(%q) = %q, want %q", variation, result, expected)
		}
	}
}
