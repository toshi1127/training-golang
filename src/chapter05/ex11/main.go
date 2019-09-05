package main

import (
	"fmt"
)

var prereqs = map[string]map[string]bool{
	"algorithms": {"data structures": true},
	"calculus":   {"linear algebra": true},

	"compilers": {
		"data structures":       true,
		"formal languages":      true,
		"computer organization": true,
	},

	"data structures":       {"discrete math": true},
	"databases":             {"data structures": true},
	"discrete math":         {"intro to programming": true},
	"formal languages":      {"discrete math": true},
	"networks":              {"operating systems": true},
	"operating systems":     {"data structures": true, "computer organization": true},
	"programming languages": {"data structures": true, "computer organization": true},
	"linear algebra":        {"calculus": true},
}

func main() {
	courseOrder := TopoSort(prereqs)
	CheckCircularReference(courseOrder)
	for i, course := range courseOrder {
		fmt.Printf("%d:\t%s\n", i, course)
	}
}

func TopoSort(m map[string]map[string]bool) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(map[string]bool)

	visitAll = func(items map[string]bool) {
		for item, _ := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}

	keys := make(map[string]bool)
	for k := range m {
		keys[k] = true
	}

	visitAll(keys)
	return order
}

func CheckCircularReference(sortedCourses []string) {
	attendanceOrder := make(map[string]int)

	for order, course := range sortedCourses {
		attendanceOrder[course] = order
	}

	for course, order := range attendanceOrder {
		for prereq, _ := range prereqs[course] {
			if order < attendanceOrder[prereq] {
				fmt.Printf("circular reference found! course: %s, prereq: %s\n", course, prereq)
			}
		}
	}
}