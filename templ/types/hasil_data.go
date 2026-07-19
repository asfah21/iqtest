package types

// HasilPageData is the data required by the hasil (result) page template.
type HasilPageData struct {
	Nama             string
	RawScore         float64
	MaxPossible      float64
	Percentile       float64
	EstimatedIQ      *float64
	DomainScores     map[string]DomainScoreView
	AvgResponseMs    int
	IsReliable       bool
	ExecutiveSummary string
	Kekuatan         []string
	AreaPerhatian    []string
}

// DomainScoreView represents a single domain score for template rendering.
type DomainScoreView struct {
	Domain     string
	Percentage float64
	Label      string
}
