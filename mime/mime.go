// Package mime contains a method for parsing MIME.
package mime

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
	"mime/quotedprintable"
	"strings"
)

// MIMEPart defines a MIME part.
type MIMEPart struct {
	Headers map[string][]string // MIME part headers.
	Reader  io.Reader           // Body reader.
}

// Parse parses a MIME message into its parts.
func Parse(reader io.Reader) (headers map[string][]string, parts []MIMEPart, err error) {
	headers = make(map[string][]string)
	parts = []MIMEPart{}
	err = parseHelper(reader, headers, &parts, true)
	if err != nil {
		return nil, nil, err
	}

	return headers, parts, nil
}

// parseHelper parses a MIME message recursively.
func parseHelper(reader io.Reader, headers map[string][]string, parts *[]MIMEPart, root bool) error {
	// If not working with the message root, track local headers.
	if !root {
		headers = make(map[string][]string)
	}

	// Loop through message.
	buffer := bytes.Buffer{}
	bufReader := bufio.NewReader(reader)
	folded := false
	var lastC rune
headerLoop:
	for {
		if c, _, err := bufReader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		} else {
			if lastC == '\n' {
				// Ignore folded spaces following new lines.
				if c == '\t' || c == ' ' {
					lastC = c
					if !folded {
						_, err = buffer.WriteRune(' ')
						if err != nil {
							return err
						}
						folded = true
					}
					continue headerLoop
				}

				// Process the buffer.
				header := buffer.String()
				headerParts := strings.SplitN(header, ":", 2)
				if len(headerParts) == 2 {
					headerKey := strings.ToLower(strings.TrimSpace(headerParts[0]))
					if _, ok := headers[headerKey]; ok {
						headers[headerKey] = append(headers[headerKey], strings.TrimSpace(headerParts[1]))
					} else {
						headers[headerKey] = []string{strings.TrimSpace(headerParts[1])}
					}
				}
				folded = false
				buffer.Reset()
			}

			// Process special characters.
			switch c {
			case '\r':
				// Ignore carriage returns.
				continue headerLoop
			case '\n':
				if lastC == '\n' {
					break headerLoop
				}
			default:
				_, err = buffer.WriteRune(c)
				if err != nil {
					return err
				}
			}
			lastC = c
		}
	}
	// Process the buffer.
	header := buffer.String()
	headerParts := strings.SplitN(header, ":", 2)
	if len(headerParts) == 2 {
		headerKey := strings.ToLower(strings.TrimSpace(headerParts[0]))
		if _, ok := headers[headerKey]; ok {
			headers[headerKey] = append(headers[headerKey], strings.TrimSpace(headerParts[1]))
		} else {
			headers[headerKey] = []string{strings.TrimSpace(headerParts[1])}
		}
	}
	buffer.Reset()

	// Validate that the message is MIME-encoded.
	validMime := true
	var boundaryName string
	boundaryNameArray := headers["content-type"]
	if len(boundaryNameArray) > 0 {
		boundaryName = boundaryNameArray[0]
		boundaryPos := strings.Index(boundaryName, "boundary=")
		if boundaryPos > -1 {
			boundaryName = strings.TrimSpace(boundaryName[boundaryPos+9:])
			if len(boundaryName) > 0 {
				_, mimeVersionOk := headers["mime-version"]
				if !mimeVersionOk {
					validMime = false
				}
			} else {
				validMime = false
			}
		} else {
			validMime = false
		}
	} else {
		validMime = false
	}

	// If the body is not encoded as MIME, return as-is.
	if !validMime {
		// Handle encodings.
		if contentTypeHeader, ok := headers["content-transfer-encoding"]; ok {
			switch strings.ToLower(contentTypeHeader[0]) {
			case "base64":
				reader = base64.NewDecoder(base64.StdEncoding, bufReader)
			case "quoted-printable":
				reader = quotedprintable.NewReader(bufReader)
			default:
				reader = bufReader
			}
		}

		part := MIMEPart{
			Headers: headers,
			Reader:  reader,
		}
		*parts = append(*parts, part)
		return nil
	}

	// Determine the MIME boundary name.
	boundaryNameLength := len(boundaryName)

	// Parse the MIME parts.
	initialized := false
	inBoundaryName := false
	immediatelyAfterBoundaryName := false
	var partBuffer = bytes.Buffer{}
	twoCago := '\n'
	threeCago := '\n'
bodyLoop:
	for {
		if c, _, err := bufReader.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		} else {
			if inBoundaryName {
				_, err = buffer.WriteRune(c)
				if err != nil {
					return err
				}
				if buffer.Len() == boundaryNameLength {
					if buffer.String() == boundaryName {
						if initialized {
							// Part is complete. Process it recursively.
							err = parseHelper(bufio.NewReader(bytes.NewReader(partBuffer.Bytes())), headers, parts, false)
							if err != nil {
								return err
							}
							partBuffer.Reset()
						} else {
							initialized = true
						}

						lastC = '\n'
						inBoundaryName = false
						immediatelyAfterBoundaryName = true
					}
				}
			} else {
				if c == '-' && lastC == '-' && twoCago == '\n' {
					// If we've reached the closing boundary, stop processing.
					if immediatelyAfterBoundaryName {
						break bodyLoop
					}

					// Truncate boundary characters.
					partBufferLength := partBuffer.Len()
					if partBufferLength > 2 {
						if threeCago == '\r' {
							partBuffer.Truncate(partBufferLength - 3)
						} else {
							partBuffer.Truncate(partBufferLength - 2)
						}
					} else {
						partBuffer.Reset()
					}

					inBoundaryName = true
					buffer.Reset()
				} else {
					if immediatelyAfterBoundaryName {
						if c == '\r' || c == '\n' {
							continue bodyLoop
						}
						immediatelyAfterBoundaryName = false
					}
					_, err = partBuffer.WriteRune(c)
					if err != nil {
						return err
					}
				}
				threeCago = twoCago
				twoCago = lastC
				lastC = c
			}
		}
	}

	return nil
}
