package main

func main() {

}

type Adder interface {
	Add(a, b int) int
}

func Calculate(a, b int, adder Adder) int {
	ans := adder.Add(a, b)
	return ans
}

func Add(a, b int) int {
	return a + b
}
