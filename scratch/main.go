package main

import (
	"fmt"
	"os"
	"github.com/gocarina/gocsv"
)

type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

func importCSV(filename string, numRows int) ([]Asset, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var assets []Asset
	if err := gocsv.Unmarshal(file, &assets); err != nil {
		return nil, err
	}

	// Limit the number of rows if numRows is provided
	if numRows > 0 && numRows < len(assets) {
		assets = assets[:numRows]
	}

	return assets, nil
}

func generateAssets(numAssets int) []Asset {
	var assets []Asset
	for i := 1; i <= numAssets; i++ {
		asset := Asset{
			ID:             fmt.Sprintf("pc%d", i),
			Color:          "blue", // Example color
			Size:           5,      // Example size
			Owner:          "Owner", // Example owner
			AppraisedValue: 3000,   // Example appraised value
		}
		assets = append(assets, asset)
	}
	return assets
}

func main() {
	
	// records, err := importCSV("data.csv", 5) // Read first 10 rows
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	records := generateAssets(100000)

	for _, record := range records {
		fmt.Printf("Name: %s, Age: %d, Email: %s\n", record.ID, record.Color, record.Owner)
	}
}
