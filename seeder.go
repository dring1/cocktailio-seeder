package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dring1/cocktailio/models"
	"github.com/olivere/elastic"
)

// Seeder struct used to seed
type Seeder struct {
	esEndpoint string
	client     *elastic.Client
}

// SeedConfig struct define params for seeding
type SeedConfig struct {
	Index    string
	Mappings interface{}
	Data     []models.Cocktail
}

// NewSeeder create the es client
// test the connection ?
func NewSeeder(esEndpoint string) (*Seeder, error) {
	client, err := elastic.NewClient(elastic.SetURL(esEndpoint), elastic.SetSniff(false))
	return &Seeder{client: client, esEndpoint: esEndpoint}, err
}

// Ping the es endpoint
func (s *Seeder) Ping() error {
	ctx := context.Background()
	info, code, err := s.client.Ping(s.esEndpoint).Do(ctx)
	if err != nil {
		return err
	}
	ctx.Done()
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return nil
}

// Seed ...
// SeedConfig
//	Index
// 	Mappings
//	Data
func (s *Seeder) Seed(ctx context.Context, index string, mappings interface{}, data []models.Cocktail) error {
	exists, err := s.client.IndexExists("cocktailio").Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("creating new index cocktailio...")
		// Create a new index.
		createIndex, err := s.client.CreateIndex("cocktailio").BodyJson(mappings).Do(ctx)
		if err != nil {
			// Handle error
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			fmt.Println("create index not acknowledged")
			return err
		}
		fmt.Println("creating new index cocktailio...")
	}
	return s.Index(ctx, data)
}

// Nuke ... Will clear the indices of the elastic search endpoint
func (s *Seeder) Nuke(ctx context.Context) error {
	_, err := s.client.DeleteIndex("cocktailio").Do(ctx)
	return err
}

// Reseed ... Nuke and Seed
func (s *Seeder) Reseed() error {
	return nil
}

// Index ...
func (s *Seeder) Index(ctx context.Context, cocktails []models.Cocktail) error {
	for i, cocktail := range cocktails {
		put, err := s.client.Index().Index("cocktailio").Type("cocktail").Id(strconv.Itoa(i)).BodyJson(cocktail).Do(ctx)
		if err != nil {
			fmt.Printf("Error'd on cocktail %s", cocktail.Name)
			return err
		}
		fmt.Printf("Indexed cocktail %s to index %s, type %s\n", put.Id, put.Index, put.Type)
	}
	return nil
}
