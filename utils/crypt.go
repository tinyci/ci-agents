package utils

import (
	"strconv"
	"strings"

	"github.com/tinyci/ci-agents/errors"
)

// ParseCryptKey takes a hex-encoded string and turns it into a byte stream, typically used for encryption keys.
func ParseCryptKey(key string) ([]byte, *errors.Error) {
	key = strings.TrimSpace(key)

	if len(key) == 0 || len(key)%2 != 0 {
		return nil, errors.New("invalid key -- not a multiple of two in length")
	}

	ret := []byte{}

	for i := 0; i < len(key); i += 2 {
		i, err := strconv.ParseUint(key[i:i+2], 16, 8)
		if err != nil {
			return nil, errors.New(err).Wrap("invalid key -- must be communicated in hexadecimal format")
		}

		ret = append(ret, byte(i))
	}

	switch len(ret) {
	case 8, 16, 32:
	default:
		return nil, errors.New("invalid key size -- must equate to 8, 16, or 32 bytes to be a valid AES key")
	}

	return ret, nil
}
