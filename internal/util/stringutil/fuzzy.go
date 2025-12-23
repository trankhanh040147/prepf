package stringutil

import "strings"

// MinFuzzyMatchScore is the minimum similarity score required for a fuzzy match
const MinFuzzyMatchScore = 0.3

// FuzzyMatch finds the closest matching string using simple string similarity
// Returns the closest match, or empty string if no reasonable match found
func FuzzyMatch(input string, validKeys []string) string {
	inputLower := strings.ToLower(input)

	// Check for substring matches (prefix/suffix) - highest priority
	for _, key := range validKeys {
		keyLower := strings.ToLower(key)
		if strings.HasPrefix(keyLower, inputLower) || strings.HasPrefix(inputLower, keyLower) {
			return key
		}
	}

	// Use simple similarity scoring for fuzzy matching
	bestMatch := ""
	bestScore := 0.0

	for _, key := range validKeys {
		score := Similarity(inputLower, strings.ToLower(key))
		if score > bestScore {
			bestScore = score
			bestMatch = key
		}
	}

	// Only return a match if similarity is reasonable
	if bestScore > MinFuzzyMatchScore {
		return bestMatch
	}

	return ""
}

// Similarity calculates a simple similarity score between two strings
// Returns a value between 0.0 and 1.0, where 1.0 is identical
func Similarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	// Calculate longest common subsequence ratio
	lcs := LongestCommonSubsequence(s1, s2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	if maxLen == 0 {
		return 0.0
	}

	return float64(lcs) / float64(maxLen)
}

// LongestCommonSubsequence calculates the length of the longest common subsequence
func LongestCommonSubsequence(s1, s2 string) int {
	m, n := len(s1), len(s2)
	if m == 0 || n == 0 {
		return 0
	}

	lcsTable := make([][]int, m+1)
	for i := range lcsTable {
		lcsTable[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				lcsTable[i][j] = lcsTable[i-1][j-1] + 1
			} else {
				if lcsTable[i-1][j] > lcsTable[i][j-1] {
					lcsTable[i][j] = lcsTable[i-1][j]
				} else {
					lcsTable[i][j] = lcsTable[i][j-1]
				}
			}
		}
	}

	return lcsTable[m][n]
}
