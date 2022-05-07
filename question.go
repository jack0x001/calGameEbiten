package main

import (
	"fmt"
	"math/rand"
)

type Question struct {
	NumberA  int
	Operator string
	NumberB  int
}

func (q *Question) Answer() int {
	switch q.Operator {
	case "+":
		return q.NumberA + q.NumberB
	case "-":
		return q.NumberA - q.NumberB
		//case "*":
		//	return q.NumberA * q.NumberB
		//case "/":
		//	return q.NumberA / q.NumberB
	}
	return 0
}

func (q *Question) String() string {
	return fmt.Sprintf("%d %s %d = ?", q.NumberA, q.Operator, q.NumberB)
}

func NewQuestion() *Question {

	var q *Question
	for true {
		a := rand.Intn(maxAnswer) + 1
		b := rand.Intn(maxAnswer) + 1

		//if a <b ,swap a and b
		if a < b {
			a, b = b, a
		}
		//rand number between 0 and 1,if 1, then +, else -
		r := rand.Intn(2)
		if r == 1 {
			q = &Question{a, "+", b}
		} else {
			q = &Question{a, "-", b}
		}

		if q.Answer() <= maxAnswer {
			break
		}
	}

	return q
}
