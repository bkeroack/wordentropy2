package main

import (
	"crypto/rand"
	"log"
	"math/big"
)

func random_range(max int64) int64 {
	max_big := *big.NewInt(max)
	n, err := rand.Int(rand.Reader, &max_big)
	if err != nil {
		log.Fatalf("ERROR: cannot get random integer!\n")
	}
	return n.Int64()
}
