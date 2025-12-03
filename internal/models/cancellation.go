package models

// CancellationInstruction represents instructions for cancelling a subscription for a specific service
type CancellationInstruction struct {
	ID           int    `db:"id" json:"id"`
	ServiceName  string `db:"service_name" json:"service_name"`
	Instructions string `db:"instructions" json:"instructions"`
	Note         string `db:"note" json:"note"`
	URL          string `db:"url" json:"url"`
	Category     string `db:"category" json:"category"`
}
