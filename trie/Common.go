package trie

import "bytes"
import "encoding/binary"
import "fmt"

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

	//fmt.Println(tmp_u8)
	for i, j := 0, len(tmp_u8)-1; i < j; i, j = i+1, j-1 {
		tmp_u8[i], tmp_u8[j] = tmp_u8[j], tmp_u8[i]
	}

	//fmt.Println(tmp_u8)
	return tmp_u8, nil
}

func Find_Lista_Listb(list1 []string, list2 []string) (bool, []string) {
	retlist := []string{}
	found := false
	for _, value := range list1 {
		for _, ivalue := range list2 {
			if value == ivalue {
				found = true
				retlist = append(retlist, value)
			}
		}
	}

	fmt.Println("usr list : ", list1)
	fmt.Println("list 2 : ", list2)
	return found, retlist
}
