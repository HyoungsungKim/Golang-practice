package main

import (
	"fmt"
	"sort"
)

var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compiler": {"data structures", "formal languages", "computer organization"},

	"data structures":       {"discrete math"},
	"database":              {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"network":               {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)

	var visitAll func(items []string)
	visitAll = func(items []string) {
		fmt.Println(seen)
		for _, item := range items {
			fmt.Println(item)
			fmt.Println(seen[item])
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	//Sort strings
	sort.Strings(keys)
	fmt.Println("###############")
	fmt.Println(keys)
	fmt.Println("###############")

	visitAll(keys)
	return order
}
