package entity

func lineRow(fields []DotAttribute, col int) *Row {
	r := blankRow()

	fno := len(fields)
	if fno == 0 {
		return r
	}
	spaceBetweenFields := fno - 1
	usableBlankSpace := col - spaceBetweenFields - rowStartEmptyBlockNum - rowEndEmptyBlockNum
	colSpan := usableBlankSpace / fno

	for i, f := range fields {
		d := column(f.Name(), f.Port(), f.Color())
		d.ColSpan = colSpan

		r.Data = append(r.Data, d)
		if i+1 < fno {
			r.Data = append(r.Data, blankColumn())
		}
	}

	for i := 0; i < rowEndEmptyBlockNum; i++ {
		r.Data = append(r.Data, blankColumn())
	}

	return r
}

func countRow(num, rowCapacity int) int {
	row := num / rowCapacity
	if num%rowCapacity > 0 {
		row += 1
	}
	return row
}

func rowFirstPadding(source []DotAttribute, row, column int) [][]DotAttribute {
	completionSrc := make([]DotAttribute, len(source))
	copy(completionSrc, source)

	blankNum := row*column - len(source)
	for i := 0; i < blankNum; i++ {
		completionSrc = append(completionSrc, nil)
	}
	target := make([][]DotAttribute, row)
	for i := 0; i < row; i++ {
		target[i] = completionSrc[i*column : (i+1)*column]
	}

	return target
}

func columnFirstPadding(source []DotAttribute, row, column int) [][]DotAttribute {
	var completionSrc []DotAttribute
	sLen := len(source)

	sMaxCol := sLen / row
	if sLen%row > 0 {
		sMaxCol += 1
	}

	if sMaxCol <= column {
		for i := 0; i < row; i++ {
			for j := 0; j < sMaxCol; j++ {
				idx := i*sMaxCol + j
				if idx < sLen {
					completionSrc = append(completionSrc, source[idx])
				} else {
					completionSrc = append(completionSrc, nil)
				}
			}
			for z := 0; z < column-sMaxCol; z++ {
				completionSrc = append(completionSrc, nil)
			}
		}
	}

	target := make([][]DotAttribute, row)
	for i := 0; i < row; i++ {
		target[i] = completionSrc[i*column : (i+1)*column]
	}

	return target
}

func leftRightRow(n *Node, e DotElement, maxLeft, maxRight int) {
	as := e.Attributes()
	attrsLen := len(as)

	if attrsLen == 2 {
		commands := as[0].([]DotAttribute)
		attrs := as[1].([]DotAttribute)

		mLen := len(commands)
		aLen := len(attrs)
		mRowNum := countRow(mLen, maxLeft)
		aRowNum := countRow(aLen, maxRight)

		var mRows, aRows [][]DotAttribute
		var row int
		if mRowNum >= aRowNum && mRowNum > 0 {
			row = mRowNum
			mRows = rowFirstPadding(commands, row, maxLeft)
			aRows = columnFirstPadding(attrs, row, maxRight)
		} else if aRowNum >= mRowNum && aRowNum > 0 {
			row = aRowNum
			aRows = rowFirstPadding(attrs, row, maxRight)
			mRows = columnFirstPadding(commands, row, maxLeft)
		}

		for rn := 0; rn < row; rn++ {
			ms := mRows[rn]
			as := aRows[rn]
			r := buildNodeRow(ms, as, rn == 0, e, maxRight, row)

			n.Rows = append(n.Rows, r)
		}
		n.Rows = append(n.Rows, blankRow())

	} else if attrsLen == 1 {
		methods := as[0].([]DotAttribute)
		mLen := len(methods)

		mRowNum := countRow(mLen, maxLeft)
		mRows := rowFirstPadding(methods, mRowNum, maxLeft)
		for rn := 0; rn < mRowNum; rn++ {
			ms := mRows[rn]
			r := buildNodeRow(ms, nil, rn == 0, e, maxRight, mRowNum)

			n.Rows = append(n.Rows, r)
		}
		n.Rows = append(n.Rows, blankRow())
	}
}

func buildNodeRow(ms, as []DotAttribute, isFirstRow bool, obj DotElement, right, row int) *Row {
	r := blankRow()

	lenMs := len(ms)
	reverseMs := make([]DotAttribute, lenMs)
	for i := 0; i < lenMs; i++ {
		reverseMs[i] = ms[lenMs-1-i]
	}

	for _, m := range reverseMs {
		var d *Data
		if m != nil {
			d = column(m.Name(), m.Port(), m.Color())
		} else {
			d = blankColumn()
		}
		r.Data = append(r.Data, d)
	}

	for i := 0; i < emptyBlockBeforeDomainObjNum; i++ {
		r.Data = append(r.Data, blankColumn())
	}

	if isFirstRow {
		d := column(obj.Name(), obj.Port(), obj.Color())
		d.RowSpan = row
		r.Data = append(r.Data, d)
	}
	if as != nil {
		for _, a := range as {
			var d *Data
			if a != nil {
				d = column(a.Name(), a.Port(), a.Color())
			} else {
				d = blankColumn()
			}
			r.Data = append(r.Data, d)
		}
	} else {
		for j := 0; j < right; j++ {
			d := blankColumn()
			r.Data = append(r.Data, d)
		}
	}
	for i := 0; i < rowEndEmptyBlockNum; i++ {
		r.Data = append(r.Data, blankColumn())
	}

	return r
}
