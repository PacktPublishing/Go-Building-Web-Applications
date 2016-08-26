package main

import(
	"net/http"
	"html/template"
	"fmt"
	"github.com/gorilla/mux"
	"sync"
	"strconv"
)

type User struct {
	Name string
	Times map[int] bool
	DateHTML template.HTML
}

type Page struct {
	Title string
	Body template.HTML
	Users map[string] User
}

var usersInit map[string] bool
var userIndex int
var validTimes []int
var mutex sync.Mutex
var Users map[string]User
var templates = template.Must(template.New("template").ParseFiles("view_users.html", "register.html"))


func register(w http.ResponseWriter, r *http.Request){
	fmt.Println("Request to /register")
	params := mux.Vars(r)
 	name := params["name"]	

	if _,ok := Users[name]; ok {
		t,_ := template.ParseFiles("generic.txt")
		page := &Page{ Title: "User already exists", Body: template.HTML("User " + name + " already exists")}
		t.Execute(w, page)
 	} 	else {
 		newUser := User { Name: name }
 		initUser(&newUser)
 		Users[name] = newUser
	 	t,_ := template.ParseFiles("generic.txt")
		page := &Page{ Title: "User created!", Body: template.HTML("You have created user "+name)}
		t.Execute(w, page)		
 	}



}

func dismissData(st1 int, st2 bool) {


}

func formatTime(hour int) string {
	hourText := hour
	ampm := "am"
	if (hour > 11) {
		ampm = "pm"
	}
	if (hour > 12) {
		hourText = hour - 12;
	}
fmt.Println(ampm)
	outputString := strconv.FormatInt(int64(hourText),10) + ampm
	

	return  outputString
}

func (u User) FormatAvailableTimes() template.HTML {
  HTML := ""
  HTML +=  "<b>"+u.Name+"</b> - "

  for k,v := range u.Times {
  	dismissData(k,v)

  	if (u.Times[k] == true) {
  		formattedTime := formatTime(k)
  		HTML += "<a href='/schedule/"+u.Name+"/"+strconv.FormatInt(int64(k),10)+"' class='button'>"+formattedTime+"</a> "

  		} else {

  		}

  }
  return template.HTML(HTML)
}

func users(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to /users")



	t,_ := template.ParseFiles("users.txt")
	page := &Page{ Title: "View Users", Users: Users}
	t.Execute(w, page)
}

func schedule(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request to /schedule")
	params := mux.Vars(r)
 	name := params["name"]	
 	time := params["hour"]
 	timeVal,_ := strconv.ParseInt( time, 10, 0 )
 	intTimeVal := int(timeVal)

 	createURL := "/register/"+name

	if _,ok := Users[name]; ok {
		if Users[name].Times[intTimeVal] == true {
			mutex.Lock()
				Users[name].Times[intTimeVal] = false
			mutex.Unlock()				
				fmt.Println("User exists, variable should be modified")
				t,_ := template.ParseFiles("generic.txt")
				page := &Page{ Title: "Successfully Scheduled!", Body: template.HTML("This appointment has been scheduled. <a href='/users'>Back to users</a>")}

				t.Execute(w, page)
		
		}	else {
			fmt.Println("User exists, spot is taken!")
			t,_ := template.ParseFiles("generic.txt")
			page := &Page{ Title: "Booked!", Body: template.HTML("Sorry, "+name+" is booked for "+time+" <a href='/users'>Back to users</a>")}
			t.Execute(w, page)

		}

 	} 	else {
 		fmt.Println("User does not exist")		
		t,_ := template.ParseFiles("generic.txt")
		page := &Page{ Title: "User Does Not Exist!", Body: template.HTML( "Sorry, that user does not exist. Click <a href='"+createURL+"'>here</a> to create it.  <a href='/users'>Back to users</a>")}
		t.Execute(w, page)
 	}
 	fmt.Println(name,time)
}

func defaultPage(w http.ResponseWriter, r *http.Request) {

}

func initUser(user *User) {

	user.Times = make(map[int] bool)
	for i := 9; i < 18; i ++ {
		user.Times[i] = true
	}

}

func main() {
	Users = make(map[string] User)
	userIndex = 0
	bill := User {	Name: "Bill"	}
	initUser(&bill)
	Users["Bill"] = bill
	userIndex++

    r := mux.NewRouter()
    r.HandleFunc("/", defaultPage)
    r.HandleFunc("/users", users)
    r.HandleFunc("/register/{name:[A-Za-z]+}", register)   
    r.HandleFunc("/schedule/{name:[A-Za-z]+}/{hour:[0-9]+}", schedule)      
    http.Handle("/", r)

    err := http.ListenAndServe(":1900", nil)
    if err != nil {
       // log.Fatal("ListenAndServe:", err)
    }	


}
