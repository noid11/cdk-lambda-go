package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-xray-sdk-go/xraylog"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/lambda"
	"golang.org/x/net/context/ctxhttp"

	"github.com/aws/aws-xray-sdk-go/xray"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var lambdaSvc = lambda.New(sess)
var ec2Svc = ec2.New(sess)

func init() {
	xray.Configure(xray.Config{
		ServiceVersion: "1.0.0",
	})

	xray.SetLogger(xraylog.NewDefaultLogger(os.Stdout, xraylog.LogLevelDebug))

	xray.AWS(lambdaSvc.Client)
	xray.AWS(ec2Svc.Client)
}

func callLambda(ctx context.Context) (string, error) {
	input := &lambda.GetAccountSettingsInput{}
	req, err := lambdaSvc.GetAccountSettingsWithContext(ctx, input)
	if err != nil {
		log.Print(err.Error())
	}
	output, _ := json.MarshalIndent(req, "", "  ")
	return string(output), err
}

func callEc2(ctx context.Context) (string, error) {
	req, err := ec2Svc.DescribeRegionsWithContext(ctx, nil)
	if err != nil {
		log.Print(err.Error())
	}
	output, _ := json.MarshalIndent(req, "", "  ")
	return string(output), err
}

func getExample(ctx context.Context) (int, error) {
	resp, err := ctxhttp.Get(ctx, xray.Client(nil), "https://www.example.com/")
	if err != nil {
		log.Print(err.Error())
	}
	return resp.StatusCode, err
}

func handleRequest(ctx context.Context, event events.SQSEvent) (string, error) {

	// event
	eventJSON, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("EVENT: %s", eventJSON)

	// environment variables
	log.Printf("REGION: %s", os.Getenv("AWS_REGION"))
	log.Println("ALL ENV VARS:")
	for _, element := range os.Environ() {
		log.Println(element)
	}

	// request context
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("REQUEST ID: %s", lc.AwsRequestID)

	// global variable
	log.Printf("FUNCTION NAME: %s", lambdacontext.FunctionName)

	// context method
	deadline, _ := ctx.Deadline()
	log.Printf("DEADLINE: %s", deadline)

	// AWS SDK call
	lambdaUsage, err := callLambda(ctx)
	if err != nil {
		return "ERROR", err
	}
	log.Printf("LAMBDA USAGE: %s", lambdaUsage)

	ec2Usage, err := callEc2(ctx)
	if err != nil {
		return "ERROR", err
	}
	log.Printf("EC2 USAGE: %s", ec2Usage)

	// http request
	httpReq, err := getExample(ctx)
	if err != nil {
		return "ERROR", err
	}
	log.Printf("HTTP REQUESTED STATUS CODE: %d", httpReq)

	_, subseg := xray.BeginSubsegment(ctx, "my-subsegment")
	subseg.Close(nil)

	return "Hello World!", nil
}

func main() {
	runtime.Start(handleRequest)
}
