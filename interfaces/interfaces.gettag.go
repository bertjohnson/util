package interfaces

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// GetTagValue returns whether the tag exists. If so, its value.
func GetTagValue(tags string, targetTagName string) (tagValue string, exists bool, err error) {
	// Loop through tags.
	reader := bufio.NewReader(strings.NewReader(tags))
	inTagValue := false
	inTagTitle := true
	tagName := ""
	tagNameBuffer := bytes.Buffer{}
	tagValueBuffer := bytes.Buffer{}
	for {
		if c, _, err := reader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", false, err
			}
		} else {
			switch {
			case inTagValue:
				if c == '"' {
					inTagValue = false

					if tagName == targetTagName {
						tagValueParts := strings.SplitN(tagValueBuffer.String(), ",", 2)
						return tagValueParts[0], true, nil
					}

					tagValueBuffer.Reset()
				} else {
					_, err = tagValueBuffer.WriteRune(c)
					if err != nil {
						return "", false, err
					}
				}
			case inTagTitle:
				if c == ':' {
					inTagTitle = false
					tagName = strings.ToLower(strings.TrimSpace(tagNameBuffer.String()))
					tagNameBuffer.Reset()
				} else {
					_, err = tagNameBuffer.WriteRune(c)
					if err != nil {
						return "", false, err
					}
				}
			default:
				switch c {
				case '"':
					inTagValue = true
				case ' ':
					inTagTitle = true
				}
			}
		}
	}

	return "", false, nil
}
