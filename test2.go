package main

//2 переменные диап 012, сложить и перевести в римское число

func arToRim(i int) string {

	var ar = map[int]string{
		1: "I",
		2: "II",
		3: "III",
		4: "IV",
	}

	var _ = map[string]int{
		"I":   1,
		"II":  2,
		"III": 3,
		"IV":  4,
	}

	str := ar[i]

	return str
}

func main() {
	a := 0
	b := 1
	sum := a + b
	sumR := arToRim(sum)
	println(sumR)

}
