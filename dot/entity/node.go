package entity

import "errors"

func (n *Node) Build(els DotElements) error {
	firstRow := blankRow()
	firstRow.Data[0].Port = firstRowPort
	n.Rows = append(n.Rows, firstRow)

	maxLeft, maxRight := els.indicators()
	cols := rowStartEmptyBlockNum +
		maxLeft + emptyBlockBeforeDomainObjNum + domainObjBlockNum +
		maxRight + rowEndEmptyBlockNum

	if ValidateElements(els) {
		for _, e := range els {
			if isLeftRightStructure(e) {
				leftRightRow(n, e, maxLeft, maxRight)
			} else {
				var attrs []DotAttribute
				for _, a := range e.Attributes() {
					attrs = append(attrs, a.(DotAttribute))
				}

				maxAttrNum := (cols - rowStartEmptyBlockNum - rowEndEmptyBlockNum) / 2
				chunk := chunkSlice(attrs, maxAttrNum)
				for _, c := range chunk {
					r := lineRow(c, cols)
					n.Rows = append(n.Rows, r)
				}

				n.Rows = append(n.Rows, blankRow())
			}
		}

		nr := nameRow(n.Name, cols)
		n.Rows = append(n.Rows, nr)
		return nil
	}
	return errors.New("invalid elements")
}

func chunkSlice(slice []DotAttribute, chunkSize int) [][]DotAttribute {
	var chunks [][]DotAttribute
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func ValidateElements(els DotElements) bool {
	for _, e := range els {
		if isLeftRightStructure(e) {
			for _, a := range e.Attributes() {
				if _, ok := a.([]DotAttribute); !ok {
					return false
				}
			}
		} else {
			for _, a := range e.Attributes() {
				if _, ok := a.(DotAttribute); !ok {
					return false
				}
			}
		}
	}
	return true
}

func isLeftRightStructure(e DotElement) bool {
	attrs := e.Attributes()
	if len(attrs) >= 1 {
		if _, ok := attrs[0].([]DotAttribute); ok {
			return true
		}
	}
	return false
}
