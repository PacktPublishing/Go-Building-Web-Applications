package main

import (
"os"
"flag"
"fmt"
"runtime/pprof"
)

const TESTLENGTH = 100000
type CPUHog struct {
	longByte []byte
}

func makeLongByte() []byte {
	longByte := make([]byte,TESTLENGTH)

	for i:= 0; i < TESTLENGTH; i++ {
		longByte[i] = byte(i)
	}
	return longByte
}

var profile = flag.String("cpuprofile", "", "output pprof data to file")


func main() {
	var CPUHogs []CPUHog

	flag.Parse()
		if *profile != "" {
			flag,err := os.Create(*profile)
			if err != nil {
				fmt.Println("Could not create profile",err)
			}
			pprof.StartCPUProfile(flag)
			defer pprof.StopCPUProfile()

		}

	for i := 0; i < TESTLENGTH; i++ {
		hog := CPUHog{}
		hog.longByte = makeLongByte()
		_ = append(CPUHogs,hog)
	}
}
