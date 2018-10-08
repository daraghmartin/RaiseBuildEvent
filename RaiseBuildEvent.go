package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	// 	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	// 	"github.com/aws/aws-sdk-go/service/ssm"
	"fmt"
	"os"
	"time"
	// 	"strconv"
	// 	"strings"
)

func main() {
	numArgs := 7

	if len(os.Args) != (numArgs + 1) {
		gotArgs := len(os.Args) - 1
		fmt.Printf("Error need %d args - got %d\n", numArgs, gotArgs)
		os.Exit(1)
	}

	name := os.Args[1]
	version := os.Args[2]
	commit := os.Args[3]
	status := os.Args[4]
	detailType := os.Args[5]
	source := os.Args[6]
	comment := os.Args[7]

	now := time.Now()

	detail := fmt.Sprintf("{ \"BuildStatus\": \"%s\", \"BuildVersion\": \"%s\", \"BuildDate\": \"%s\", \"commit\": \"%s\", \"BuildName\": \"%s\", \"Comment\": \"%s\"}", status, version, now, commit, name, comment)

	fmt.Println(detail)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create the cloudwatch events client
	svc := cloudwatchevents.New(sess)

	result, err := svc.PutEvents(&cloudwatchevents.PutEventsInput{
		Entries: []*cloudwatchevents.PutEventsRequestEntry{
			&cloudwatchevents.PutEventsRequestEntry{
				Detail:     aws.String(detail),
				DetailType: aws.String(detailType),
				Source:     aws.String(source),
				Resources: []*string{
					aws.String(fmt.Sprintf("Build:%s", version)),
					aws.String(version),
					aws.String(name),
				},
			},
		},
	})

	if err != nil {
		fmt.Println("Error putting event: ", err)
		os.Exit(2)
	}

	fmt.Println("Ingested events:", result.Entries)
}
