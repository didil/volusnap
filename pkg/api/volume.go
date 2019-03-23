package api

type volume struct {
	ID   string
	Name string
	// size in gigabytes
	Size   float64
	Region string
	// needed for scaleway
	Organization string
}
