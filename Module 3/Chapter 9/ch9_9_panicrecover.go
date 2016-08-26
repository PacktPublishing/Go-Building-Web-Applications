package main


import
(
	"os"
	"fmt"
	"strconv"
)

func gatherPanics() {
	if rec := recover(); rec != nil {
		fmt.Println("Critical Error:", rec)
	}
}

func getFileDetails(fileName string) {
	defer gatherPanics()
	finfo,err := os.Stat(fileName)	
	if err != nil {	
		panic("Cannot access file")
	}else {
		fmt.Println("Size: ", strconv.FormatInt(finfo.Size(),10))
	}
}

func openFile(fileName string) {
	defer gatherPanics()
	if _, err := os.Stat(fileName); err != nil {
		panic("File does not exist")
	}	

}

func main() {
	var fileName string
	fmt.Print("Enter filename>")
	_,err := fmt.Scanf("%s",&fileName)
	if err != nil {}
	fmt.Println("Getting info for",fileName)

	openFile(fileName)
	getFileDetails(fileName)

}