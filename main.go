package main

import "github.com/celzero/gotrie/trie"
import "fmt"
import "runtime"
import "unsafe"
import (
    "os"
    "runtime/pprof"
)

//import "time"
//import "os"
//import "bufio"
//import "strings"

func main() {

    err, FT := trie.Build("./td", "./rank", "./basicconfig", "./blocklists")
    if err == nil {
        //[33216 32768 8192 256 4]
        //6IeA6ICA4oCAxIAE
        res := []string{"AMI","CQT","EOK","MTF"}
        usr_flag := FT.CreateUrlEncodedflag(res)
        fmt.Println(usr_flag)
        fmt.Println(FT.Urlenc_to_flag("w4DEgAQ="))

        fmt.Print("1: v1<>google.com ")
        fmt.Println(FT.DNlookup("google.com", "6IeA6ICA4oCAxIAE"))

        // version 0: all blocklists
        t := "77%2Bg77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2B%2F77%2Bg"
        fmt.Print("2: v1<>amazon.com ")
        fmt.Println(FT.DNlookup("amazon.com", t))

        // version 1: all blocklists
        l := "1:4P___________________________-D_"
        fmt.Print("3: v2<>amazon.com ")
        fmt.Println(FT.DNlookup("amazon.com", l))
        // version 1: services and native wildcard blocklists
        n := "1:IHD_D___APzAPw=="
        fmt.Println(FT.DNlookup("twitch.tv", n))
        fmt.Println(FT.DNlookup("rubbish.twitch.tv", n))
        fmt.Println(FT.DNlookup("twitter.com", n))
        fmt.Println(FT.DNlookup("block.what.ever.twitter.com", n))

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
