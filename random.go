package main

import (

	"math/rand"
	"fmt"
)



func get_random(low int, high int) int {
	return rand.Intn(high-low+1) + low
}

func jason_random_v1(low int, high int) int {
	get_times := get_random(1, 5)
	var value int
	for i := 0; i < get_times; i++ {
		value = get_random(low, high)
	}
	return value
}




func main(){

	fmt.Println(jason_random_v1(1,10))
	fmt.Println(jason_random_v1(1,10))
	fmt.Println(jason_random_v1(1,10))
	fmt.Println(jason_random_v1(1,10))
	fmt.Println(jason_random_v1(1,10))

}
