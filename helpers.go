package main

import ()

// function that turns boolean into either 1 or -1
func bool_to_int(b bool) int {
	if b {
		return 1
	} else {
		return -1
	}
}