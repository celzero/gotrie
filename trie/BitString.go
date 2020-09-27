package trie

import "math"
import "fmt"

type BS struct {
    MaskTop         map[int][]uint16
    MaskBottom      map[int][]uint16
    BitsSetTable256 [256]int
    bytes           []uint16
    length          int
    useBuffer       bool
}

func (BitString *BS) Init(str []uint16) error {
    BitString.MaskTop = make(map[int][]uint16)
    BitString.MaskTop[16] = []uint16{}
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0xffff)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x7fff)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x3fff)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x1fff)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x0fff)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x07ff)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x03ff)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x01ff)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x00ff)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x007f)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x003f)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x001f)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x000f)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x0007)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x0003)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x0001)
    BitString.MaskTop[16] = append(BitString.MaskTop[16], 0x0000)

    BitString.MaskBottom = make(map[int][]uint16)
    BitString.MaskBottom[16] = []uint16{}
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xffff)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xfffe)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xfffc)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xfff8)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xfff0)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xffe0)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xffc0)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xff80)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xff00)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xfe00)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xfc00)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xf800)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xf000)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xe000)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0xc000)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0x8000)
    BitString.MaskBottom[16] = append(BitString.MaskBottom[16], 0x0000)

    BitString.BitsSetTable256[0] = 0
    for i := 0; i < 256; i++ {
        BitString.BitsSetTable256[i] = (i & 1) + BitString.BitsSetTable256[int(math.Floor(float64(i/2)))]
    }

    BitString.bytes = str
    BitString.length = len(BitString.bytes) * W
    BitString.useBuffer = true

    //fmt.Printf("Length : %d\n",BitString.length)
    //fmt.Printf("str Length : %d\n",len(*str))
    return nil
}

func (BitsString BS) display() {
    fmt.Println(len(BitsString.bytes))
}
func (BitString BS) countSetBits(n int) int {
    return (BitString.BitsSetTable256[n&0xff] +
        BitString.BitsSetTable256[(n>>8)&0xff] +
        BitString.BitsSetTable256[(n>>16)&0xff] +
        BitString.BitsSetTable256[n>>24])
}

func (BitString BS) getData() []uint16 {
    return BitString.bytes
}

/*
func (BitString BS) encode(n int64)[]uint16 {
    var e = []uint16{};
    for i:=int64(0); i < int64(BitString.length); i = i + n {
        e = append(e,BitString.get(i,int64(math.Min(float64(BitString.length),float64(n))),false))
        //e.push(this.get(i, Math.min(this.length, n)));
    }
    return e;
}*/

func (BitString BS) get(p int, n int, debug bool) uint32 {

    if (p%W)+n <= W {
        return uint32((BitString.bytes[p/W|0] & BitString.MaskTop[int(W)][p%W]) >> (W - (p % W) - n))
        // case 2: bits lie incompletely in the given byte
    } else {
        var result uint32
        var l int
        result = uint32(BitString.bytes[p/W|0] & BitString.MaskTop[int(W)][p%W])

        tmp_count := 0 //santhosh added
        disp1 := BitString.bytes[p/W|0]
        disp2 := BitString.MaskTop[int(W)][p%W]
        var res1 = result

        l = W - p%W
        p += l
        n -= l

        for n >= W {
            tmp_count += 1
            result = (result << W) | uint32(BitString.bytes[p/W|0])
            p += W
            n -= W
        }

        var res2 = result
        if n > 0 {
            result = (result << n) | uint32(BitString.bytes[p/W|0]>>(W-n))
        }

        if debug {
            fmt.Printf("disp1: %d disp2: %d loopcount: %d res1:%d res2:%d r:%d\n", disp1, disp2, tmp_count, res1, res2, result)
        }
        return result
    }
}

func (BitString BS) count(p int, n int) int {
    count := 0
    for n >= 16 {
        count = count + (BitString.BitsSetTable256[BitString.get(p, 16, false)])
        p += 16
        n -= 16
    }

    return count + BitString.BitsSetTable256[BitString.get(p, n, false)]
}

func (BitString BS) pos0(i int, n int) int {
    if n < 0 {
        return 0
    }
    step := 16
    index := i
    for n > 0 {
        d := BitString.get(i, step, false)
        bits0 := step - BitString.countSetBits(int(d))
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
