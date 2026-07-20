package models

type Severity string

const (
	SeverityLow    Severity = "LOW"
	SeverityMedium Severity = "MEDIUM"
	SeverityHigh   Severity = "HIGH"
)

type Finding struct {
	RuleID         string
	Severity       Severity
	FilePath       string
	Message        string
	Recommendation string
}
