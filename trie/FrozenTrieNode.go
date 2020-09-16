package trie
import "fmt"
type FrozenTrieNode struct{
	trie FrozenTrie
	index int
	valCached *[]uint32
	finCached,comCached,flagCached *bool
	whCached *uint32
	fcCached,chCached *int
}

func (FTN *FrozenTrieNode)Init(FT FrozenTrie,index int)  {
	FTN.trie = FT
	FTN.index = index
	if (Debug) {
		fmt.Printf("%d :i, fc: %d tl: %d c: %t f: %t wh: %d flag: %t\n",FTN.index ,FTN.firstChild(),FTN.letter(), FTN.compressed(),FTN.final(), FTN.where() , FTN.flag())		
    }
}



func (FTN FrozenTrieNode)check(){
	fmt.Println("FrozenTrieNode td length : ",len(FTN.trie.data.bytes))
	//fmt.Println("FrozenTrieNode rd length : ",len(FTN.trie.directory))
}


func (FTN *FrozenTrieNode)final()bool{	
	if(FTN.finCached == nil){
		tmp := (FTN.trie.data.get( FTN.trie.letterStart + (FTN.index * FTN.trie.bitslen) + FTN.trie.extraBit, 1,false ) == 1)
		FTN.finCached = &tmp
	}
	return *FTN.finCached
}

func (FTN *FrozenTrieNode)where()uint32{
	if(FTN.whCached == nil){
		tmp := FTN.trie.data.get(FTN.trie.letterStart + (FTN.index * FTN.trie.bitslen) + 1 + FTN.trie.extraBit, FTN.trie.bitslen - 1 - FTN.trie.extraBit,false);
		FTN.whCached = &tmp
	}
	return *FTN.whCached
}

func (FTN *FrozenTrieNode)compressed()bool{
	if(FTN.comCached == nil){
		tmp := (FTN.trie.data.get(FTN.trie.letterStart + (FTN.index * FTN.trie.bitslen), 1,false) == 1) //(config.compress && !config.unroll) 
		FTN.comCached = &tmp
	}
	return *FTN.comCached
}

func (FTN *FrozenTrieNode)flag()bool{
	if(FTN.flagCached == nil){
		tmp := (FTN.compressed() && FTN.final()) //(config.valueNode) ?
		FTN.flagCached = &tmp
	}
	return *FTN.flagCached
}

func (FTN *FrozenTrieNode)letter()uint32{
	return FTN.where()
} 

func (FTN *FrozenTrieNode)firstChild()int{
	if(FTN.fcCached == nil){
		tmp := FTN.trie.directory.selectRD( 0, FTN.index+1 ) - FTN.index
		FTN.fcCached = &tmp
	}
	return *FTN.fcCached
}

func (FTN *FrozenTrieNode)childOfNextNode()int{
	if(FTN.chCached == nil){
		tmp :=  FTN.trie.directory.selectRD( 0, FTN.index + 2 ) - FTN.index - 1
		FTN.chCached = &tmp
	}
	return *FTN.chCached
}

func (FTN *FrozenTrieNode)childCount()int{
	return FTN.childOfNextNode() - FTN.firstChild();
}

func (FTN *FrozenTrieNode)value()[]uint32{
	if (FTN.valCached == nil) {
		//let valueChain = this;
		value := []uint32{};
		i := 0;
		j := 0;
		//if (config.debug) console.log("thisnode: index/vc/ccount ", this.index, this.letter(), this.childCount())
		for (i < FTN.childCount()) {
			valueChain := FTN.getChild(i);
			//if (config.debug) console.log("vc no-flag end vlet/vflag/vindex/val ", i, valueChain.letter(), valueChain.flag(), valueChain.index, value)
			if (!valueChain.flag()) {
				break;
			}
			if (i % 2 == 0) {
				value = append(value,valueChain.letter() << 8);
			} else {
				value[j] = (value[j] | valueChain.letter());
				j += 1;
			}
			i += 1;
		}
		FTN.valCached = &value;
	}

	return *FTN.valCached;
}

func (FTN *FrozenTrieNode)getChildCount()int{
	return FTN.childCount();
}

func (FTN *FrozenTrieNode)getChild(index int)FrozenTrieNode{
	return FTN.trie.getNodeByIndex(FTN.firstChild() + index);
}