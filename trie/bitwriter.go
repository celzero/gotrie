package trie

//import "fmt"

type Bitwriter struct {
	bits   []byte
	bytes  []byte
	bits16 []byte
	Top    int32
}

/*
func (wr *Bitwriter) Write16(data int16, numBits int16) {
    // todo: throw error?
    if (numBits > 16) {
        fmt.Println("write16 can only writes lsb16 bits, out of range: ");
        return;
    }
    n := data;
    brim := int16(16 - (wr.Top % 16));
    cur := int16((wr.Top / 16) | 0);
    e := wr.bits16[cur] | 0;
    remainingBits := int16(0);
    // clear msb
    b := n & BitString.MaskTop[16][16 - numBits];

    // shift to bit pos to be right at brim-th bit
    if (brim >= numBits) {
        b = b << (brim - numBits);
    } else {
        // shave right most bits if there are too many bits than
        // what the current element at the brim can accomodate
        remainingBits = numBits - brim;
        b = b >> remainingBits;
    }
    // overlay b on current element, e.
    b = e | b;
    wr.bits16[cur] = b;

    // account for the left-over bits shaved off by brim
    if (remainingBits > 0) {
        b = n & BitString.MaskTop[16][16 - remainingBits];
        b = b << (16 - remainingBits);
        wr.bits16[cur + 1] = b;
    }

    // update top to reflect the bits included
    wr.Top = wr.Top + int32(numBits);
}
*/
