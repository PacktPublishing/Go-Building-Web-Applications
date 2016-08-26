package main

/*
	#include <stdio.h>
	#include <stdlib.h>

	void Coutput (char* str) {
		printf("%s",str);
	}
*/
import "C"
import "unsafe"

func main() {
	v := C.CString("Don't Forget My Memory Is Not Visible To Go!")
	C.Coutput(v)
	C.free(unsafe.Pointer(v))
}