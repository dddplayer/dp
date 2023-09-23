package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
)

func (n *Node) Build(els []arch.Element) error {
	firstRow := blankRow()
	firstRow.Data[0].Port = firstRowPort
	n.Table.Rows = append(n.Table.Rows, firstRow)

	maxLeft, maxRight := ElementsIndicators(els)
	cols := rowStartEmptyBlockNum +
		maxLeft + emptyBlockBeforeDomainObjNum + domainObjBlockNum +
		maxRight + rowEndEmptyBlockNum

	for _, e := range els {
		if isLeftRightStructure(e) {
			leftRightRow(n, e, maxLeft, maxRight)
		} else {
			var attrs []arch.Node
			for _, a := range e.Children() {
				for _, n := range a {
					attrs = append(attrs, n)
				}
			}

			maxAttrNum := (cols - rowStartEmptyBlockNum - rowEndEmptyBlockNum) / 2
			chunk := chunkSlice(attrs, maxAttrNum)
			for _, c := range chunk {
				r := lineRow(c, cols)
				n.Table.Rows = append(n.Table.Rows, r)
			}

			n.Table.Rows = append(n.Table.Rows, blankRow())
		}
	}

	nr := nameRow(n.Name, cols)
	n.Table.Rows = append(n.Table.Rows, nr)

	return nil
}

func chunkSlice(slice []arch.Node, chunkSize int) [][]arch.Node {
	var chunks [][]arch.Node
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func isLeftRightStructure(e arch.Element) bool {
	children := e.Children()
	if len(children) >= 2 {
		return true
	}
	return false
}
