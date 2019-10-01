package main

import (
	"fmt"
)

func main() {
	var s1 = [2]*int{new(int),new(int)}
	*s1[0] = 2
	*s1[1] = 4
	fmt.Println(s1)
	var s2 = make([]int,2)
	for _,v := range s1{
		s2 = append(s2,*v)
		
	}
	fmt.Println(s2)
}

type Runner interface {
	Run()
}

type run1 struct {

}

func (p *run1)Run()  {
	fmt.Println("run1")
}

type run2 struct {

}

func (p run2)Run()  {
	fmt.Println("run2")
}