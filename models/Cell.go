package models

// on memory
type Cell struct {
	ID       string
	PlayerID string
	CellType string
	Rank     int
	Hp       int
	X        float64
	Y        float64
	Power    int //Rankが高くなれば、Power(生産量が上がる)
}
