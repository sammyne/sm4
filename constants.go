package sm4

const (
	BlockSize = 16
	KeySize   = 16
)

const Round = 32

var fk = [KeySize / 4]uint32{0xa3b1bac6, 0x56aa3350, 0x677d9197, 0xb27022dc}
