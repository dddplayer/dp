package valueobject

type IntimacyElement struct {
	Key int
	Val float64
}
type IntimacyElements []*IntimacyElement

func (ele IntimacyElements) Len() int           { return len(ele) }
func (ele IntimacyElements) Swap(i, j int)      { ele[i], ele[j] = ele[j], ele[i] }
func (ele IntimacyElements) Less(i, j int) bool { return ele[i].Val >= ele[j].Val }
