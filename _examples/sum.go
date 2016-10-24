package main

func sum(s []int, c chan int){
	sum := 0
	for _,v := range s{
		sum += v
	}
	c <- sum
}

