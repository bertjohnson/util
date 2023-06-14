package crypto

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGenerateX509Certificate tests GenerateX509Certificate().
func TestGenerateX509Certificate(t *testing.T) {
	hostname, err := os.Hostname()
	assert.NoError(t, err)
	crtBytes, keyBytes, err := GenerateX509Certificate(strings.ToLower(hostname), hostname, time.Now(), 2*365*24*time.Hour, false, 2048, "p256")
	assert.NoError(t, err)
	assert.NotNil(t, crtBytes)
	assert.NotNil(t, keyBytes)
}
