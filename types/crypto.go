package types

// I am no crypto expert, I tried to do the right thing though. Comments and
// patches to this section are welcome from people willing to cite
// authoritative sources for their changes.

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"errors"

	"github.com/tinyci/ci-agents/utils"
)

// EncryptToken encrypts an oauth token with gcm+aes.
func EncryptToken(key []byte, tok *OAuthToken) ([]byte, error) {
	// I know the storage management of this is terrible, but I don't have a better
	// solution (or know someone who has one).
	//
	// Offline key usage is the problem here. We could trivially associate the
	// mac with the (encrypted by other means) session cookies in the UI, but
	// that won't resolve the case where CI jobs need the token.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, utils.WrapError(err, "while setting up encrypter")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, utils.WrapError(err, "while setting up encryption cipher suite")
	}

	nonce := make([]byte, gcm.NonceSize())
	c, err := rand.Reader.Read(nonce)
	if err != nil {
		return nil, utils.WrapError(err, "reading entropy")
	}

	if c < gcm.NonceSize() {
		return nil, errors.New("could not read enough entropy to encrypt token")
	}

	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(tok); err != nil {
		return nil, err
	}

	outbuf := gcm.Seal(nil, nonce, buf.Bytes(), nil)

	// the price of laziness
	return []byte(base64.StdEncoding.EncodeToString(append(nonce, outbuf...))), nil
}

// DecryptToken decrypts a token message that was encrypted by EncryptToken.
func DecryptToken(key, tokenBytes []byte) (*OAuthToken, error) {
	if len(tokenBytes) == 0 {
		return &OAuthToken{}, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(string(tokenBytes))
	if err != nil {
		return nil, err
	}

	if len(decoded) == 0 {
		return &OAuthToken{}, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, utils.WrapError(err, "while setting up decrypter")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, utils.WrapError(err, "while setting up decryption cipher suite")
	}

	if len(decoded) < gcm.NonceSize()+1 {
		return nil, errors.New("invalid token bytes")
	}

	tb, err := gcm.Open(nil, decoded[:gcm.NonceSize()], decoded[gcm.NonceSize():], nil)
	if err != nil {
		return nil, utils.WrapError(err, "decrypting token")
	}

	tok := OAuthToken{}
	if err := json.Unmarshal(tb, &tok); err != nil {
		return nil, utils.WrapError(err, "While decrypting and unmarshalling token")
	}

	return &tok, nil
}

// DecodeKey takes a hexadecimal key and turns it into a byte array where the
// bytes represent the by-two hexadecimal numerical values present in the
// original key.
func DecodeKey(key string) []byte {
	bytKey := []byte{}

	i := 0
	for i < len(key) {
		var byt byte
		fmt.Fscanf(bytes.NewBuffer([]byte(key[i:i+2])), "%02x", &byt)
		bytKey = append(bytKey, byt)
		i += 2
	}

	return bytKey
}

// EncodeKey takes a byte stream and yields a hexadecimal set of values
// associated with it. See `xxd(1)`.
func EncodeKey(key []byte) string {
	var res string

	for _, byt := range key {
		res += fmt.Sprintf("%02x", byt)
	}

	return res
}
