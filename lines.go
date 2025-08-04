package main

import (
	"fmt"
	"strings"
)

type lines []line

type line struct {
	string
	empty bool
	index uint
}

func (l line) String() string {
	return fmt.Sprintf("%d: %s", l.index, l.string) + func() string {
		if l.empty {
			return "(empty)"
		} else {
			return ""
		}
	}()
}

func linesFromString(s string) lines {
	l := lines{}
	for i, x := range strings.Split(s, "\n") {
		l = append(l, line{empty: x == "", index: uint(i + 1), string: x})
	}
	return l
}

func (l lines) toBlocks() []lines {
	blocks := []lines{lines{}}
	blockIndex := 0
	for _, line := range l {
		if line.empty {
			blockIndex += 1
			blocks = append(blocks, lines{})
			continue
		}
		blocks[blockIndex] = append(blocks[blockIndex], line)
	}
	return blocks
}

func (l lines) String() string {
	var ret string
	for _, v := range l {
		ret += v.String() + "\n"
	}
	return ret
}
