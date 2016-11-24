package main

//var c0 chan int = make(chan int)
//var c1 chan int = make(chan int)

var s0 []int = make([]int, 16)
var s1 []int = make([]int, 16)

var a int = 0
var b int = 0

func f(){
	sum := 0
	for _,v := range s0{
		sum += v
	}
	a = sum
}

func g(){
	mm := 1
	for _,v := range s1{
		mm *= v
	}
	b = mm
}

func h(){
	go f()
	go g()
}
