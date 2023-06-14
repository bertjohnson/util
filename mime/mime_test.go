package mime

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseMultiPartMixed tests Parse() with a multipart/mixed email.
func TestParseMultiPartMixed(t *testing.T) {
	// Load the message bytes.
	message, err := ioutil.ReadFile("testdata/multipart-mixed.eml")
	assert.NoError(t, err)

	// Parse the message.
	headers, parts, err := Parse(bytes.NewReader(message))
	assert.NoError(t, err)
	expectedHeaders := map[string][]string{
		"content-type": {"multipart/mixed; boundary=5c5a8a19485f5e85f51cb58f063716ec97215e4df774d91ed40f6efdf7a3"},
		"date":         {"Sat, 30 Jun 2018 08:20:19 -0500"},
		"from":         {"administrator@mail.example.com"},
		"message-id":   {"<1530364819181864600.57176.4191787481235805846@Francium>"},
		"mime-version": {"1.0"},
		"subject":      {"Test Message"},
		"to":           {"895a7a8074db4237bfebd5221ab846c5@mail.example.com"},
	}
	assert.Equal(t, expectedHeaders, headers)
	assert.Equal(t, 2, len(parts))

	// Message body.
	messageBody, err := ioutil.ReadAll(parts[0].Reader)
	assert.NoError(t, err)
	assert.Equal(t, "Hello world", string(messageBody))

	// Attachment.
	attachmentBody, err := ioutil.ReadAll(parts[1].Reader)
	assert.NoError(t, err)
	assert.Equal(t, 1433, len(attachmentBody))
}

// TestParseTextHTML tests Parse() with a text/html email.
func TestParseTextHTML(t *testing.T) {
	// Load the message bytes.
	message, err := ioutil.ReadFile("testdata/text-html.eml")
	assert.NoError(t, err)

	// Parse the message.
	headers, parts, err := Parse(bytes.NewReader(message))
	assert.NoError(t, err)
	assert.Equal(t, 45, len(headers))
	assert.Equal(t, 1, len(parts))

	// Message body.
	messageBody, err := ioutil.ReadAll(parts[0].Reader)
	assert.NoError(t, err)
	assert.Equal(t, 2336, len(messageBody))
}
