package entity

import (
	"fmt"
)

const (
	DotPortJoiner = ":"
	DotJoiner     = "_"
)

type DotGraph interface {
	Name() string
	Nodes() []DotNode
	Edges() []DotEdge
}

type DotNode interface {
	Name() string
	Elements() []DotElement
}

type DotElement interface {
	Name() string
	Port() string
	Color() string
	Attributes() []any // could be DotAttribute, []DotAttribute
}

type DotAttribute interface {
	Name() string
	Port() string
	Color() string
}

type DotEdge interface {
	From() string
	To() string
}

func DotAttrPort(edgeName, attrName string) string {
	return fmt.Sprintf("%s%s%s", edgeName, DotPortJoiner, attrName)
}

type DotElements []DotElement

func (des DotElements) First() DotElement {
	if len(des) > 0 {
		return des[0]
	}
	return nil
}

func (des DotElements) indicators() (maxLeft int, maxRight int) {
	maxLeft, maxRight = 0, 0
	left, right := 0, 0

	for _, e := range des {
		for i, a := range e.Attributes() {
			if as, ok := a.([]DotAttribute); ok {
				if i == 0 {
					left = len(as)
					if left > maxLeft {
						maxLeft = left
					}
				} else if i == 1 {
					right = len(as)
					if right > maxRight {
						maxRight = right
					}
				} else {
					break
				}

			}
		}
	}

	if maxLeft >= maxNum {
		maxLeft = maxNum
	} else if maxLeft <= minNum {
		maxLeft = minNum
	}
	if maxRight >= maxNum {
		maxRight = maxNum
	} else if maxRight <= minNum {
		maxRight = minNum
	}

	return
}
