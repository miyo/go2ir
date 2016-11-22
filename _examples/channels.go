package main


//var c chan int = make(chan int)
	
func f(s []int, c chan int){
	sum := 0
	for _,v := range s{
		sum += v
	}
	c <- sum
}

func g(s []int, c chan int){
	sum := 1
	for _,v := range s{
		sum *= v
	}
	c <- sum
}

func h(s []int, c0 chan int, c1 chan int){

	go f(s, c)
	go g(s, c)

}
