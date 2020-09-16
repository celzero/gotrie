package trie

import "math"
import "fmt"
type RankDirectory struct{
	Directory BS
	ValueDir BS
	Data BS
	l1Size int
	l2Size int
	l1Bits int
	l2Bits int
	sectionBits int
	numBits int
}

func (RD *RankDirectory) Init(directoryData []uint16, trieData []uint16, numBits int, l1Size int, l2Size int, valueDir []uint16)  {
	RD.Directory = BS{}
	RD.Directory.Init(directoryData)

	if(valueDir != nil){
		RD.ValueDir = BS{}
		RD.ValueDir.Init(valueDir)
	}
	RD.Data = BS{}
	RD.Data.Init(trieData)

	RD.l1Size = l1Size;
	RD.l2Size = l2Size;
	RD.l1Bits = int(math.Ceil(math.Log2(float64(numBits))));
	RD.l2Bits = int(math.Ceil(math.Log2(float64(l1Size))));
	RD.sectionBits = (l1Size / l2Size - 1) * RD.l2Bits + RD.l1Bits;
	RD.numBits = numBits;
}

func (RD RankDirectory) display(){
	fmt.Println("RankDirectory rd length : ",len(RD.Directory.bytes))
	fmt.Println("RankDirectory td length : ",len(RD.Data.bytes))
	fmt.Println("Num Bits : ",RD.numBits)
	fmt.Println("L1size : ",RD.l1Size)
	fmt.Println("l1bits : ",RD.l1Bits)
	fmt.Println("l2bits : ",RD.l2Bits)
}

func (RD *RankDirectory) selectRD(which int,y int)int{
	which = 0
	return RD.rank(0, y); //if (config.selectsearch) { }	
}


func (RD *RankDirectory) rank(which int,x int)int{
	//if (config.selectsearch) { }
	var temp uint32
	rank := -1;
	sectionPos := 0;
	if (x >= RD.l2Size) {
		sectionPos = (x / RD.l2Size | 0) * RD.l1Bits
		temp = RD.Directory.get(sectionPos - RD.l1Bits, RD.l1Bits,false)
		rank = int(temp);
		x = x % RD.l2Size;
	}
	var ans = 0
	if(x > 0){
		ans = RD.Data.pos0( rank + 1, x )
	}else{
		ans = rank
	}
	//if (config.debug) console.log("ans: " + ans + " " + rank + ":r, x: " + x + " " + sectionPos + ":s, o: " + o);
	if(Debug){
		fmt.Printf("ans: %d %d:r, x: %d %d:s %d:l1 %t:ifcheck\n",ans,temp,x,sectionPos,RD.l1Bits,x >= RD.l2Size)
	}
	return ans;	
}