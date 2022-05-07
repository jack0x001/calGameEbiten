package main

import (
	"math/rand"
	"time"
)

func floorMod(x, y int) int {
	return x - floorDiv(x, y)*y
}
func floorDiv(x, y int) int {
	d := x / y
	if d*y == x || x >= 0 {
		return d
	}
	return d - 1
}

// RandomInt 返回可以在负数和正数之间的随机数.
// 内置的rand.Intn()函数只能在0和正数之间返回随机数.
func randInt(lower, upper int) int {
	rand.Seed(time.Now().UnixNano())
	//rand.Intn() does not accept anything less than zero
	//so lower and upper should not be the same and upper should always be larger than lower.
	rng := upper - lower
	if rng < 0 {
		panic("upper must be larger than lower")
	}
	if rng == 0 {
		return lower
	}
	return rand.Intn(rng) + lower
}
