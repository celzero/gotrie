package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"unsafe"

	"github.com/celzero/gotrie/trie"
)

//import "time"
//import "os"
//import "bufio"
//import "strings"

func main() {
	mychannel := make(chan bool)
	go loadbuild(mychannel)
	go loadbuild(mychannel)
	go loadbuild(mychannel)
	go loadbuild(mychannel)
	go loadbuild(mychannel)
	<-mychannel
	<-mychannel
	<-mychannel
	<-mychannel
	<-mychannel
}

func loadbuild(mychannel chan bool) {
	FT, err := trie.Build("./td", "./rank", "./basicconfig", "./blocklists")
	if err == nil {
		//[33216 32768 8192 256 4]
		//6IeA6ICA4oCAxIAE
		res := []string{"AMI", "CQT", "EOK", "MTF"}
		usr_flag := FT.CreateUrlEncodedflag(res)
		fmt.Println(usr_flag)
		fmt.Println(FT.Urlenc_to_flag("w4DEgAQ="))

		fmt.Print("1: v1<>google.com (ex:false) ")
		fmt.Println(FT.DNlookup("google.com", "6IeA6ICA4oCAxIAE"))

		// version 0: all blocklists
		t := "77%2Bg77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2Bg"
		fmt.Print("2: v1<>amazon.com (ex:true) ")
		fmt.Println(FT.DNlookup("amazon.com", t))

		// version 1: all blocklists
		l := "1:4P___________________________-D_"
		fmt.Print("3: v2<>amazon.com (ex: true) ")
		fmt.Println(FT.DNlookup("amazon.com", l))
		// version 1: services and native wildcard blocklists
		n := "1:IHD_D___APzAPw=="
		fmt.Print("4: v2<>twitch.tv (ex:true) ")
		fmt.Println(FT.DNlookup("twitch.tv", n))
		fmt.Print("5: v2<>rubbish.twitch.tv (ex:true) ")
		fmt.Println(FT.DNlookup("rubbish.twitch.tv", n))
		fmt.Print("6: v2<>twitter.com (ex:true) ")
		fmt.Println(FT.DNlookup("twitter.com", n))
		fmt.Print("7: v2<>block.what.ever.twitter.com (ex:true) ")
		fmt.Println(FT.DNlookup("block.what.ever.twitter.com", n))
		fmt.Print("8: v2<>www.aws.amazon.com (ex:true) ")
		fmt.Println(FT.DNlookup("www.aws.amazon.com", n)) // true

		fmt.Print("8: v2<>amazon.com (ex:false) ")
		fmt.Println(FT.DNlookup("amazon.com", "1:AIAAwA==")) // false

		PrintMemUsage()

		fmt.Println("ft: %d, td: %d, rd: %d",
			unsafe.Sizeof(FT), unsafe.Sizeof(FT.GetData()), unsafe.Sizeof(FT.GetDir()))
		/*
		   for{
		       reader := bufio.NewReader(os.Stdin)
		       fmt.Print("Enter text: ")
		       text, _ := reader.ReadString('\n')
		       text = strings.TrimSpace(text)
		       if(text == "exit"){break;}
		       start := time.Now()
		       fmt.Println(FT.DNlookup(text,usr_flag))
		       elapsed := time.Since(start)
		       fmt.Printf("Time Diff %s\n",elapsed)
		   }*/
	} else {
		fmt.Println("Error at trie Build")
	}

	/*err,FT := trie.Build()
	  if(err == nil){
	      res := []string{"DAH","ADH","BXW", "BQJ"}
	      fmt.Println("Base64 to flag : ",FT.Urlenc_to_flag(FT.CreateUrlEncodedflag(res)))
	  }*/

	/*
	   start := time.Now()
	   trie.CheckDN1()
	   elapsed := time.Since(start)
	   fmt.Printf("Time Diff %s\n",elapsed)*/
	f, err := os.Create("./amx")
	if err != nil {
		fmt.Println("could not create memory profile: ", err)
	}
	defer f.Close() // error handling omitted for example
	if err := pprof.WriteHeapProfile(f); err != nil {
		fmt.Println("could not write memory profile: ", err)
	}
	mychannel <- true
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
