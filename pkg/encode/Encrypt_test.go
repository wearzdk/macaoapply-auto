package encode

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestEncode(t *testing.T) {
	res := Encode("https://www.rymoe.com", "Hello World!")
	assert.Equal(t, res, "yhEBCli0+r0gr4NJ")
}
