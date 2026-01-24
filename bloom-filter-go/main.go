package main

import (
	"fmt"
)

func main() {
	bf := NewBloomFilter(50, 3) 

	items := []string{
		"apple", "banana", "cherry", "date", "elderberry",
		"fig", "grape", "honeydew", "kiwi", "lemon",
		"mango", "nectarine", "orange", "papaya", "quince",
	}

	fmt.Println("Adding items to Bloom Filter...")
	for _, item := range items {
		bf.Add(item)
		fmt.Printf("Added: %s\n", item)
	}

	fmt.Println("\n--- Testing for false positives ---")

	fmt.Println("\nItems that WERE added:")
	fmt.Printf("apple: %v\n", bf.MightContain("apple"))
	fmt.Printf("banana: %v\n", bf.MightContain("banana"))

	fmt.Println("\nItems that were NOT added (false positives will show true):")
	testItems := []string{
		"strawberry", "blueberry", "raspberry", "watermelon",
		"peach", "pear", "plum", "apricot", "coconut", "dragonfruit",
	}

	falsePositives := 0
	for _, item := range testItems {
		result := bf.MightContain(item)
		fmt.Printf("%s: %v", item, result)
		if result {
			fmt.Print(" ← FALSE POSITIVE!")
			falsePositives++
		}
		fmt.Println()
	}

	fmt.Printf("\nFalse positive rate: %d/%d (%.1f%%)\n",
		falsePositives, len(testItems),
		float64(falsePositives)/float64(len(testItems))*100)
}
