package strings

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBase58ToHex tests Base58ToHex().
func TestBase58ToHex(t *testing.T) {
	var input, result, expectedResult string
	var err error

	input = "JxF12TrwXzT5jvT"
	expectedResult = "48656c6c6f20776f726c64"
	result, err = Base58ToHex(input)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

// TestCapitalize tests Capitalize().
func TestCapitalize(t *testing.T) {
	var input, result, expectedResult string

	input = "processed"
	expectedResult = "Processed"
	result = Capitalize(input)
	assert.Equal(t, expectedResult, result)
}

// TestEnsureAlphanumeric tests EnsureAlphanumeric().
func TestEnsureAlphanumeric(t *testing.T) {
	var input string
	var expectedResult string
	var result string

	input = "1. Hydrogen, 2. Helium, 3. Lithium"
	expectedResult = "1Hydrogen2Helium3Lithium"
	result = EnsureAlphanumeric(input)
	assert.Equal(t, expectedResult, result, "Ensure alphanumeric.")
}

// TestExtractLinks tests ExtractLinks().
func TestExtractLinks(t *testing.T) {
	var input string
	var expectedResult []string
	var result []string
	var err error

	// Multiple values.
	input = "Let's go to https://google.com and [https://youtube.com] and {https://facebook.com} and (https://wikipedia.org) and <https://reddit.com>"
	expectedResult = []string{
		"https://google.com",
		"https://youtube.com",
		"https://facebook.com",
		"https://wikipedia.org",
		"https://reddit.com",
	}
	result, err = ExtractLinks(input)
	assert.NoError(t, err, "Extract links.")
	assert.Equal(t, expectedResult, result)

	// Starting and trailing values.
	input = "https://google.com and https://youtube.com"
	expectedResult = []string{
		"https://google.com",
		"https://youtube.com",
	}
	result, err = ExtractLinks(input)
	assert.NoError(t, err, "Extract links.")
	assert.Equal(t, expectedResult, result)

	// HTML values.
	input = "<a href=\"https://google.com\">Google</a> and <a href='https://youtube.com'>YouTube</a> and <a href=https://facebook.com>Facebook</a>"
	expectedResult = []string{
		"https://google.com",
		"https://youtube.com",
		"https://facebook.com",
	}
	result, err = ExtractLinks(input)
	assert.NoError(t, err, "Extract links.")
	assert.Equal(t, expectedResult, result)

	// Enclosed values.
	input = "[https://google.com], (https://youtube.com), <https://facebook.com>, {https://wikipedia.org}, <dc:identifier>https://reddit.com</dc:identifier>"
	expectedResult = []string{
		"https://google.com",
		"https://youtube.com",
		"https://facebook.com",
		"https://wikipedia.org",
		"https://reddit.com",
	}
	result, err = ExtractLinks(input)
	assert.NoError(t, err, "Extract links.")
	assert.Equal(t, expectedResult, result)

	// Complex enclosed values.
	input = "&lt;a href=&#34;https://golang.org/wiki/Go-Release-Cycle&#34; target=&#34;_blank&#34;&gt;usual schedule&lt;/a&gt;"
	expectedResult = []string{
		"https://golang.org/wiki/Go-Release-Cycle",
	}
	result, err = ExtractLinks(input)
	assert.NoError(t, err, "Extract links.")
	assert.Equal(t, expectedResult, result)

	// Malformed values.
	input = "[https://google.com (https://youtube.com <https://facebook.com {https://wikipedia.org"
	expectedResult = []string{
		"https://google.com",
		"https://youtube.com",
		"https://facebook.com",
		"https://wikipedia.org",
	}
	result, err = ExtractLinks(input)
	assert.NoError(t, err, "Extract links.")
	assert.Equal(t, expectedResult, result)

	// Empty value.
	input = ""
	expectedResult = []string{}
	result, err = ExtractLinks(input)
	assert.NoError(t, err, "Extract links.")
	assert.Equal(t, expectedResult, result)
}

// TestExtractLinksWithValidator tests ExtractLinksWithValidator().
func TestExtractLinksWithValidator(t *testing.T) {
	var input string
	var expectedResult []string
	var result []string
	var err error

	// Multiple values.
	input = "Relative files are a.png and b.jpeg and c.gif"
	expectedResult = []string{
		"a.png",
		"c.gif",
	}
	result, err = ExtractLinksWithValidator(input, func(input string) (bool, string) {
		if strings.HasSuffix(input, ".png") || strings.HasSuffix(input, ".gif") {
			return true, input
		}

		return false, ""
	})
	assert.NoError(t, err, "Extract links.")
	assert.Equal(t, expectedResult, result)
}

// TestFoldWhitespace tests FoldWhitespace().
func TestFoldWhitespace(t *testing.T) {
	var input string
	var expectedResult string
	var result string
	var err error

	// Test without lagging whitespace.
	input = "\tThere  is    \t a \rlot \nof \r\n\tvariety."
	expectedResult = "There is a lot of variety."
	result, err = FoldWhitespace(input)
	assert.NoError(t, err, "Fold without lagging whitespace.")
	assert.Equal(t, expectedResult, result, "Fold without lagging whitespace.")

	// Test with lagging whitespace.
	input = "\tThere  is    \t a \rlot \nof \r\n\tvariety.  "
	expectedResult = "There is a lot of variety."
	result, err = FoldWhitespace(input)
	assert.NoError(t, err, "Fold with lagging whitespace.")
	assert.Equal(t, expectedResult, result, "Fold with lagging whitespace.")
}

// TestFormatString tests FormatString().
func TestFormatString(t *testing.T) {
	var input string
	var replacements []string
	var expectedResult string
	var result string

	input = "Planets are %1 and %2."
	replacements = []string{"Mercury", "Venus", "Earth"}
	expectedResult = "Planets are Mercury and Venus."
	result = FormatString(input, replacements)
	assert.Equal(t, expectedResult, result, "Format string.")
}

// TestHexToBase58 tests HexToBase58().
func TestHexToBase58(t *testing.T) {
	var input, result, expectedResult string
	var err error

	input = "4865-6c6c-6f20-776f-726c64"
	expectedResult = "JxF12TrwXzT5jvT"
	result, err = HexToBase58(input)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

// TestIsAlphanumeric tests IsAlphanumeric().
func TestIsAlphanumeric(t *testing.T) {
	var input string

	// Test match.
	input = "1Mercury2Venus3Earth"
	assert.True(t, IsAlphanumeric(input), "Alphanumeric.")

	// Test mismatch.
	input = "1: Mercury, 2: Venus, 3: Earth"
	assert.False(t, IsAlphanumeric(input), "Not alphanumeric.")
}

// TestIsWhitespace tests IsWhitespace().
func TestIsWhitespace(t *testing.T) {
	// All known whitespace runes.
	whitespaceOrds := map[int]bool{
		9:        true,
		10:       true,
		11:       true,
		13:       true,
		32:       true,
		160:      true,
		49824:    true,
		14785152: true,
		14844032: true,
		14844033: true,
		14844034: true,
		14844035: true,
		14844036: true,
		14844037: true,
		14844038: true,
		14844039: true,
		14844040: true,
		14844041: true,
		14844042: true,
		14844043: true,
		14844079: true,
		14844319: true,
		14909568: true,
	}

	// Test first 10,000 runes.
	for i := 0; i < 10000; i++ {
		result := IsWhitespace(rune(i))
		expectedResult := whitespaceOrds[i]
		assert.Equal(t, expectedResult, result, "Is whitespace.")
	}

	// Test all whitespace characters.
	for ord := range whitespaceOrds {
		assert.True(t, IsWhitespace(rune(ord)), "Is whitespace.")
	}
}

// TestNormalizeCharSet tests NormalizeCharSet().
func TestNormalizeCharSet(t *testing.T) {
	var input string
	var expectedResult string
	var result string

	// Replacement test for CP1252.
	input = "ansi_x3.110-1983"
	expectedResult = "iso-8859-1"
	result = NormalizeCharSet(input)
	assert.Equal(t, expectedResult, result, "Normalized character set.")

	// Replacement test for CP1252.
	input = "CP1252"
	expectedResult = "windows-1252"
	result = NormalizeCharSet(input)
	assert.Equal(t, expectedResult, result, "Normalized character set.")

	// As-is test.
	input = "none"
	expectedResult = "none"
	result = NormalizeCharSet(input)
	assert.Equal(t, expectedResult, result, "Normalized character set.")
}

// TestNormalizePhone tests NormalizePhone().
func TestNormalizePhone(t *testing.T) {
	var input string
	var expectedResult string
	var result string
	var err error

	// Simple test.
	input = "(800) 867-5309"
	expectedResult = "8008675309"
	result, err = NormalizePhone(input)
	assert.NoError(t, err, "Normalized phone number.")
	assert.Equal(t, expectedResult, result, "Normalized phone number.")

	// Complex test.
	input = "+1 (800) 867-[5309] x001"
	expectedResult = "18008675309x001"
	result, err = NormalizePhone(input)
	assert.NoError(t, err, "Normalized phone number.")
	assert.Equal(t, expectedResult, result, "Normalized phone number.")
}

// TestPermutations tests Permutations().
func TestPermutations(t *testing.T) {
	var input []string
	var expectedResultCount int
	var resultCount int

	input = []string{"aa", "bb", "cc", "dd"}
	expectedResultCount = 24
	resultCount = len(Permutations(input))
	assert.Equal(t, expectedResultCount, resultCount, "Determined permutations.")
}

// TestRemoveAny tests RemoveAny().
func TestRemoveAny(t *testing.T) {
	var haystack, needles string
	var expectedResult string
	var result string
	var err error

	// No vowels test.
	haystack = "xxyyzz"
	needles = "aeiou"
	expectedResult = "xxyyzz"
	result, err = RemoveAny(haystack, needles)
	assert.NoError(t, err, "Removed needle characters.")
	assert.Equal(t, expectedResult, result, "Removed needle characters.")

	// Vowels test.
	haystack = "abcdefghijklmnopqrstuvwxyz"
	needles = "aeiou"
	expectedResult = "bcdfghjklmnpqrstvwxyz"
	result, err = RemoveAny(haystack, needles)
	assert.NoError(t, err, "Removed needle characters.")
	assert.Equal(t, expectedResult, result, "Removed needle characters.")
}

// TestRemoveUnreadableCharacters tests RemoveUnreadableCharacters().
func TestRemoveUnreadableCharacters(t *testing.T) {
	var input string
	var expectedResult string
	var result string
	var err error

	// No unreadable characters test.
	input = "do re mi fa so             la ti do"
	expectedResult = "do re mi fa so la ti do"
	result, err = RemoveUnreadableCharacters(input)
	assert.NoError(t, err, "Removed unreadable characters.")
	assert.Equal(t, expectedResult, result, "Removed unreadable characters.")
}

// TestRemoveWhitespace tests RemoveWhitespace().
func TestRemoveWhitespace(t *testing.T) {
	var input string
	var expectedResult string
	var result string
	var err error

	input = "\tThere  is    \t a \rlot \nof \r\n\tvariety."
	expectedResult = "Thereisalotofvariety."
	result, err = RemoveWhitespace(input)
	assert.NoError(t, err, "Remove whitespace.")
	assert.Equal(t, expectedResult, result, "Remove whitespace.")
}

// TestReturnAllBetween tests ReturnAllBetween().
func TestReturnAllBetween(t *testing.T) {
	var haystack, startNeedle, endNeedle string
	var expectedResult []string
	var result []string

	// Simple test with matches.
	haystack = "doremifasolatido"
	startNeedle = "o"
	endNeedle = "i"
	expectedResult = []string{"rem", "lat"}
	result = ReturnAllBetween(haystack, startNeedle, endNeedle)
	assert.Equal(t, expectedResult, result)

	// Simple test without matches.
	haystack = "doremifasolatido"
	startNeedle = "gold"
	endNeedle = "silver"
	expectedResult = []string{}
	result = ReturnAllBetween(haystack, startNeedle, endNeedle)
	assert.Equal(t, expectedResult, result)
}

// TestReturnAllBetweenArray tests ReturnAllBetweenArray().
func TestReturnAllBetweenArray(t *testing.T) {
	var haystack string
	var startNeedles, endNeedles, expectedResult, result []string

	// Simple test with matches.
	haystack = "ab[cd{ef(gh)(ij)}kl]"
	startNeedles = []string{"{", "[", "("}
	endNeedles = []string{"}", "]", ")"}
	expectedResult = []string{"cd{ef(gh", "ij"}
	result = ReturnAllBetweenArray(haystack, startNeedles, endNeedles)
	assert.Equal(t, expectedResult, result)

	// Simple test without matches.
	haystack = "abcdefg"
	startNeedles = []string{"a", "b", "c"}
	endNeedles = []string{"x", "y", "z"}
	expectedResult = []string{}
	result = ReturnAllBetweenArray(haystack, startNeedles, endNeedles)
	assert.Equal(t, expectedResult, result)
}

// TestReturnAllBetweenInsensitive tests ReturnAllBetweenInsensitive().
func TestReturnAllBetweenInsensitive(t *testing.T) {
	var haystack, startNeedle, endNeedle string
	var expectedResult []string
	var result []string

	// Simple test with matches.
	haystack = "DoREmiFasOLatIDo"
	startNeedle = "o"
	endNeedle = "I"
	expectedResult = []string{"REm", "Lat"}
	result = ReturnAllBetweenInsensitive(haystack, startNeedle, endNeedle)
	assert.Equal(t, expectedResult, result)

	// Simple test without matches.
	haystack = "doremifasolatido"
	startNeedle = "gold"
	endNeedle = "silver"
	expectedResult = []string{}
	result = ReturnAllBetweenInsensitive(haystack, startNeedle, endNeedle)
	assert.Equal(t, expectedResult, result)
}

// TestReturnAndRemoveAllBetween tests ReturnAndRemoveAllBetween().
func TestReturnAndRemoveAllBetween(t *testing.T) {
	var haystack, startNeedle, endNeedle string
	var expectedOutput, output string
	var expectedResult, result []string

	// HTML test, inclusive.
	haystack = "<i>This</i> <b>is</b> <a>a</a> <s>simple</s> <u>test</u>"
	startNeedle = "<"
	endNeedle = ">"
	expectedOutput = "This is a simple test"
	expectedResult = []string{"i", "/i", "b", "/b", "a", "/a", "s", "/s", "u", "/u"}
	output, result = ReturnAndRemoveAllBetween(haystack, startNeedle, endNeedle, true)
	assert.Equal(t, expectedOutput, output, "HTML test.")
	assert.Equal(t, expectedResult, result)

	// Parentheses test, inclusive.
	haystack = "Some words are (big)... other words are (BIGGER)"
	startNeedle = "("
	endNeedle = ")"
	expectedOutput = "Some words are ()... other words are ()"
	expectedResult = []string{"big", "BIGGER"}
	output, result = ReturnAndRemoveAllBetween(haystack, startNeedle, endNeedle, false)
	assert.Equal(t, expectedOutput, output, "Parentheses test.")
	assert.Equal(t, expectedResult, result)
}

// TestReturnBetween tests ReturnBetween().
func TestReturnBetween(t *testing.T) {
	var haystack, startNeedle, endNeedle string
	var expectedResult string
	var result string

	// Simple test with matches.
	haystack = "HHeLi"
	startNeedle = "H"
	endNeedle = "Li"
	expectedResult = "He"
	result = ReturnBetween(haystack, startNeedle, endNeedle)
	assert.Equal(t, expectedResult, result, "Returned between.")

	// Simple test without matches.
	haystack = "HHeLi"
	startNeedle = "B"
	endNeedle = "C"
	expectedResult = ""
	result = ReturnBetween(haystack, startNeedle, endNeedle)
	assert.Equal(t, expectedResult, result, "Returned between.")
}

// TestReturnBetweenInsensitive tests ReturnBetweenInsensitive().
func TestReturnBetweenInsensitive(t *testing.T) {
	var haystack, startNeedle, endNeedle string
	var expectedResult string
	var result string

	// Simple test with matches.
	haystack = "HheLI"
	startNeedle = "h"
	endNeedle = "li"
	expectedResult = "he"
	result = ReturnBetweenInsensitive(haystack, startNeedle, endNeedle)
	assert.Equal(t, expectedResult, result, "Returned between, case insensitive.")

	// Simple test without matches.
	haystack = "HheLI"
	startNeedle = "B"
	endNeedle = "C"
	expectedResult = ""
	result = ReturnBetweenInsensitive(haystack, startNeedle, endNeedle)
	assert.Equal(t, expectedResult, result, "Returned between, case insensitive.")
}

// TestSplitAny tests SplitAny().
func TestSplitAny(t *testing.T) {
	var haystack, needle string
	var expectedResult []string
	var result []string

	// Simple test.
	haystack = "Mercury;Venus,Earth,Mars;Jupiter"
	needle = ",;"
	expectedResult = []string{"Mercury", "Venus", "Earth", "Mars", "Jupiter"}
	result = SplitAny(haystack, needle)
	assert.Equal(t, expectedResult, result)
}

// TestLowestNonNegative tests lowestNonNegative().
func TestLowest(t *testing.T) {
	min, index := lowestNonNegative([]int{25, -1, 41, -1, -1, 9})
	assert.Equal(t, 9, min, "Lowest.")
	assert.Equal(t, 5, index, "Lowest.")
}
