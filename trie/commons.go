package trie

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"unsafe"
)

func bit0p(n int, p int) map[string]int {
	var ret = make(map[string]int)
	if p == 0 {
		ret["index"] = 0
		ret["scanned"] = 0
		return ret
	}

	if n == 0 && p == 1 {
		ret["index"] = 1
		ret["scanned"] = 1
		return ret
	}
	var c = 0
	var i = 0
	//var m = n
	for c = 0; n > 0 && p > c; n = n >> 1 {
		// increment c when nth lsb (bit) is 0
		if n < (n ^ 0x1) {
			c = c + 1
		} else {
			c = c + 0
		}
		//c = c + (n < (n ^ 0x1)) ? 1 : 0;
		i += 1
	}
	ret["scanned"] = i
	if p == c {
		ret["index"] = i
	} else {
		ret["index"] = 0
	}
	//console.log("      ", String.fromCharCode(m).charCodeAt(0).toString(2), m, i, p, c);
	return ret
}

func Bit0(n int, p int, pad int) int {
	var r = bit0p(n, p)
	if r["scanned"] <= 0 {
		return r["scanned"]
	} // r.index
	if r["index"] > 0 {
		return r["scanned"]
	} // r.index
	if pad > r["scanned"] {
		return (r["scanned"] + 1)
	} else {
		return 0
	}
}

func DEC16(str string, index int) int {
	res := int([]rune(str)[index])
	return res
}
func CHR16(charcode int) string {
	return string(charcode)
}

func FlagSubstring(str string, si int, ei int) string {
	res := []rune(str)
	resstr := ""
	if ei == 0 {
		resstr = string(res[si:])
	} else {
		resstr = string(res[si:ei])
	}
	return resstr
}

func Flag_to_uint(str string) []uint32 {
	runedata := []rune(str)
	resp := make([]uint32, len(runedata))
	for key, value := range runedata {
		resp[key] = uint32(value)
	}
	return resp
}

func TxtEncode(str string) ([]uint8, error) {
	strbytes := []byte(str)
	r := bytes.NewReader(strbytes)
	tmp_u8 := make([]uint8, len(strbytes))
	err := binary.Read(r, binary.LittleEndian, &tmp_u8)
	if err != nil {
		fmt.Println("Error At byte to uint16 conversion : common.go -> TxtEncode()")
		return nil, err
	}

	for i, j := 0, len(tmp_u8)-1; i < j; i, j = i+1, j-1 {
		tmp_u8[i], tmp_u8[j] = tmp_u8[j], tmp_u8[i]
	}

	return tmp_u8, nil
}

func Find_Lista_Listb(list1 []string, list2 []string) (bool, []string) {
	retlist := []string{}
	if len(list1) <= 0 || len(list2) <= 0 {
		// fmt.Printf("return false and empty retlist\n")
		return false, retlist
	}
	//fmt.Printf("----list len %d %d\n",len(list1), len(list2))

	found := false
	for _, value := range list1 {
		// FIXME: FT.usr_bl can contain empty strings and that
		// shoudln't be matched at all. But then again commons
		// isn't the place to account for such edge-cases.
		if len(value) <= 0 {
			continue
		}
		for _, ivalue := range list2 {
			if value == ivalue {
				found = true
				retlist = append(retlist, value)
			}
		}
	}

	return found, retlist
}

// from: stackoverflow.com/a/25286918
func MD5Hex(b []byte) string {
	hash := md5.Sum(b)
	return hex.EncodeToString(hash[:])
}

// stackoverflow.com/a/74643597
func castToBytes[T any](t *[]T) []byte {
	s := *t
	l := len(s)
	if l == 0 {
		return nil
	}

	sz := int(unsafe.Sizeof(s[0])) * l
	return unsafe.Slice((*byte)(unsafe.Pointer(&s[0])), sz)
}
