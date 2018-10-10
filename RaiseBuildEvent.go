package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"os"
	"time"
)

func main() {
	opts := make(map[string]*string)

	opts["status"] = flag.String("BuildStatus", "", "Status")
	opts["version"] = flag.String("BuildVersion", "", "Version")
	opts["buildDate"] = flag.String("BuildDate", "", "Optional")
	opts["commit"] = flag.String("Commit", "", "Commit")
	opts["name"] = flag.String("BuildName", "", "Name")
	opts["comment"] = flag.String("Comment", "", "Optional")
	opts["detailType"] = flag.String("DetailType", "build.notification", "DetailType (Optional)")
	opts["source"] = flag.String("Source", "", "Source")

	verbose := flag.Bool("v", false, "Verbose")

	flag.Parse()

	options := make(map[string]string)

	// Convert pointers map to string map to save on my typing
	for k, v := range opts {
		options[k] = *v
	}

	// Required...
	validators := []string{"status", "version", "commit", "name", "source"}
	for _, validator := range validators {
		if options[validator] == "" {
			fmt.Println("Not all options provided")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if options["buildDate"] == "" {
		options["buildDate"] = time.Now().String()
	}

	// Check for status = 0 (bad) or 1 (good)
	switch options["status"] {
	case "0":
		options["status"] = "Failed"
	case "1":
		options["status"] = "Succeeded"
	}

	if *verbose {
		for k, v := range options {
			fmt.Printf("%s:[%s]\n", k, v)
		}
	}

	detail := fmt.Sprintf("{ \"BuildStatus\": \"%s\", \"BuildVersion\": \"%s\", \"BuildDate\": \"%s\", \"commit\": \"%s\", \"BuildName\": \"%s\", \"Comment\": \"%s\"}",
		options["status"], options["version"], options["buildDate"], options["commit"], options["name"], options["comment"])

	if *verbose {
		fmt.Println(detail)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create the cloudwatch events client
	svc := cloudwatchevents.New(sess)

	result, err := svc.PutEvents(&cloudwatchevents.PutEventsInput{
		Entries: []*cloudwatchevents.PutEventsRequestEntry{
			&cloudwatchevents.PutEventsRequestEntry{
				Detail:     aws.String(detail),
				DetailType: aws.String(options["detailType"]),
				Source:     aws.String(options["source"]),
				Resources: []*string{
					aws.String(fmt.Sprintf("Build:%s", options["version"])),
					aws.String(options["version"]),
					aws.String(options["name"]),
				},
			},
		},
	})

	if err != nil {
		fmt.Println("Error putting event: ", err)
		os.Exit(2)
	}

	if *verbose {
		fmt.Println("Ingested events:", result.Entries)
	}
}
