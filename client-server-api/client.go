package main

import (
	"context"
)

func client() {
	ctx := context.Background()

	Do(ctx, "GET", "http://localhost:8080")
}
