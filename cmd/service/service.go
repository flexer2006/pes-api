package main

import (
	"context"
	"fmt"
	"os"

	"github.com/flexer2006/case-person-enrichment-go/internal/service/app"
)

const envs = "./deploy/.env"

func main() {
	if err := app.Run(context.Background(), envs); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "service startup failed:", err)
		os.Exit(1)
	}
}
