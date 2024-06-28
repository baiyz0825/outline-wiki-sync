package xrandonm

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateRandomColor generates a random hex color string
func GenerateRandomColor() string {
	// Create a new random generator with a seed based on the current time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("#%06x", r.Intn(0xFFFFFF))
}
