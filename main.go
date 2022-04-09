package main

import "fmt"

func sumf(c chan int) {
	num1 := <-c
	num2 := <-c 
	c <- num1 + num2
}
func minusf(c chan int) {
	num1 := <-c
	num2 := <-c
	c <- num1 - num2
}
func mpf(c chan int) {
	num1 := <-c
	num2 := <-c
	c <- num1 * num2
}
func dvf(c chan int) {
	num1 := <-c
	num2 := <-c
	c <- num1 / num2
}

func main() {
	fmt.Println("Start")
	sum := make(chan int, 1)
	minus := make(chan int, 1)
	mp := make(chan int, 1)
	dv := make(chan int, 1)

	go sumf(sum)
	go minusf(minus)
	go mpf(mp)
	go dvf(dv)

	sum <- 2
	sum <- 2

	sumval := <-sum

	go sumf(sum)
	sum <- 3
	sum <- 6

	mp <- 7
	mp <- 7

	dv <- 9
	dv <- 3

	sumval2, mpval, dvval := <-sum, <-mp, <-dv

	fmt.Printf("2 + 2 = %v, 3 + 6 = %v, 7 * 7 = %v, 9 / 3 = %v", sumval, sumval2, mpval, dvval)

}