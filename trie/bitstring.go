package trie

import (
	"fmt"
	"math"
)

type BStr struct {
	bytes  []uint16
	length int
}

func NewBStr(str *[]uint16) *BStr {
	if str == nil {
		return nil
	}

	bs := new(BStr)

	bs.bytes = *str
	bs.length = len(bs.bytes) * W

	//fmt.Printf("Length : %d\n",BitString.length)
	//fmt.Printf("str Length : %d\n",len(*str))
	return bs
}

func (bs *BStr) get(p int, n int, debug bool) uint32 {

	if (p%W)+n <= W {
		return uint32((bs.bytes[p/W] & MaskTop[int(W)][p%W]) >> (W - (p % W) - n))
		// case 2: bits lie incompletely in the given byte
	} else {
		var result uint32
		var l int
		result = uint32(bs.bytes[p/W] & MaskTop[int(W)][p%W])

		tmp_count := 0 //santhosh added
		disp1 := bs.bytes[p/W]
		disp2 := MaskTop[int(W)][p%W]
		var res1 = result

		l = W - p%W
		p += l
		n -= l

		for n >= W {
			tmp_count += 1
			result = (result << W) | uint32(bs.bytes[p/W])
			p += W
			n -= W
		}

		var res2 = result
		if n > 0 {
			result = (result << n) | uint32(bs.bytes[p/W]>>(W-n))
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
		i := bs.get(p, 16, false)
		count = count + (BitsSetTable256[i])
		p += 16
		n -= 16
	}
	return count + BitsSetTable256[bs.get(p, n, false)]
}

func (bs *BStr) pos0(i int, n int) int {
	if n < 0 {
		return 0
	}
	step := 16
	index := i
	for n > 0 {
		d := bs.get(i, step, false)
		bits0 := step - CountSetBits(int(d))
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

var MaskTop = make(map[int][]uint16)
var MaskBottom = make(map[int][]uint16)
var BitsSetTable256 [256]int

func CountSetBits(n int) int {
	return (BitsSetTable256[n&0xff] +
		BitsSetTable256[(n>>8)&0xff] +
		BitsSetTable256[(n>>16)&0xff] +
		BitsSetTable256[n>>24])
}

func init() {
	MaskTop[16] = []uint16{}
	MaskTop[16] = append(MaskTop[16], 0xffff)
	MaskTop[16] = append(MaskTop[16], 0x7fff)
	MaskTop[16] = append(MaskTop[16], 0x3fff)
	MaskTop[16] = append(MaskTop[16], 0x1fff)
	MaskTop[16] = append(MaskTop[16], 0x0fff)
	MaskTop[16] = append(MaskTop[16], 0x07ff)
	MaskTop[16] = append(MaskTop[16], 0x03ff)
	MaskTop[16] = append(MaskTop[16], 0x01ff)
	MaskTop[16] = append(MaskTop[16], 0x00ff)
	MaskTop[16] = append(MaskTop[16], 0x007f)
	MaskTop[16] = append(MaskTop[16], 0x003f)
	MaskTop[16] = append(MaskTop[16], 0x001f)
	MaskTop[16] = append(MaskTop[16], 0x000f)
	MaskTop[16] = append(MaskTop[16], 0x0007)
	MaskTop[16] = append(MaskTop[16], 0x0003)
	MaskTop[16] = append(MaskTop[16], 0x0001)
	MaskTop[16] = append(MaskTop[16], 0x0000)

	MaskBottom[16] = []uint16{}
	MaskBottom[16] = append(MaskBottom[16], 0xffff)
	MaskBottom[16] = append(MaskBottom[16], 0xfffe)
	MaskBottom[16] = append(MaskBottom[16], 0xfffc)
	MaskBottom[16] = append(MaskBottom[16], 0xfff8)
	MaskBottom[16] = append(MaskBottom[16], 0xfff0)
	MaskBottom[16] = append(MaskBottom[16], 0xffe0)
	MaskBottom[16] = append(MaskBottom[16], 0xffc0)
	MaskBottom[16] = append(MaskBottom[16], 0xff80)
	MaskBottom[16] = append(MaskBottom[16], 0xff00)
	MaskBottom[16] = append(MaskBottom[16], 0xfe00)
	MaskBottom[16] = append(MaskBottom[16], 0xfc00)
	MaskBottom[16] = append(MaskBottom[16], 0xf800)
	MaskBottom[16] = append(MaskBottom[16], 0xf000)
	MaskBottom[16] = append(MaskBottom[16], 0xe000)
	MaskBottom[16] = append(MaskBottom[16], 0xc000)
	MaskBottom[16] = append(MaskBottom[16], 0x8000)
	MaskBottom[16] = append(MaskBottom[16], 0x0000)

	BitsSetTable256[0] = 0
	for i := 0; i < 256; i++ {
		j := int(math.Floor(float64(i / 2))) // j = i >> 1
		BitsSetTable256[i] = (i & 1) + BitsSetTable256[j]
	}
}
