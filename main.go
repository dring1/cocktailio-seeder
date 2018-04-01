package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dring1/cocktailio/models"
)

// 1 provision
// 2 nuke and provision

func main() {

	// parse command line args
	// seed
	seedCommand := flag.NewFlagSet("seed", flag.ExitOnError)
	esEndpointPtr := seedCommand.String("es-endpoint", "http://localhost:9200", "elasticsearch url")
	mappingsPathPtr := seedCommand.String("es-mappings-path", "mappings.json", "path to es mappings")
	seedDataPathPtr := seedCommand.String("seed-data-path", "seed.json", "path to es seed data")

	// nuke
	nukeCommand := flag.NewFlagSet("nuke", flag.ExitOnError)
	nukeEsEndpointPtr := nukeCommand.String("es-endpoint", "http://localhost:9200", "elasticsearch url")

	flag.Parse()

	switch os.Args[1] {
	case "seed":
		seedCommand.Parse(os.Args[2:])
		seeder, err := NewSeeder(*esEndpointPtr)
		if err != nil {
			panic(err)
		}
		err = seed(seeder, *esEndpointPtr, *mappingsPathPtr, *seedDataPathPtr)
		if err != nil {
			panic(err)
		}
	case "nuke":
		fmt.Println("Nuking index...")
		nukeCommand.Parse(os.Args[2:])
		seeder, err := NewSeeder(*nukeEsEndpointPtr)
		if err != nil {
			panic(err)
		}
		err = seeder.Nuke(context.Background())
		if err != nil {
			panic(err)
		}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	// create
}

func seed(seeder *Seeder, esURL string, mappingsPath string, seedDataPath string) error {
	ctx := context.Background()
	mappings, err := ReadMappings(mappingsPath)
	if err != nil {
		return err
	}
	seedData, err := ReadSeedData(seedDataPath)
	if err != nil {
		return err
	}
	err = seeder.Seed(ctx, "cocktailio", mappings, seedData)
	if err != nil {
		return err
	}
	return nil
}

// ReadMappings ...
func ReadMappings(path string) (interface{}, error) {
	plan, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data interface{}
	err = json.Unmarshal(plan, &data)
	return data, err
}

// ReadSeedData ...
func ReadSeedData(path string) ([]models.Cocktail, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cocktails []models.Cocktail
	err = json.Unmarshal(file, &cocktails)
	return cocktails, err
}
