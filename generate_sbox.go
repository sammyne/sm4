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

func mustWriteTo(arr [256]byte, varName string, file *os.File) {
	mustWriteString("var "+varName+" = [256]byte{", file)
	for i, v := range arr {
		if i%16 == 0 {
			//fmt.Printf("%02x 0x%02x\n",i, v)
			mustWriteString("\n\t", file)
		}
		mustWriteString(fmt.Sprintf("0x%02x, ", v), file)
	}
	mustWriteString("}\n", file)
}

func main() {
	/*
		for i := 0; i <= 0x0f; i++ {
			fmt.Printf("\t%x", i)
		}
		fmt.Println()

		for i := uint8(0); i <= 0x0f; i++ {
			fmt.Printf("%x", i)
			for j := uint8(0); j <= 0x0f; j++ {
				fmt.Printf("\t%02x", change(uint8(inverse(change((i<<4)|j), poly))))
			}
			fmt.Println()
		}
	*/

	var sbox, sboxInv [256]byte
	for i := 0; i <= 255; i++ {
		ii := uint8(i)
		j := change(uint8(inverse(change(ii), poly)))
		sbox[i], sboxInv[j] = j, ii
	}

	//fmt.Println(sbox)

	file, err := os.OpenFile("sbox.go", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	mustWriteString("// AUTO GENERATED. DO NOT EDIT. see generate_sbox.go\n\n", file)
	mustWriteString("package sm4\n\n", file)

	mustWriteTo(sbox, "sbox", file)
	mustWriteString("\n",file)
	mustWriteTo(sboxInv, "sboxInv", file)
}
