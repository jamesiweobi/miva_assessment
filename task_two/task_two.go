package main

import (
	"fmt"
	"strconv"
)

func isPrime(num int) bool {
	if num < 2 {
		return false
	}
	if num == 2 {
		return true
	}
	if num%2 == 0 {
		return false
	}
	
	for i := 3; i*i <= num; i += 2 {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func isPalindrome(num int) bool {
	if num < 10 {
		return true
	}
	
	str := strconv.Itoa(num)
	left, right := 0, len(str)-1
	
	for left < right {
		if str[left] != str[right] {
			return false
		}
		left++
		right--
	}
	return true
}

func findPrimePalindromes(N int) int {
	results := make(chan int, N)
	
	go func() {
		defer close(results)
		count := 0
		num := 2
		
		for count < N {
			if isPrime(num) && isPalindrome(num) {
				results <- num
				count++
			}
			if num == 2 {
				num = 3
			} else {
				num += 2
			}
		}
	}()
	
	sum := 0
	for prime := range results {
		sum += prime
	}
	
	return sum
}

func completeTaskTwo() {
	var N int
	fmt.Print("Enter N (1-50): ")
	fmt.Scan(&N)
	
	if N < 1 || N > 50 {
		fmt.Println("N must be between 1 and 50")
		return
	}
	
	sum := findPrimePalindromes(N)
	fmt.Printf("Sum of first %d prime palindromic numbers: %d\n", N, sum)
}

func main() {
	fmt.Println("Starting Task Two")
completeTaskTwo()
fmt.Println("Completed Task Two")
}

