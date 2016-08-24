package main

import 
(
	"flag"
	"fmt"
	"os"
	"github.com/couchbaselabs/go-couchbase"	
	"io"	
	"crypto/md5"
	"encoding/hex"
	"strconv"	
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

var liveFolder = "/wamp/www/shared/"
var backupFolder = "/wamp/www/backup/"

func generateHash(name string) string {

	hash := md5.New()
	io.WriteString(hash,name)
	hashString := hex.EncodeToString(hash.Sum(nil))

	return hashString
}

func main() {
	revision := flag.Int("r",0,"Number of versions back")
	fileName := flag.String("f","","File Name")
	flag.Parse()

	if *fileName == "" {

		fmt.Println("Provide a file name to use!")
		os.Exit(0)
	}


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

	hashString := generateHash(*fileName)
	checkFile := File{}		
	bucketerr := bucket.Get(hashString,&checkFile)
	if bucketerr != nil {

	}else {
		backupLocation := backupFolder + checkFile.Name + "." + strconv.FormatInt(int64(checkFile.Version-*revision),10)
		newLocation := liveFolder + checkFile.Name
		fmt.Println(backupLocation)
		org,_ := os.Open(backupLocation)
			defer org.Close()
		cpy,_ := os.Create(newLocation)
			defer cpy.Close()
		io.Copy(cpy,org)
		fmt.Println("Revision complete")
	}



}