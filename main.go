package main

import (
	"cloud.google.com/go/logging"
	cplogging "commentparser/logging"
	"commentparser/models"
	"commentparser/server"
	"commentparser/services"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// entry point for the application, see readme.md for instructions
func main() {

	fmt.Printf("Starting Comment Parser %v\n\n", time.Now())
	// look for server mode
	if len(os.Args) < 2 {
		os.Stderr.WriteString("Two parameters required: package_name and (comma-seperated) search_terms or " +
			"'server' optionally followed by configuration path location")
	} else if os.Args[1] == "server" {
		ctx := context.Background()

		var config server.Configuration
		// this is the default dir to look at the configuration
		configFilePath := os.Getenv("HOME") + "/configuration/development.json"
		if len(os.Args) > 2 && len(os.Args[2]) > 0 {
			configFilePath = os.Args[2]
		}

		configFile, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(configFile, &config)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse the configuration file at %s", configFilePath)
			panic(err)
		}

		fmt.Printf("Will start server at %s", config.Address)

		credFilePath := config.GoogleCloudCredFile
		if strings.Index(credFilePath, "~/") == 0 {
			credFilePath = credFilePath[1:]
			credFilePath = os.Getenv("HOME") + credFilePath
		}
		client, err := logging.NewClient(
			ctx,
			config.GoogleCloudProjectID,
			option.WithCredentialsFile(
				credFilePath))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		loggerGC := client.Logger(config.LogName)

		loggingGC := cplogging.NewStackdriverLogger(loggerGC)

		loggingGC.Debug("Starting server")

		// measurement
		measurementGC := server.NewMeasurementStackdriver(loggerGC)

		server.CommentParserHttpServer(config, loggingGC, measurementGC)

		if err := client.Close(); err != nil {
			log.Fatalf("Failed to close client: %v", err)
		}
	} else {
		searchTerms := strings.Split(os.Args[2], ",")
		request := models.CommentParsingRequest{
			PackageName: os.Args[1],
			Tokens:      searchTerms,
		}
		res, err := services.ExtractRelevantComments(request, cplogging.NewConsoleLogging())

		for _, matches := range res.Matches {
			for _, match := range matches {
				fmt.Fprintf(
					os.Stdout,
					"%s:%v:\n%s\n",
					match.FileName,
					match.LineNumber,
					match.LineContent)
			}
		}

		if err != nil {
			panic(err)
		}
	}
	os.Exit(0)
}
