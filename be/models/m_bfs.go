package models


type Node struct {
	Name    string
	IsBegin bool
}

type VisitedVal struct {
	Path    []string
	Depth   int
	IsBegin bool
}