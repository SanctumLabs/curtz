package encoding

import (
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"log"

	"github.com/sanctumlabs/curtz/app/pkg"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"

	"github.com/google/uuid"
	nanoid "github.com/matoous/go-nanoid"
)

//Encode encodes the uuid to a base64 string that is url-safe.
func Encode(id uuid.UUID) string {
	return b64.RawURLEncoding.EncodeToString(id.NodeID())
}

//Decode decodes a base64 string to a raw uuid.
func Decode(id string) (uuid.UUID, error) {
	dec, err := b64.RawURLEncoding.DecodeString(id)

	if err != nil {
		return uuid.UUID{}, err
	}

	decoded, err := uuid.FromBytes(dec)
	if err != nil {
		return uuid.UUID{}, err
	}

	return decoded, nil
}

//GetUniqueShortCode returns a random unique short code.
func GetUniqueShortCode() (string, error) {
	l := pkg.ShortCodeLength
	alphanumeric := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return nanoid.Generate(alphanumeric, l)
}

//GenUniqueID returns a random but unique id.
func GenUniqueID() identifier.ID {
	return identifier.New()
}

//GenHexKey generates a crypto-random key with byte length len and hex-encodes it to a string.
func GenHexKey(len int) string {
	bytes := make([]byte, len)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(bytes)
}
