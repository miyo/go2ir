package main

import "fmt"

var c chan int = make(chan int)
	
func sum(s []int, c chan int){
	sum := 0
	for _,v := range s{
		sum += v
	}
	c <- sum
}

func main(){
	s := []int{7,2,8,-9,4,0}


	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)

	x, y := <-c, <-c

	fmt.Println(x, y, x + y)
}
