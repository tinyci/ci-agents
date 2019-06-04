package types

import (
	"crypto/rand"
	"fmt"
	"math/big"

	check "github.com/erikh/check"
)

// copied over from testutil because testutil depends on types
func randString(len int) string {
	ret := ""
	charlen := 'z' - 'a'

	for i := 0; i < len; i++ {
		rnd, _ := rand.Int(rand.Reader, big.NewInt(int64(charlen)))
		ret += fmt.Sprintf("%c", 'a'+rune(rnd.Int64()))
	}

	return ret
}

func (ts *typesSuite) TestDecodeKey(c *check.C) {
	fmt.Print("gathering entropy...")
	bytSize := 1024 * 1024
	chunkSize := 64
	bigbytes := make([]byte, bytSize)

	n, err := rand.Reader.Read(bigbytes)
	c.Assert(err, check.IsNil)
	c.Assert(n, check.Equals, bytSize)

	fmt.Println("done.")

	for i := 0; i < bytSize; i += chunkSize {
		region := bigbytes[i : i+chunkSize]
		enc := EncodeKey(region)
		c.Assert(len(enc), check.Equals, chunkSize*2, check.Commentf("Key: %s", enc))
		c.Assert(DecodeKey(enc), check.DeepEquals, region, check.Commentf("Key: %s", enc))
	}
}

func (ts *typesSuite) TestTokenEncryption(c *check.C) {
	iterations := 10000

	for i := 0; i < iterations; i++ {
		keyBytes := make([]byte, 32)
		n, err := rand.Reader.Read(keyBytes)
		c.Assert(err, check.IsNil)
		c.Assert(n, check.Equals, 32)

		l, err := rand.Int(rand.Reader, big.NewInt(10))
		c.Assert(err, check.IsNil)
		scopes := make([]string, l.Int64())

		for x := 0; int64(x) < l.Int64(); x++ {
			scopes[x] = randString(10)
		}

		tok := OAuthToken{
			Token:    randString(128),
			Username: randString(8),
			Scopes:   scopes,
		}

		text, err := EncryptToken(keyBytes, &tok)
		c.Assert(err, check.IsNil)
		c.Assert(len(text), check.Not(check.Equals), 0)
		tok2, err := DecryptToken(keyBytes, text)
		c.Assert(err, check.IsNil)
		c.Assert(tok, check.DeepEquals, *tok2)
	}
}
