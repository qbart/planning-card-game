package main

func intFind(vals []int, val int) int {
	for i := 0; i < len(vals); i++ {
		if vals[i] == val {
			return i
		}
	}
	return -1
}
