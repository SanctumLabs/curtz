package encoding

import (
	b64 "encoding/base64"
	"testing"
)

func Test_Encode(t *testing.T) {
	id := GenUniqueID()
	want := b64.RawURLEncoding.EncodeToString(id.NodeID())
	got := Encode(id)

	if got != want {
		t.Fatalf("got %v want %v", got, want)
	}
}
