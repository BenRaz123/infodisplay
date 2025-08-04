package main

import "fmt"

type opt[T any] struct {
	Val T
	Has bool
}

func wrap[T any](v T) opt[T] {
	return opt[T]{Val: v, Has: true}
}

func none[T any]() opt[T] {
	return opt[T]{Has: false}
}

func (o opt[T]) String() string {
	switch {
	case o.Has:
		return fmt.Sprintf("%s", o.Val)
	default:
		return fmt.Sprintf("<nil>")
	}
}

func (o opt[T]) use(f func(T)) {
	if o.Has {
		f(o.Val)
	}
}

func (o opt[T]) Or(v T) T {
	if o.Has {
		return o.Val
	} else {
		return v
	}
}
