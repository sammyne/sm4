package sm4

import "testing"

func TestExpandKey(t *testing.T) {
	testVector := []struct {
		key    []byte
		expect [roundKeysLen]uint32
	}{
		{
			[]byte{
				0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
				0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10,
			},
			[roundKeysLen]uint32{
				0xf12186f9,
				0x41662b61,
				0x5a6ab19a,
				0x7ba92077,
				0x367360f4,
				0x776a0c61,
				0xb6bb89b3,
				0x24763151,
				0xa520307c,
				0xb7584dbd,
				0xc30753ed,
				0x7ee55b57,
				0x6988608c,
				0x30d895b7,
				0x44ba14af,
				0x104495a1,
				0xd120b428,
				0x73b55fa3,
				0xcc874966,
				0x92244439,
				0xe89e641f,
				0x98ca015a,
				0xc7159060,
				0x99e1fd2e,
				0xb79bd80c,
				0x1d2115b0,
				0x0e228aeb,
				0xf1780c81,
				0x428d3654,
				0x62293496,
				0x01cf72e5,
				0x9124a012,
			},
		},
	}

	for i, c := range testVector {
		if got := expandKey(c.key); got != c.expect {
			t.Fatalf("#%d: expect %v, got %v", i, c.expect, got)
		}
	}
}