package sm4

import (
	"crypto/cipher"
	"encoding/binary"
)

//go:generate go run generate_sbox.go

type sm4Cipher struct {
	expandedKey [round]uint32
}

// BlockSize returns the cipher's block size.
func (c *sm4Cipher) BlockSize() int {
	return BlockSize
}

// Decrypt decrypts the first block in src into dst.
// Dst and src must overlap entirely or not at all.
func (c *sm4Cipher) Decrypt(dst []byte, src []byte) {
	_ = src[15] // early bounds check

	_ = src[15] // early bounds check

	x0 := binary.BigEndian.Uint32(src[0:4])
	x1 := binary.BigEndian.Uint32(src[4:8])
	x2 := binary.BigEndian.Uint32(src[8:12])
	x3 := binary.BigEndian.Uint32(src[12:16])

	for r := round - 1; r >= 0; r-- {
		x0, x1, x2, x3 = x1, x2, x3, x0^rotBlockWord(subw(x1^x2^x3^c.expandedKey[r]))
	}

	_ = dst[15] // early bounds check
	binary.BigEndian.PutUint32(dst[0:4], x3)
	binary.BigEndian.PutUint32(dst[4:8], x2)
	binary.BigEndian.PutUint32(dst[8:12], x1)
	binary.BigEndian.PutUint32(dst[12:16], x0)

	_ = dst[15] // early bounds check
}

// Encrypt encrypts the first block in src into dst.
// Dst and src must overlap entirely or not at all.
func (c *sm4Cipher) Encrypt(dst []byte, src []byte) {
	_ = src[15] // early bounds check

	x0 := binary.BigEndian.Uint32(src[0:4])
	x1 := binary.BigEndian.Uint32(src[4:8])
	x2 := binary.BigEndian.Uint32(src[8:12])
	x3 := binary.BigEndian.Uint32(src[12:16])

	for r := 0; r < round; r++ {
		x0, x1, x2, x3 = x1, x2, x3, x0^rotBlockWord(subw(x1^x2^x3^c.expandedKey[r]))
	}

	_ = dst[15] // early bounds check
	binary.BigEndian.PutUint32(dst[0:4], x3)
	binary.BigEndian.PutUint32(dst[4:8], x2)
	binary.BigEndian.PutUint32(dst[8:12], x1)
	binary.BigEndian.PutUint32(dst[12:16], x0)
}

func NewCipher(key []byte) (cipher.Block, error) {
	if keyLen := len(key); keyLen != KeySize {
		return nil, KeySizeError(keyLen)
	}

	c := &sm4Cipher{
		expandedKey: expandKey(key),
	}

	return c, nil
}

func expandKey(key []byte) [round]uint32 {
	// Encryption key setup.
	k0 := binary.BigEndian.Uint32(key[0:4]) ^ fk[0]
	k1 := binary.BigEndian.Uint32(key[4:8]) ^ fk[1]
	k2 := binary.BigEndian.Uint32(key[8:12]) ^ fk[2]
	k3 := binary.BigEndian.Uint32(key[12:16]) ^ fk[3]

	var expandedKey [round]uint32
	for i := range expandedKey {
		k0, k1, k2, k3 = k1, k2, k3, k0^rotKeyWord(subw(k1^k2^k3^ck[i]))
		expandedKey[i] = k3
	}

	return expandedKey
}
