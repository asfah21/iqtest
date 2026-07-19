package types

// HasilPageData is the data required by the hasil (result) page template.
type HasilPageData struct {
	Nama                string
	Narsisme            int
	Machiavellian       int
	Psikopati           int
	ExecutiveSummary    string
	RelationshipProfile string
	Kekuatan            []string
	AreaPerhatian       []string
	RelationshipInsight string
	CompatibilityNotes  string
	ReflectionQuestions []string
}
