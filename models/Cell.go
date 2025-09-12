package models

// on memory
type Cell struct {
	ID       string
	PlayerID string
	Rank     int
	Progress int
	X        float64
	Y        float64
	// CellType string
	// Power    int
	Active bool
}
