package encoding

import (
	"encoding/base64"
	"testing"
)

func TestEncode(t *testing.T) {
	id := GenUniqueID()
	expected := base64.RawURLEncoding.EncodeToString(id.Bytes())
	actual := Encode(id)

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
