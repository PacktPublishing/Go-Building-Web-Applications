package main

import
(
  "fmt"
)


type User struct {
  UserName string
  UserID int64
}

type SharedFile struct {
  FileName string
  LastModified int64
  LastModifiedUser User
  Version int64
  Contents string
}

func commitFile(data string) {

}

func versionFile(sf SharedFile) {

}

func getUserDetails(userName string) {

}

func getFileDetails(fileName string) {

}

func startWebServer() {

}

func startCLIServer() {

}

func main() {

  fmt.Println("File Sharing Web Interface Starting")
  startWebServer()
  startCLIServer()
}
