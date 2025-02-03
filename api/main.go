package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	DefaultPort           = 5000
	DefaultFileCreateMode = 0700
)

var (
	DefaultResultsDataDir = getEnvDefault("RESULTS_DATA_DIR", "../data/results")
)

type ArbitraryJSON map[string]interface{}

func main() {
	// Fiber instance
	app := fiber.New()

	// Routes
	app.Get("/result", getResults)
	app.Post("/result", postResult)

	// Start server
	listenString := fmt.Sprintf(":%d", DefaultPort)
	log.Fatal(app.Listen(listenString))
}

// getResults returns all of the results that have been previously stored using postResult
func getResults(c *fiber.Ctx) error {
	files, err := os.ReadDir(DefaultResultsDataDir)
	if err != nil {
		return errors.Wrap(err, "unable to read results")
	}

	// We do not know the length yet here because files might contain invalid JSON (and therefore be ignored)
	results := make([]ArbitraryJSON, 0)

	for _, file := range files {
		filename := fmt.Sprintf("%s/%s", DefaultResultsDataDir, file.Name())
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Printf(fmt.Sprintf("skipping result from %s because of error: %v", file.Name(), err))
			continue
		}

		result, err := parseArbitraryJSON(content)
		if err != nil {
			log.Printf(fmt.Sprintf("skipping result from %s, file does not contain valid JSON: %v", file.Name(), err))
			continue
		}

		results = append(results, result)
	}

	return c.JSON(results)
}

// getResult stores a new result
func postResult(c *fiber.Ctx) error {
	content := c.Body()

	if err := testIfValidJSON(content); err != nil {
		return errors.Wrap(err, "unable to parse body as JSON")
	}

	// Make sure the data dir exists
	makeSurePathExists(DefaultResultsDataDir)

	// Write body to file
	filename := filepath.Join(DefaultResultsDataDir, getRandomFilename())
	err := os.WriteFile(filename, content, DefaultFileCreateMode)
	if err != nil {
		return errors.Wrap(err, "unable to write result file")
	}

	return c.JSON("OK")
}

// getEnvDefault return the value of the environment variable specified by name, or the defaultValue if not set
func getEnvDefault(name string, defaultValue string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return defaultValue
}

func parseArbitraryJSON(value []byte) (ArbitraryJSON, error) {
	var result ArbitraryJSON

	err := json.Unmarshal(value, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func testIfValidJSON(value []byte) error {
	contentJSON := struct{}{}
	return json.Unmarshal(value, &contentJSON)
}

func makeSurePathExists(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, DefaultFileCreateMode)

		if err != nil {
			return err
		}
	}

	return nil
}

func getRandomFilename() string {
	id := uuid.New()
	return id.String()
}
