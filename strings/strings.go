// Package strings provides functions for manipulating strings.
package strings

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"strconv"
	"strings"

	"github.com/shengdoushi/base58"
)

// Base58ToHex converts a base58-encoded string to hex.
func Base58ToHex(input string) (output string, err error) {
	// Decode hex.
	inputBytes, err := base58.Decode(input, base58.BitcoinAlphabet)
	if err != nil {
		return "", err
	}

	// Encode hex.
	return hex.EncodeToString(inputBytes), nil
}

// Capitalize capitalizes the first letter of a string.
func Capitalize(input string) string {
	// Validate input.
	if input == "" {
		return input
	}

	return strings.ToUpper(input[0:1]) + input[1:]
}

// EnsureAlphanumeric removes any non-alphanumeric characters from a string.
func EnsureAlphanumeric(input string) string {
	output := bytes.Buffer{}
	inputLength := len(input)
	for i := 0; i < inputLength; i++ {
		r := rune(input[i])
		charCode := int(r)
		if (charCode >= 48 && charCode <= 57) || (charCode >= 65 && charCode <= 90) || (charCode >= 97 && charCode <= 122) {
			output.WriteRune(r)
		}
	}
	return output.String()
}

// ExtractLinks returns links from a plain/text body.
func ExtractLinks(haystack string) ([]string, error) {
	return ExtractLinksWithValidator(haystack, func(input string) (bool, string) {
		protocolPos := strings.Index(input, "://")
		if protocolPos > 0 {
			// Find non-letter prefixes.
			truncated := false
			prefix := ""
			for i := protocolPos - 1; i >= 0; i-- {
				charCode := int(rune(input[i]))
				if charCode < 65 || (charCode > 90 && charCode < 97) || charCode > 122 {
					prefix = input[0:i]
					input = input[i:]
					truncated = true
					break
				}
			}

			// Remove enclosing characters.
			matched := true
			for matched {
				switch {
				case strings.HasPrefix(input, "<"), strings.HasPrefix(input, "="):
					pos := strings.Index(input, ">")
					if pos > -1 {
						input = input[1:pos]
					} else {
						input = input[1:]
					}
				case strings.HasPrefix(input, ">"):
					pos := strings.Index(input, "<")
					if pos > -1 {
						input = input[1:pos]
					} else {
						input = input[1:]
					}
				case strings.HasPrefix(input, "("):
					pos := strings.Index(input, ")")
					if pos > -1 {
						input = input[1:pos]
					} else {
						input = input[1:]
					}
				case strings.HasPrefix(input, "["):
					pos := strings.Index(input, "]")
					if pos > -1 {
						input = input[1:pos]
					} else {
						input = input[1:]
					}
				case strings.HasPrefix(input, "{"):
					pos := strings.Index(input, "}")
					if pos > -1 {
						input = input[1:pos]
					} else {
						input = input[1:]
					}
				default:
					if truncated {
						pos := strings.Index(input[1:], string(input[0:1]))
						if pos > -1 {
							input = input[1 : pos+1]
						} else {
							input = input[1:]
						}
					}
					matched = false
				}
				truncated = false
			}

			// Remove suffix if it matches prefix.
			if prefix != "" {
				prefixLength := len(prefix)
			prefixLoop:
				for j := prefixLength - 1; j >= 0; j-- {
					switch prefix[j] {
					case '=', ':':
						prefix = prefix[j+1:]
						break prefixLoop
					}
				}

				prefixLength = len(prefix)
				if prefixLength > 0 {
					if strings.HasSuffix(input, prefix) {
						input = input[0 : len(input)-prefixLength]
					}
				}
			}

			return true, input
		}

		return false, ""
	})
}

// ExtractLinksWithValidator returns links from a plain/text body, using a custom function to validate links.
func ExtractLinksWithValidator(haystack string, isLink func(input string) (bool, string)) ([]string, error) {
	// Handle empty string.
	if haystack == "" {
		return []string{}, nil
	}

	reader := bufio.NewReader(strings.NewReader(haystack))
	outputBuffer := bytes.Buffer{}
	links := []string{}
	var c rune
	var err error
	for {
		if c, _, err = reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		} else {
			if IsWhitespace(c) {
				if islink, url := isLink(outputBuffer.String()); islink {
					links = append(links, url)
				}

				outputBuffer.Reset()
				if err != nil {
					return nil, err
				}
			} else {
				_, err = outputBuffer.WriteRune(c)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if outputBuffer.Len() > 0 {
		if islink, url := isLink(outputBuffer.String()); islink {
			links = append(links, url)
		}
	}

	return links, nil
}

// FoldWhitespace normalizes whitespace.
func FoldWhitespace(input string) (string, error) {
	var buffer bytes.Buffer

	reader := bufio.NewReader(strings.NewReader(input))
	inWhitespace := true
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		} else {
			// Avoid redundant whitespace.
			isWhitespace := IsWhitespace(c)
			if isWhitespace {
				if inWhitespace {
					continue
				}
				inWhitespace = true
				buffer.WriteRune(' ')
			} else {
				inWhitespace = false
				buffer.WriteRune(c)
			}
		}
	}

	if inWhitespace {
		substring := buffer.String()
		if len(substring) > 0 {
			return substring[:len(substring)-1], nil
		}
	}
	return buffer.String(), nil
}

// FormatString formats a string, replacing % numeric expressions with their corresponding array values.
func FormatString(input string, replacements []string) string {
	replacementsLength := len(replacements)
	for i := 1; i <= replacementsLength; i++ {
		input = strings.Replace(input, "%"+strconv.Itoa(i), replacements[i-1], -1)
	}

	return input
}

// HexToBase58 converts a hex-encoded string to base58.
func HexToBase58(input string) (output string, err error) {
	input = strings.Replace(input, "-", "", -1)

	// Decode hex.
	inputBytes, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}

	// Encode base58.
	return base58.Encode(inputBytes, base58.BitcoinAlphabet), nil
}

// IsAlphanumeric returns true if the string consists of only letters and numbers.
func IsAlphanumeric(input string) bool {
	inputLength := len(input)
	for i := 0; i < inputLength; i++ {
		r := rune(input[i])
		charCode := int(r)
		if !((charCode >= 48 && charCode <= 57) || (charCode >= 65 && charCode <= 90) || (charCode >= 97 && charCode <= 122)) {
			return false
		}
	}
	return true
}

// IsWhitespace returns true if the rune is whitespace.
func IsWhitespace(r rune) bool {
	switch r {
	case 9: // Horizontal tab.
		return true
	case 10: // Line feed.
		return true
	case 11: // Vertical tab.
		return true
	case 13: // Carriage return.
		return true
	case 32: // Space.
		return true
	case 160: // Non-breaking space.
		return true
	case 49824:
		return true
	case 14785152:
		return true
	case 14844032:
		return true
	case 14844033:
		return true
	case 14844034:
		return true
	case 14844035:
		return true
	case 14844036:
		return true
	case 14844037:
		return true
	case 14844038:
		return true
	case 14844039:
		return true
	case 14844040:
		return true
	case 14844041:
		return true
	case 14844042:
		return true
	case 14844043:
		return true
	case 14844079:
		return true
	case 14844319:
		return true
	case 14909568:
		return true
	}

	return false
}

// NormalizeCharSet replaces aliased character sets with their more common representation.
func NormalizeCharSet(charSet string) string {
	switch strings.ToLower(charSet) {
	case "ansi_x3.110-1983":
		return "iso-8859-1"
	case "cp1252":
		return "windows-1252"
	}
	return charSet
}

// NormalizePhone formats phone numbers consistently.
func NormalizePhone(input string) (string, error) {
	var buffer bytes.Buffer

	reader := bufio.NewReader(strings.NewReader(input))
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		} else {
			switch c {
			// Eliminate unnecessary formatting characters
			case '+', '(', ')', '[', ']', '{', '}', '<', '>', '-', '_', '.', ' ', '\t':
			default:
				buffer.WriteRune(c)
			}
		}

	}

	return buffer.String(), nil
}

// Permutations returns all permutations made through combinations of input nodes.
func Permutations(input []string) [][]string {
	var helper func([]string, int)
	output := [][]string{}

	helper = func(input []string, n int) {
		if n == 1 {
			tmp := make([]string, len(input))
			copy(tmp, input)
			output = append(output, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(input, n-1)
				if n%2 == 1 {
					tmp := input[i]
					input[i] = input[n-1]
					input[n-1] = tmp
				} else {
					tmp := input[0]
					input[0] = input[n-1]
					input[n-1] = tmp
				}
			}
		}
	}
	helper(input, len(input))
	return output
}

// RemoveAny removes any matching needle characters.
func RemoveAny(haystack string, needles string) (string, error) {
	runes := make(map[rune]bool)
	for _, c := range needles {
		runes[c] = true
	}

	var buffer bytes.Buffer
	reader := strings.NewReader(haystack)
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		} else {
			if _, ok := runes[c]; !ok {
				buffer.WriteRune(c)
			}
		}
	}

	return buffer.String(), nil
}

// RemoveUnreadableCharacters removes non-human readable characters.
func RemoveUnreadableCharacters(input string) (string, error) {
	var buffer bytes.Buffer

	reader := bufio.NewReader(strings.NewReader(input))
	inWhitespace := true
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		} else {
			// Avoid redundant whitespace.
			isWhitespace := IsWhitespace(c)
			if isWhitespace {
				if inWhitespace {
					continue
				}
				inWhitespace = true
				buffer.WriteRune(' ')
			} else {
				if c >= 32 && c < 127 {
					buffer.WriteRune(c)
				}
				inWhitespace = false
			}
		}
	}

	return strings.Trim(buffer.String(), " \t\r\n"), nil
}

// RemoveWhitespace removes whitespace.
func RemoveWhitespace(input string) (string, error) {
	var buffer bytes.Buffer

	reader := bufio.NewReader(strings.NewReader(input))
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		} else {
			// Avoid redundant whitespace.
			if !IsWhitespace(c) {
				buffer.WriteRune(c)
			}
		}
	}

	return buffer.String(), nil
}

// ReturnAllBetween returns all substrings in haystack between startNeedle and endNeedle.
func ReturnAllBetween(haystack string, startNeedle string, endNeedle string) []string {
	results := []string{}

	prevStartPos := 0
	startPos := 0
	startNeedleLength := len(startNeedle)
	endNeedleLength := len(endNeedle)
	for startPos > -1 {
		prevStartPos = startPos
		startPos = strings.Index(haystack[startPos:], startNeedle)
		if startPos > -1 {
			endPos := strings.Index(haystack[prevStartPos+startPos+startNeedleLength:], endNeedle)
			if endPos > -1 {
				results = append(results, haystack[prevStartPos+startPos+startNeedleLength:prevStartPos+startPos+startNeedleLength+endPos])
				startPos = prevStartPos + startPos + startNeedleLength + endPos + endNeedleLength
			} else {
				startPos = -1
			}
		}
	}
	return results
}

// ReturnAllBetweenArray returns all substrings in haystack between entries of startNeedles and entries of endNeedles.
func ReturnAllBetweenArray(haystack string, startNeedles []string, endNeedles []string) []string {
	results := []string{}

	haystackLength := len(haystack)
	endNeedlesLength := len(endNeedles)
	endPositions := make([]int, endNeedlesLength)
iLoop:
	for i := 0; i < haystackLength; i++ {
		window := haystack[i:]
		for _, startNeedle := range startNeedles {
			if strings.HasPrefix(window, startNeedle) {
				startNeedleLength := len(startNeedle)
				for j := 0; j < endNeedlesLength; j++ {
					endPositions[j] = strings.Index(window[startNeedleLength:], endNeedles[j])
				}

				endPos, index := lowestNonNegative(endPositions)
				if endPos > -1 {
					results = append(results, window[startNeedleLength:startNeedleLength+endPos])

					i = i + startNeedleLength + endPos + len(endNeedles[index]) - 1
					continue iLoop
				}
			}
		}
	}

	return results
}

// ReturnAllBetweenInsensitive returns all substrings in haystack between startNeedle and endNeedle, case insensitive.
func ReturnAllBetweenInsensitive(haystack string, startNeedle string, endNeedle string) []string {
	haystackLowercase := strings.ToLower(haystack)
	results := []string{}

	startNeedleLowercase := strings.ToLower(startNeedle)
	startNeedleLength := len(startNeedle)
	endNeedleLowercase := strings.ToLower(endNeedle)
	endNeedleLength := len(endNeedle)

	prevStartPos := 0
	startPos := 0
	for startPos > -1 {
		prevStartPos = startPos
		startPos = strings.Index(haystackLowercase[startPos:], startNeedleLowercase)
		if startPos > -1 {
			endPos := strings.Index(haystackLowercase[prevStartPos+startPos+startNeedleLength:], endNeedleLowercase)
			if endPos > -1 {
				results = append(results, haystack[prevStartPos+startPos+startNeedleLength:prevStartPos+startPos+startNeedleLength+endPos])
				startPos = prevStartPos + startPos + startNeedleLength + endPos + endNeedleLength
			} else {
				startPos = -1
			}
		}
	}
	return results
}

// ReturnAndRemoveAllBetween returns all substrings in haystack between startNeedle and endNeedle, and removes them.
func ReturnAndRemoveAllBetween(haystack string, startNeedle string, endNeedle string, removeNeedles bool) (string, []string) {
	outputBuffer := bytes.Buffer{}
	results := []string{}

	prevStartPos := 0
	startPos := 0
	startNeedleLength := len(startNeedle)
	endNeedleLength := len(endNeedle)
	for startPos > -1 {
		prevStartPos = startPos
		startPos = strings.Index(haystack[startPos:], startNeedle)
		if startPos > -1 {
			outputBuffer.WriteString(haystack[prevStartPos : prevStartPos+startPos])
			if !removeNeedles {
				outputBuffer.WriteString(startNeedle)
			}

			endPos := strings.Index(haystack[prevStartPos+startPos+startNeedleLength:], endNeedle)
			if endPos > -1 {
				if !removeNeedles {
					outputBuffer.WriteString(endNeedle)
				}

				results = append(results, haystack[prevStartPos+startPos+startNeedleLength:prevStartPos+startPos+startNeedleLength+endPos])
				startPos = prevStartPos + startPos + startNeedleLength + endPos + endNeedleLength
			} else {
				outputBuffer.WriteString(haystack[startPos+startNeedleLength:])
				startPos = -1
			}
		} else {
			outputBuffer.WriteString(haystack[prevStartPos:])
		}
	}
	return outputBuffer.String(), results
}

// ReturnBetween returns the substring in haystack between startNeedle and endNeedle.
func ReturnBetween(haystack string, startNeedle string, endNeedle string) string {
	startPos := strings.Index(haystack, startNeedle)
	startNeedleLength := len(startNeedle)
	if startPos > -1 {
		endPos := strings.Index(haystack[startPos+startNeedleLength:], endNeedle)
		if endPos > -1 {
			return haystack[startPos+startNeedleLength : startPos+startNeedleLength+endPos]
		}
	}
	return ""
}

// ReturnBetweenInsensitive returns the substring in haystack between startNeedle and endNeedle, case insensitive.
func ReturnBetweenInsensitive(haystack string, startNeedle string, endNeedle string) string {
	haystackLowercase := strings.ToLower(haystack)
	startPos := strings.Index(haystackLowercase, startNeedle)
	startNeedleLength := len(startNeedle)
	if startPos > -1 {
		endPos := strings.Index(haystackLowercase[startPos+startNeedleLength:], endNeedle)
		if endPos > -1 {
			return haystack[startPos+startNeedleLength : startPos+startNeedleLength+endPos]
		}
	}
	return ""
}

// SplitAny splits haystack on any matching needle characters.
func SplitAny(haystack string, needles string) []string {
	runes := make(map[rune]bool)
	for _, c := range needles {
		runes[c] = true
	}

	var buffer bytes.Buffer
	reader := strings.NewReader(haystack)
	output := []string{}
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil
			}
		} else {
			if _, ok := runes[c]; ok {
				if buffer.Len() > 0 {
					output = append(output, buffer.String())
				}
				buffer.Reset()
			} else {
				buffer.WriteRune(c)
			}
		}
	}
	if buffer.Len() > 0 {
		output = append(output, buffer.String())
	}

	return output
}

// lowestNonNegative returns the lowest index of the collection that is greater than -1.
func lowestNonNegative(indexes []int) (int, int) {
	min := -1
	minIndex := -1
	for i, index := range indexes {
		if index > -1 && (min == -1 || index < min) {
			min = index
			minIndex = i
		}
	}

	return min, minIndex
}
