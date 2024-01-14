package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Структура Calculator, содержащая методы
type Calculator struct{}

// Парсинг ввода
func (c *Calculator) parseInput(input string) (int, int, string, string, error) {

	// Удаляем все пробельные символы из входной строки
	input = strings.ReplaceAll(strings.TrimSpace(input), " ", "")

	// Проверяем наличие математической операции и находим индекс оператора
	operatorIndex := -1
	for i, char := range input {
		if strings.ContainsRune("+-*/", char) {
			operatorIndex = i
			break
		}
	}

	// Выдача паники, если строка не является математической операцией
	if operatorIndex == -1 {
		panic("строка не является математической операцией")
	}

	// Проверяем формат чисел и определяем тип
	numberType, err := c.detectNumberType(input)
	if err != nil {
		panic(err)
	}

	// Разбиваем числа и оператор на отдельные переменные
	firstStr := input[:operatorIndex]        // строка до оператора
	secondStr := input[operatorIndex+1:]     // строка после оператора
	operator := string(input[operatorIndex]) // оператор

	// Проверка на десятичные числа
	if strings.Contains(firstStr, ".") || strings.Contains(secondStr, ".") {
		return 0, 0, "", "", errors.New("Ошибка: калькулятор умеет работать только с целыми числами!")
	}

	var firstNumber, secondNumber int

	// Получаем целые числа
	switch numberType {
	case "arabic":
		// Разбираем числа
		firstNumber, err = strconv.Atoi(firstStr)
		if err != nil {
			return 0, 0, "", "", errors.New("Ошибка: неверный формат первого операнда")
		}
		secondNumber, err = strconv.Atoi(secondStr)
		if err != nil {
			return 0, 0, "", "", errors.New("Ошибка: неверный формат второго операнда")
		}
	case "roman":
		// Переводим римские числа в арабские
		firstNumber, err = c.romanToArabic(firstStr)
		if err != nil {
			return 0, 0, "", "", err
		}
		secondNumber, err = c.romanToArabic(secondStr)
		if err != nil {
			return 0, 0, "", "", err
		}
	default:
		panic("невозможно определить тип чисел")
	}

	// Проверка диапазона чисел
	if firstNumber < 1 || firstNumber > 10 || secondNumber < 1 || secondNumber > 10 {
		return 0, 0, "", "", errors.New("Ошибка: числа должны быть в диапазоне от 1 до 10 включительно")
	}

	return firstNumber, secondNumber, operator, numberType, nil
}

// Метод определения типа цифр (арабские или римские)
func (c *Calculator) detectNumberType(input string) (string, error) {
	// Удаляем пробельные символы из входной строки
	input = strings.ReplaceAll(strings.TrimSpace(input), " ", "")

	r1 := "[0-9IVXLCDM]" // операнд
	r2 := "[-+*/]"       // оператор
	r3 := "[0-9]"        // арифметический операнд
	r4 := "[IVXLCDM]"    // римский операнд

	// Проверяем, соответствует ли количество операндов и операторов требуемому количеству
	isManyOperandsAndOperators := !regexp.MustCompile(r2+"{1}").MatchString(input) || !regexp.MustCompile("^"+r1+"+"+r2+r1+"+$").MatchString(input)

	// Проверяем, используются ли одновременно арабские и римские цифры
	isBoth := regexp.MustCompile(r3+"+").MatchString(input) && regexp.MustCompile(r4+"+").MatchString(input)

	// Проверяем, соответствует ли строка арабским цифрам с помощью регулярного выражения
	isArabic := regexp.MustCompile("^" + r3 + "+" + r2 + r3 + "+$").MatchString(input)

	// Проверяем, соответствует ли строка римским цифрам с помощью регулярного выражения
	isRoman := regexp.MustCompile("^" + r4 + "+" + r2 + r4 + "+$").MatchString(input)

	// Если нарушено количество операндов и операторов
	if isManyOperandsAndOperators {
		panic("необходимо использовать только два операнда и один оператор (+, -, /, *)")
	}

	// Если оба условия выполнены одновременно, возвращаем ошибку
	if isBoth {
		panic("использование одновременно арабских и римских цифр")
	}

	// Если только арабские цифры
	if isArabic {
		return "arabic", nil
	}

	// Если только римские цифры
	if isRoman {
		return "roman", nil
	}

	// Если ни одно из условий не выполнено
	panic("невозможно определить тип чисел")
}

func (c *Calculator) romanToArabic(input string) (int, error) {
	// Создаем карту, сопоставляющую римские цифры с их арабскими значениями
	romanNumerals := map[string]int{
		"I":    1,
		"II":   2,
		"III":  3,
		"IV":   4,
		"V":    5,
		"VI":   6,
		"VII":  7,
		"VIII": 8,
		"IX":   9,
		"X":    10,
	}

	return romanNumerals[input], nil
}

func (c *Calculator) arabicToRoman(number int) (string, error) {
	// Проверка римских чисел на положительность
	if number <= 0 {
		panic("результатом работы калькулятора с римскими числами могут быть только положительные числа")
	}

	// Создаем карту, сопоставляющую арабские значения с их римскими представлениями
	arabicNumerals := map[int]string{
		1:   "I",
		4:   "IV",
		5:   "V",
		9:   "IX",
		10:  "X",
		40:  "XL",
		50:  "L",
		90:  "XC",
		100: "C",
	}

	// Создаем слайс для хранения ключей в убывающем порядке
	var keys []int
	for key := range arabicNumerals {
		keys = append(keys, key)
	}
	// Сортируем слайс ключей в убывающем порядке
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	var result string

	if number > 0 {
		// Перебираем каждую арабскую цифру в обратном порядке, для этого используем упорядоченные ключи для итерации в убывающем порядке
		for _, key := range keys {
			for number >= key {
				result += arabicNumerals[key]
				number -= key
			}
		}
	}

	// Возвращаем соответствующее римское представление
	return result, nil
}

// Методы арифметических операций
func (c *Calculator) add(a, b int) (int, error) {
	return a + b, nil
}

func (c *Calculator) subtract(a, b int) (int, error) {
	return a - b, nil
}

func (c *Calculator) multiply(a, b int) (int, error) {
	return a * b, nil
}

func (c *Calculator) divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("Ошибка: деление на ноль")
	}
	return a / b, nil
}

// Метод Calculate
func (c *Calculator) Calculate(input string) (string, error) {
	firstNumber, secondNumber, operator, numberType, err := c.parseInput(input)
	if err != nil {
		return "", err
	}

	var result int

	// Выполняем математическую операцию
	switch operator {
	case "+":
		result, err = c.add(firstNumber, secondNumber)
	case "-":
		result, err = c.subtract(firstNumber, secondNumber)
	case "*":
		result, err = c.multiply(firstNumber, secondNumber)
	case "/":
		result, err = c.divide(firstNumber, secondNumber)
	default:
		return "", errors.New("Ошибка: неверный оператор")
	}

	if err != nil {
		return "", err
	}

	// Преобразуем число в строку
	switch numberType {
	case "arabic":
		return strconv.Itoa(result), nil
	case "roman":
		resultStr, _ := c.arabicToRoman(result)
		return resultStr, nil
	}

	return "", err
}

func main() {
	c := &Calculator{}

	// Чтение из консоли и выполнение операции
	// Обернуть в for {}, если нужен бесконечный цикл
	fmt.Println("Введите математическую операцию (например, 1 + 2):")

	// чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	result, err := c.Calculate(input)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Результат:", result)
	}
}
