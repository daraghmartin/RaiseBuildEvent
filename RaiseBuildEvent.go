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
	status := flag.String("BuildStatus", "", "Status")
	version := flag.String("BuildVersion", "", "Version")
	buildDate := flag.String("BuildDate", "", "Optional")
	commit := flag.String("Commit", "", "Commit")
	name := flag.String("BuildName", "", "Name")
	comment := flag.String("Comment", "", "Optional")
	detailType := flag.String("DetailType", "build.notification", "DetailType (Optional)")
	source := flag.String("Source", "", "Source")
	verbose := flag.Bool("v", false, "Verbose")

	flag.Parse()

	options := make(map[string]string)

	// Required...
	options["status"] = *status
	options["version"] = *version
	options["commit"] = *commit
	options["name"] = *name
	options["source"] = *source

	// so far these are all required
	for _, v := range options {
		if v == "" {
			fmt.Println("Not all options provided")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	// Optional...
	options["buildDate"] = *buildDate
	if options["buildDate"] == "" {
		options["buildDate"] = time.Now().String()
	}
	options["comment"] = *comment
	options["detailType"] = *detailType

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
