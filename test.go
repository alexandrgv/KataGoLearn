package main

// Массив - сумма чётных и нечётных элементов

func handleArray(arr []int) (int, int) {
  var i, j int

  for _, el := range arr {
    if el%2 == 0 { // целое
      i += el
    } else { // чётное
      j += el
    }
  }
  return i, j
}

func main() {
  arr := []int{1, 2, 3, 4, 5, 6, 7}
  println(handleArray(arr))
}
