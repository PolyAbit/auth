package main

import (
	"github.com/PolyAbit/auth/internal/config"
)

func main() {
	cfg := config.MustLoad()

	_ = cfg
}
