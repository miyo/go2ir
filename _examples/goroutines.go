package main


//var c chan int = make(chan int)
	
func f(s []int) int{
	sum := 0
	for _,v := range s{
		sum += v
	}
	return sum
}

func g(s []int){
	mm := 1
	for _,v := range s{
		mm *= v
	}
	return mm
}


var S0 []int = make([]int, 100)
var S1 []int = make([]int, 100)

func h(){
	

	go f(s0, c)
	go g(s1, c)

}
