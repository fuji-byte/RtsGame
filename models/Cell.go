package models

type Cell struct {
	ID       string
	PlayerID string
	Rank     int
	Cell     int
	X        float32
	Y        float32
	Power    int //Rankが高くなれば、Power(生産量が上がる)
}
