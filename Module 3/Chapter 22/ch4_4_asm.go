package main

/*
#include <stdio.h>

void asm() {

__asm__( "" );
    printf("I come from a %s","C function with embedded asm\n");

}
*/
import "C"

func main() {
    
    C.asm()

}