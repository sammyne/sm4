package sm4

import "math/bits"

func rotBlockWord(w uint32) uint32 {
	return w ^ bits.RotateLeft32(w, 2) ^
		bits.RotateLeft32(w, 10) ^
		bits.RotateLeft32(w, 18) ^
		bits.RotateLeft32(w, 24)
}

func rotKeyWord(w uint32) uint32 {
	return w ^ bits.RotateLeft32(w, 13) ^ bits.RotateLeft32(w, 23)
}

// Apply sbox to each byte in w.
func subw(w uint32) uint32 {
	return uint32(sbox[w>>24])<<24 |
		uint32(sbox[w>>16&0xff])<<16 |
		uint32(sbox[w>>8&0xff])<<8 |
		uint32(sbox[w&0xff])
}
