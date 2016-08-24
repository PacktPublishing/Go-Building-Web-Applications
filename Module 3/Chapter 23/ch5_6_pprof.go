package main

import
(
"fmt"
"math/rand"
	"runtime/pprof"
	"flag"
	"runtime"
	"os"	
)

var profile = flag.String("cpuprofile", "", "output pprof data to file")

func generateString(length int, seed *rand.Rand, chHater chan string) string  {
	bytes := make([]byte, length)
	for i:=0;i<length;i++ {
		bytes[i] = byte(rand.Int())
	}
	chHater <- string(bytes[:length])
	return string(bytes[:length])
}

func generateChannel() <-chan int {
	ch := make(chan int)
	return ch
}

func main() {
	iterations := 99999
	goodbye := make(chan bool,iterations)
	channelThatHatesLetters := make(chan string)

	runtime.GOMAXPROCS(2)
	flag.Parse()
	if *profile != "" {
		flag,err := os.Create(*profile)
		if err != nil {
			fmt.Println("Could not create profile",err)
		}
		pprof.StartCPUProfile(flag)
		defer pprof.StopCPUProfile()

	}

	seed := rand.New(rand.NewSource(19))	

	initString := ""



	for i:= 0; i < iterations; i++ {
		go func() {
			initString = generateString(300,seed,channelThatHatesLetters)
			goodbye <- true			
		}()

	}
	select {
		case <-channelThatHatesLetters:

	}
	<-goodbye

	fmt.Println(initString)

}