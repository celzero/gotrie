package trie

import (
	"fmt"
	"math"
)

type BStr struct {
	MaskTop         map[int][]uint16
	MaskBottom      map[int][]uint16
	BitsSetTable256 [256]int
	bytes           []uint16
	length          int
}

func NewBStr(str *[]uint16) *BStr {
	if str == nil {
		return nil
	}

	bs := new(BStr)
	bs.MaskTop = make(map[int][]uint16)
	bs.MaskTop[16] = []uint16{}
	bs.MaskTop[16] = append(bs.MaskTop[16], 0xffff)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x7fff)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x3fff)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x1fff)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x0fff)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x07ff)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x03ff)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x01ff)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x00ff)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x007f)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x003f)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x001f)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x000f)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x0007)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x0003)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x0001)
	bs.MaskTop[16] = append(bs.MaskTop[16], 0x0000)

	bs.MaskBottom = make(map[int][]uint16)
	bs.MaskBottom[16] = []uint16{}
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xffff)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xfffe)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xfffc)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xfff8)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xfff0)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xffe0)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xffc0)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xff80)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xff00)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xfe00)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xfc00)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xf800)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xf000)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xe000)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0xc000)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0x8000)
	bs.MaskBottom[16] = append(bs.MaskBottom[16], 0x0000)

	bs.BitsSetTable256[0] = 0
	for i := 0; i < 256; i++ {
		bs.BitsSetTable256[i] = (i & 1) + bs.BitsSetTable256[int(math.Floor(float64(i/2)))]
	}

	bs.bytes = *str
	bs.length = len(bs.bytes) * W

	//fmt.Printf("Length : %d\n",BitString.length)
	//fmt.Printf("str Length : %d\n",len(*str))
	return bs
}

func (bs *BStr) display() {
	fmt.Println(len(bs.bytes))
}
func (BitString *BStr) countSetBits(n int) int {
	return (BitString.BitsSetTable256[n&0xff] +
		BitString.BitsSetTable256[(n>>8)&0xff] +
		BitString.BitsSetTable256[(n>>16)&0xff] +
		BitString.BitsSetTable256[n>>24])
}

func (bs *BStr) getData() []uint16 {
	return bs.bytes
}

func (bs *BStr) get(p int, n int, debug bool) uint32 {

	if (p%W)+n <= W {
		return uint32((bs.bytes[p/W|0] & bs.MaskTop[int(W)][p%W]) >> (W - (p % W) - n))
		// case 2: bits lie incompletely in the given byte
	} else {
		var result uint32
		var l int
		result = uint32(bs.bytes[p/W|0] & bs.MaskTop[int(W)][p%W])

		tmp_count := 0 //santhosh added
		disp1 := bs.bytes[p/W|0]
		disp2 := bs.MaskTop[int(W)][p%W]
		var res1 = result

		l = W - p%W
		p += l
		n -= l

		for n >= W {
			tmp_count += 1
			result = (result << W) | uint32(bs.bytes[p/W|0])
			p += W
			n -= W
		}

		var res2 = result
		if n > 0 {
			result = (result << n) | uint32(bs.bytes[p/W|0]>>(W-n))
		}

		if debug {
			fmt.Printf("disp1: %d disp2: %d loopcount: %d res1:%d res2:%d r:%d\n", disp1, disp2, tmp_count, res1, res2, result)
		}
		return result
	}
}

func (bs *BStr) count(p int, n int) int {
	count := 0
	for n >= 16 {
		count = count + (bs.BitsSetTable256[bs.get(p, 16, false)])
		p += 16
		n -= 16
	}

	return count + bs.BitsSetTable256[bs.get(p, n, false)]
}

func (bs *BStr) pos0(i int, n int) int {
	if n < 0 {
		return 0
	}
	step := 16
	index := i
	for n > 0 {
		d := bs.get(i, step, false)
		bits0 := step - bs.countSetBits(int(d))
		diff := 0
		if n-bits0 < 0 {
			step = int(math.Max(float64(n), float64(step/2|0)))
			continue
		}
		n -= bits0
		i += step
		if n == 0 {
			diff = Bit0(int(d), 1, step)
		} else {
			diff = 1
		}
		index = i - diff // 1;
	}

	return index
}
