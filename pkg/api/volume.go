package api

// Volume struct
type Volume struct {
	ID   string
	Name string
	// size in gigabytes
	Size         float64
	Region       string
	// needed for scaleway
	Organization string
}
