package dot

type EdgeArrowHead string
type EdgeType string

const (
	PortJoiner = ":"
	Joiner     = "_"

	EdgeArrowHeadNormal  EdgeArrowHead = "normal"
	EdgeArrowHeadONormal EdgeArrowHead = "onormal"
	EdgeArrowHeadDiamond EdgeArrowHead = "diamond"
	EdgeArrowHeadNone    EdgeArrowHead = "none"

	EdgeTypeSolid EdgeType = "solid"
	EdgeTypeDot   EdgeType = "dotted"
	EdgeTypeDash  EdgeType = "dashed"
)
