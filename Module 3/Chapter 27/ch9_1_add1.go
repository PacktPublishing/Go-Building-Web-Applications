package main

import
(
	"errors"
	"reflect"
	"log"
)

type Alpha struct {

}

type Numeric struct {

}

func (a Alpha) Add(x string, y string) (string, error) {
	var err error
	xType := reflect.TypeOf(x).Kind()
	yType := reflect.TypeOf(y).Kind()
	if xType != reflect.String || yType != reflect.String {
		err = errors.New("Incorrect type for strings a or b!")
	}
	finalString := x + y
	return finalString, err
}

func (n Numeric) Add(x int, y int) (int, error) {
	var err error

	xType := reflect.TypeOf(x).Kind()
	yType := reflect.TypeOf(y).Kind()
	if xType != reflect.Int || yType != reflect.Int {
		err = errors.New("Incorrect type for integer a or b!")
	}
	return x + y, err
}

func main() {
	n1 := Numeric{}
	a1 := Alpha{}
	z,err := n1.Add(5,2)	
	if err != nil {
		log.Println("Error",err)
	}
	log.Println(z)

	y,err := a1.Add("super","lative")
	if err != nil {
		log.Println("Error",err)
	}
	log.Println(y)
}