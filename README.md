# Goroutines pool
[![Build Status](https://travis-ci.org/DukeLog/gopool.svg?branch=master)](https://travis-ci.org/DukeLog/gopool)

Managing gouroutines pool

### Example
```
    // Create iterable (slice, array)
	slice := [][]int{1,2,3,4,5}
	// Create goroutines pool
	p := New(5)
	// Run function f on each element of the slice
	p.Map(f, slice)
	// Get results of each function as a slice
	result := p.Join()
