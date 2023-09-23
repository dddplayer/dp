package entity

import "github.com/dddplayer/dp/internal/domain/arch"

func ElementsIndicators(els []arch.Element) (maxLeft int, maxRight int) {
	maxLeft, maxRight = 0, 0
	left, right := 0, 0

	for _, e := range els {
		for i, nodes := range e.Children() {
			if i == 0 {
				left = len(nodes)
				if left > maxLeft {
					maxLeft = left
				}
			} else if i == 1 {
				right = len(nodes)
				if right > maxRight {
					maxRight = right
				}
			} else {
				break
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
