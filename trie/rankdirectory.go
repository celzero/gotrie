// Copyright (c) 2025 RethinkDNS and its authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package trie

import (
	"fmt"
	"math"
)

type RankDirectory struct {
	Directory   *BStr
	Data        *BStr
	l1Size      int
	l2Size      int
	l1Bits      int
	l2Bits      int
	sectionBits int
	numBits     int
}

func NewRankDir(rd, td *BStr, numBits int, l1Size int, l2Size int) *RankDirectory {
	l1Bits := int(math.Ceil(math.Log2(float64(numBits))))
	l2Bits := int(math.Ceil(math.Log2(float64(l1Size))))
	rdir := &RankDirectory{
		Directory:   rd,
		Data:        td,
		l1Size:      l1Size,
		l2Size:      l2Size,
		l1Bits:      l1Bits,
		l2Bits:      l2Bits,
		sectionBits: (l1Size/l2Size-1)*l2Bits + l1Bits,
		numBits:     numBits,
	}
	if Debug {
		rdir.display()
	}
	return rdir
}

func (RD *RankDirectory) display() {
	fmt.Println("sz(rd): ", RD.Directory.Size())
	fmt.Println("sz(td): ", RD.Data.Size())
	fmt.Println("numBits: ", RD.numBits)
	fmt.Println("L1size: ", RD.l1Size)
	fmt.Println("l1bits: ", RD.l1Bits)
	fmt.Println("l2bits: ", RD.l2Bits)
}

func (rdir *RankDirectory) rank(which, x int) int {
	var temp uint32
	rank := -1
	sectionPos := 0
	if x >= rdir.l2Size {
		sectionPos = (x / rdir.l2Size) * rdir.l1Bits
		temp = rdir.Directory.get(sectionPos-rdir.l1Bits, rdir.l1Bits, false)
		rank = int(temp)
		x = x % rdir.l2Size
	}
	var ans = 0
	if x > 0 {
		ans = rdir.Data.pos0(rank+1, x)
	} else {
		ans = rank
	}
	if Debug {
		fmt.Printf("ans: %d %d:r, x: %d %d:s %d:l1 %t:ifcheck\n", ans, temp, x, sectionPos, rdir.l1Bits, x >= rdir.l2Size)
	}
	return ans
}
