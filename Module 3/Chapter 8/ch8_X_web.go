package main

import
(
	"net"
	"net/http"
	"html/template"
	"log"
	"io"
	"os"
	"io/ioutil"
	"github.com/couchbaselabs/go-couchbase"
	"time"	
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)



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

type Page struct {
	Title string
	Files map[string] File
}

type ItemWrapper struct {

	Items []File
	CurrentTime int64
	PreviousTime int64

}

type Message struct {
	Hash string "json:hash"
	Action string "json:action"
	Location string "json:location"
	Name string "json:name"	
	Version int "json:version"
}

var listenFolder = "/wamp/www/shared/"
var Files map[string] File
var webTemplate = template.Must(template.ParseFiles("ch8_html.html"))
var fileChange chan File
var lastChecked int64

func generateHash(name string) string {

	hash := md5.New()
	io.WriteString(hash,name)
	hashString := hex.EncodeToString(hash.Sum(nil))

	return hashString
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
	}else {
		Files[hashString] = checkFile
	}
}

func listen(conn net.Conn) {
	for {

	    messBuff := make([]byte,1024)
		n, err := conn.Read(messBuff)
		if err != nil {

		}
		message := string(messBuff[:n])
		message = message[0:]

		resultMessage := Message{}
		json.Unmarshal(messBuff[:n],&resultMessage)
		
		updateHash := resultMessage.Hash
		tmp := Files[updateHash]
		tmp.LastModified = time.Now().Unix()
		Files[updateHash] = tmp
	}

}

func main() {
	lastChecked := time.Now().Unix()
	Files = make(map[string]File)
	fileChange = make(chan File)
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

	conn, err := net.Dial("tcp","127.0.0.1:9000")
	if err != nil {
		fmt.Println("Could not connect to File Listener!")
	}
	go listen(conn)


	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		apiOutput := ItemWrapper{}
		apiOutput.PreviousTime = lastChecked
		lastChecked = time.Now().Unix()
		apiOutput.CurrentTime = lastChecked

		for i:= range Files {
			apiOutput.Items = append(apiOutput.Items,Files[i])
		}
		output,_ := json.Marshal(apiOutput)
		fmt.Fprintln(w,string(output))

	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output := Page{Files:Files,Title:"File Manager"}
		tmp, _ := template.ParseFiles("ch8_html.html")
		tmp.Execute(w, output)
	})
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10000000)
		if err != nil {
			return
		}
		form := r.MultipartForm

		files := form.File["file"]
		for i, _ := range files {
			newFileName := listenFolder + files[i].Filename
			org,_:= files[i].Open()
			defer org.Close()
			cpy,_ := os.Create(newFileName)
			defer cpy.Close()
			io.Copy(cpy,org)
		}
	})	

	log.Fatal(http.ListenAndServe(":8080",nil))

}