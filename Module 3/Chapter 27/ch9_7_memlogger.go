package main

import
(
	"github.com/couchbaselabs/go-couchbase"
	"io"	
	"time"
	"fmt"
	"os"
	"net/http"
	"crypto/md5"
	"encoding/hex"	
)
type LogItem struct {
	ServerID string "json:server_id"
	Goroutine int "json:goroutine"
	Timestamp time.Time "json:time"
	Message string "json:message"
	Page string "json:page"
}

var currentGoroutine int

func (li LogItem) logRequest(bucket *couchbase.Bucket) {

	hash := md5.New()
	io.WriteString(hash,li.ServerID+li.Page+li.Timestamp.Format("Jan 1, 2014 12:00am"))
	hashString := hex.EncodeToString(hash.Sum(nil))
	bucket.Set(hashString,0,li)
	currentGoroutine = 0
}

func main() {
	hostName, _ := os.Hostname()
	currentGoroutine = 0
	
	logClient, err := couchbase.Connect("http://localhost:8091/")
		if err != nil {
			fmt.Println("Error connecting to logging client", err)
		}
	logPool, err := logClient.GetPool("default")
		if err != nil {
			fmt.Println("Error getting pool",err)
		}
	logBucket, err := logPool.GetBucket("logs")
		if err != nil {
			fmt.Println("Error getting bucket",err)
		}	

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		request := LogItem{}
		request.Goroutine = currentGoroutine
		request.ServerID = hostName
		request.Timestamp = time.Now()
		request.Message = "Request to " + r.URL.Path
		request.Page = r.URL.Path
		go request.logRequest(logBucket)

	})		

	http.ListenAndServe(":8080",nil)

}