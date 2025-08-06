package analyzer

import (
	"regexp"
	"strings"
)

// NameNormalizer handles contributor name normalization
type NameNormalizer struct {
	punctuationRegex  *regexp.Regexp
	whitespaceRegex   *regexp.Regexp
	diacriticReplacer *strings.Replacer
}

// TODO: Michal Pekny is not working correctly

// NewNameNormalizer creates a new name normalizer
func NewNameNormalizer() *NameNormalizer {
	// Common diacritic mappings - extend as needed
	diacriticMappings := []string{
		"á", "a", "à", "a", "â", "a", "ä", "a", "ā", "a", "ã", "a", "å", "a",
		"é", "e", "è", "e", "ê", "e", "ë", "e", "ē", "e", "ę", "e",
		"í", "i", "ì", "i", "î", "i", "ï", "i", "ī", "i",
		"ó", "o", "ò", "o", "ô", "o", "ö", "o", "ō", "o", "õ", "o",
		"ú", "u", "ù", "u", "û", "u", "ü", "u", "ū", "u",
		"ý", "y", "ÿ", "y",
		"ñ", "n",
		"ç", "c", "č", "c",
		"š", "s", "ž", "z", "ř", "r", "ď", "d", "ť", "t", "ň", "n",
		// Add uppercase versions
		"Á", "A", "À", "A", "Â", "A", "Ä", "A", "Ā", "A", "Ã", "A", "Å", "A",
		"É", "E", "È", "E", "Ê", "E", "Ë", "E", "Ē", "E", "Ę", "E",
		"Í", "I", "Ì", "I", "Î", "I", "Ï", "I", "Ī", "I",
		"Ó", "O", "Ò", "O", "Ô", "O", "Ö", "O", "Ō", "O", "Õ", "O",
		"Ú", "U", "Ù", "U", "Û", "U", "Ü", "U", "Ū", "U",
		"Ý", "Y", "Ÿ", "Y",
		"Ñ", "N",
		"Ç", "C", "Č", "C",
		"Š", "S", "Ž", "Z", "Ř", "R", "Ď", "D", "Ť", "T", "Ň", "N",
	}

	return &NameNormalizer{
		punctuationRegex:  regexp.MustCompile(`[^\p{L}\p{N}\s]+`),
		whitespaceRegex:   regexp.MustCompile(`\s+`),
		diacriticReplacer: strings.NewReplacer(diacriticMappings...),
	}
}

// NormalizeName normalizes a contributor name by:
// 1. Removing diacritics/accents using a replacer
// 2. Converting to lowercase
// 3. Removing punctuation (dots, underscores, etc.)
// 4. Normalizing whitespace (multiple spaces to single space, trim)
func (nn *NameNormalizer) NormalizeName(name string) string {
	if name == "" {
		return name
	}

	// Step 1: Remove diacritics
	normalized := nn.diacriticReplacer.Replace(name)

	// Step 2: Convert to lowercase
	normalized = strings.ToLower(normalized)

	// Step 3: Remove punctuation (keep only letters, numbers, and spaces)
	normalized = nn.punctuationRegex.ReplaceAllString(normalized, " ")

	// Step 4: Normalize whitespace
	normalized = strings.TrimSpace(normalized)
	normalized = nn.whitespaceRegex.ReplaceAllString(normalized, " ")

	return normalized
}
