package main

import
(
	"fmt"
	"github.com/howeyc/fsnotify"
	"net"
	"time"
	"io"	
	"io/ioutil"
	"github.com/couchbaselabs/go-couchbase"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"	
	"strings"
	
)

var listenFolder = "/wamp/www/shared"

type Client struct {
	ID int
	Connection *net.Conn	
}

type File struct {
	Hash string "json:hash"
	Name string "json:file_name"
	Created int64 "json:created"
	CreatedUser  int "json:created_user"
	LastModified int64 "json:last_modified"
	LastModifiedUser int "json:last_modified_user"
	Revisions int "json:revisions"
	Version int "json:version"
}

type Message struct {
	Hash string "json:hash"
	Action string "json:action"
	Location string "json:location"	
	Name string "json:name"
	Version int "json:version"
}

func generateHash(name string) string {

	hash := md5.New()
	io.WriteString(hash,name)
	hashString := hex.EncodeToString(hash.Sum(nil))

	return hashString
}

func alertServers(hash string, name string, action string, location string, version int) {

	msg := Message{Hash:hash,Action:action,Location:location,Name:name,Version:version}
	msgJSON,_ := json.Marshal(msg)

	fmt.Println(string(msgJSON))

	for i := range Clients {
		fmt.Println("Sending to clients")
		fmt.Fprintln(*Clients[i].Connection,string(msgJSON))
	}
}

func startServer(listener net.Listener) {
	for {	
		conn,err := listener.Accept()
		if err != nil {

		}
		currentClient := Client{ ID: 1, Connection: &conn}
		Clients = append(Clients,currentClient)
	    for i:= range Clients {
	    	fmt.Println("Client",Clients[i].ID)
	    }		
	}	

}

func removeFile(name string, bucket *couchbase.Bucket) {
	bucket.Delete(generateHash(name))
}

func updateExistingFile(name string, bucket *couchbase.Bucket) int {
	fmt.Println(name,"updated")
	hashString := generateHash(name)
	
	thisFile := Files[hashString]
	thisFile.Hash = hashString
	thisFile.Name = name
	thisFile.Version = thisFile.Version + 1
	thisFile.LastModified = time.Now().Unix()
	Files[hashString] = thisFile
	bucket.Set(hashString,0,Files[hashString])
	return thisFile.Version
}

func evalFile(event *fsnotify.FileEvent, bucket *couchbase.Bucket) {
	fmt.Println(event.Name,"changed")
	create := event.IsCreate()
	fileComponents := strings.Split(event.Name,"\\")
	fileComponentSize := len(fileComponents)
	trueFileName := fileComponents[fileComponentSize-1]
	hashString := generateHash(trueFileName)

	if create == true {
		updateFile(trueFileName,bucket)
		alertServers(hashString,event.Name,"CREATE",event.Name,0)
	}
	delete := event.IsDelete()
	if delete == true {
		removeFile(trueFileName,bucket)
		alertServers(hashString,event.Name,"DELETE",event.Name,0)		
	}
	modify := event.IsModify()
	if modify == true {
		newVersion := updateExistingFile(trueFileName,bucket)
		fmt.Println(newVersion)
		alertServers(hashString,trueFileName,"MODIFY",event.Name,newVersion)
	}
	rename := event.IsRename()
	if rename == true {


	}
}

func updateFile(name string, bucket *couchbase.Bucket) {
	thisFile := File{}
	hashString := generateHash(name)
	
	thisFile.Hash = hashString
	thisFile.Name = name
	thisFile.Created = time.Now().Unix()
	thisFile.CreatedUser = 0
	thisFile.LastModified = time.Now().Unix()
	thisFile.LastModifiedUser = 0
	thisFile.Revisions = 0
	thisFile.Version = 1

	Files[hashString] = thisFile

	checkFile := File{}
	err := bucket.Get(hashString,&checkFile)
	if err != nil {
		fmt.Println("New File Added",name)
		bucket.Set(hashString,0,thisFile)
	}
}

var Clients []Client
var Files map[string] File


func main() {
	Files = make(map[string]File)
	endScript := make(chan bool)

	couchbaseClient, err := couchbase.Connect("http://localhost:8091/")
		if err != nil {
			fmt.Println("Error connecting to Couchbase", err)
		}
	pool, err := couchbaseClient.GetPool("default")
		if err != nil {
			fmt.Println("Error getting pool",err)
		}
	bucket, err := pool.GetBucket("file_manager")
		if err != nil {
			fmt.Println("Error getting bucket",err)
		}	


	files, _ := ioutil.ReadDir(listenFolder)
	for _, file := range files {
		updateFile(file.Name(),bucket)
	}

    dirSpy, err := fsnotify.NewWatcher()
    defer dirSpy.Close()

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println ("Could not start server!",err)
	}

	go func() {
        for {
            select {
            case ev := <-dirSpy.Event:
                evalFile(ev,bucket)
            case err := <-dirSpy.Error:
                fmt.Println("error:", err)
            }
        }
    }()
    err = dirSpy.Watch(listenFolder)	
	startServer(listener)

	<-endScript
}