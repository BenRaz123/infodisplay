package main

import (
	"fmt"
	"log"
)

type parseError struct {
	l uint
	r errReason
	d string
	t string
}

type errReason int

const (
	invalidDirective errReason = iota
	onlyDirectives
	invalidNum
	redundantDirective
	noArgs
	takesArgs
	uknownDirective
	illegalUseOfDT
	illegalUseofBold
	//TODO:fix this, this limitation sucks
	cannotMixBulletsAndLines
)

func (e parseError) Error() string {
	ret := fmt.Sprintf("line %d: ", e.l)
	switch e.r {
	case onlyDirectives:
		ret += fmt.Sprintf(`non-directive in a directive-only block: %q`, e.t)
	case invalidNum:
		ret += fmt.Sprintf(`invalid number: %q`, e.t)
	case uknownDirective:
		ret += fmt.Sprintf(`uknown directive: %q`, e.d)
	case invalidDirective:
		ret += fmt.Sprintf(`directive !%s is not allowed in this context`, e.d)
	case illegalUseOfDT:
		ret += fmt.Sprintf(`use of DateTime wrapper within directive %q that does not allow it: %q`, e.d, e.t)
	case illegalUseofBold:
		ret += fmt.Sprintf(`use of bold text formatting within directive %q that does not allow it: %q`, e.d, e.t)
	case noArgs:
		ret += fmt.Sprintf(`directive !%s takes no arguments: %q`, e.d, e.t)
	case takesArgs:
		ret += fmt.Sprintf(`line: directive !%s takes arguments but none were given: %q`, e.d, e.t)
	case redundantDirective:
		ret += fmt.Sprintf(`redundant directive: %q`, e.d)
	case cannotMixBulletsAndLines:
		ret += fmt.Sprintf("cannot mix bullets and lines: %q", e.t)
	default:
		log.Fatal("invalid enum variant: %d", e.r)
	}
	return ret
}
