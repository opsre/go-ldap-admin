package model

// Role permission rule
type RoleCasbin struct {
	Keyword string `json:"keyword"` // Role keyword
	Path    string `json:"path"`    // Request path
	Method  string `json:"method"`  // Request method
}
