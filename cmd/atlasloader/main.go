package main

import (
	"io"
	"log/slog"
	"os"

	"monitoring-energy-service/internal/domain/entities"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&entities.ExampleEntity{},
		&entities.EnergyPlants{},
		// Add more entities here as needed
	)
	if err != nil {
		slog.Error("Failed to load gorm schema", "err", err.Error())
		os.Exit(1)
	}
	_, err = io.WriteString(os.Stdout, stmts)
	if err != nil {
		slog.Error("Failed to write gorm schema", "err", err.Error())
		os.Exit(1)
	}
}
