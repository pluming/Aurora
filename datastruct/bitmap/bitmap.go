package bitmap

type bitmap []byte

func NewBitmap() *bitmap {
	b := bitmap(make([]byte, 0))
	return &b
}

func (b *bitmap) SetBit(offset int64, val byte) {
	byteIndex := offset >> 3
	bitOffset := offset % 8
	//check Or grow bitmap
	mask := byte(1 << bitOffset)
	b.grow(offset + 1)
	if val > 0 {
		// set bit
		(*b)[byteIndex] |= mask
	} else {
		//clear bit
		(*b)[byteIndex] = ^mask
	}
}

func (b *bitmap) BitSize() int {
	return len(*b) * 8
}

func (b *bitmap) ToBytes() []byte {
	return *b
}

func (b *bitmap) GetBit(offset int64) byte {
	bytesIndex := offset >> 3
	bitIndex := offset % 8
	if bytesIndex > int64(len(*b)) {
		return 0
	}
	return ((*b)[bytesIndex] >> bitIndex) & 0x01
}

func (b *bitmap) ForEachBit(s, e int64, cb Callback) {
	offset := s
	bytesIndex := offset >> 3
	bitIndex := offset % 8

	for bytesIndex < int64(len(*b)) {
		vb := (*b)[bytesIndex]
		for bitIndex < 8 {
			bi := byte(vb>>bitIndex) & 0x01
			if !cb(offset, bi) {
				return
			}
			offset++
			bitIndex++
		}
		if e > 0 && offset == e {
			return
		}
		bytesIndex++
		bitIndex = 0
	}
}

func (b *bitmap) ForEachByte(s, e int, cb Callback) {
	if e == 0 {
		e = len(*b)
	} else if e > len(*b) {
		e = len(*b)
	}
	for i := s; i < e; i++ {
		if !cb(int64(i), (*b)[i]) {
			return
		}
	}
}

func (b *bitmap) grow(bitSize int64) {
	var byteSize int64
	if bitSize%8 == 0 {
		byteSize = bitSize / 8
	}
	byteSize = bitSize/8 + 1
	gap := byteSize - int64(len(*b))
	if gap <= 0 {
		return
	}
	*b = append(*b, make([]byte, gap)...)
}

func FromBytes(bytes []byte) *bitmap {
	bm := bitmap(bytes)
	return &bm
}

type Callback func(offset int64, val byte) bool
