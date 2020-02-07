// +build ignore

package main

// reference: https://blog.csdn.net/qq_36291381/article/details/80156315

import (
	"fmt"
	"math/bits"
	"os"
)

// irreducible polynomial as m(x)=x^8+x^7+x^6+x^5+x^4+x^2+1
//const poly = 0x1f5
const poly = 1<<8 | 1<<7 | 1<<6 | 1<<5 | 1<<4 | 1<<2 | 1<<0

func change(x uint8) uint8 {
	var out uint8

	A1 := uint8(0xa7)
	for i := 0; i < 8; i++ {
		y := x & A1

		var z uint8
		for j := 0; j < 8; j++ {
			z ^= y & 1
			y >>= 1
		}

		out = out | (z << i)
		A1 = bits.RotateLeft8(A1, 1)
	}

	return out ^ 0xd3
}

// div calculates (q,r), where a = q*b+r
func div(a, b uint16) (uint16, uint16) {
	var q, r uint16

	bLen := (16 - bits.LeadingZeros16(b))
	for {
		d := (16 - bits.LeadingZeros16(a)) - bLen

		if d < 0 || a == 0 {
			r = a
			break
		}

		// q' = 2^d, a = a-b*q'
		a ^= b << d
		q |= 1 << d
	}

	return q, r
}

func generateAndWriteCKTo(file *os.File) {
	ck := generateCK()
	mustWriteCK(ck, file)
}

func generateAndWriteSboxTo(file *os.File) {
	var sbox, sboxInv [256]byte
	for i := 0; i <= 255; i++ {
		ii := uint8(i)
		j := change(uint8(inverse(change(ii), poly)))
		sbox[i], sboxInv[j] = j, ii
	}

	mustWriteSbox(sbox, "sbox", file)
	mustWriteString("\n", file)
	mustWriteSbox(sboxInv, "sboxInv", file)
}

func generateCK() [32]uint32 {
	var ck [32]uint32

	for i := range ck {
		ck[i] = ((uint32(4*i+0) * 7 % 256) << 24) |
			((uint32(4*i+1) * 7 % 256) << 16) |
			((uint32(4*i+2) * 7 % 256) << 8) |
			(uint32(4*i+3) * 7 % 256)
	}

	return ck
}

// inverse calculates the inverse of b mod a
func inverse(b uint8, a uint16) uint16 {
	// note: v is of no use
	var (
		w1 uint16 = 0
		w0 uint16 = 1
	)

	b16 := uint16(b)
	for b16 > 0 {
		q, r := div(a, b16)

		a, b16, w1, w0 = b16, r, w0, w1^mul(q, w0)
	}

	return w1
}

func mul(a, b uint16) uint16 {
	var out uint16

	for i := 0; (i < 8) && ((b >> i) != 0); i++ {
		if (b>>i)&1 == 1 {
			out ^= a << i
		}
	}

	return out
}

func mustWriteString(str string, file *os.File) {
	if _, err := file.WriteString(str); err != nil {
		panic(err)
	}
}

func mustWriteCK(ck [32]uint32, file *os.File) {
	mustWriteString("\nvar ck = [32]uint32{\n", file)
	for i := 0; i < 8; i++ {
		mustWriteString("\t", file)
		for j := 0; j < 4; j++ {
			mustWriteString(fmt.Sprintf("0x%08x, ", ck[4*i+j]), file)
		}
		mustWriteString("\n", file)
	}
	mustWriteString("}\n", file)
}

func mustWriteSbox(arr [256]byte, varName string, file *os.File) {
	mustWriteString("var "+varName+" = [256]byte{\n", file)
	for i := 0; i < 16; i++ {
		mustWriteString("\t", file)
		for j := 0; j < 16; j++ {
			mustWriteString(fmt.Sprintf("0x%02x, ", arr[i*16+j]), file)
		}
		mustWriteString("\n", file)
	}
	mustWriteString("}\n", file)
}

func main() {

	//fmt.Println(sbox)

	file, err := os.OpenFile("parameter.go", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	mustWriteString("// AUTO GENERATED. DO NOT EDIT. see generate_parameter.go\n\n", file)
	mustWriteString("package sm4\n\n", file)

	generateAndWriteSboxTo(file)

	generateAndWriteCKTo(file)
}
