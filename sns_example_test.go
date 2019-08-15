package awsoc_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"go.krak3n.codes/awsoc"
	"go.krak3n.codes/awsoc/propagation/b3"
	"go.opencensus.io/trace"
)

func ExampleSNS_PublishWithContext() {
	cfg := &aws.Config{
		Region: aws.String("eu-west-1"),
	}

	if v := os.Getenv("AWS_DEFAULT_REGION"); v != "" {
		cfg.Region = aws.String(v)
	}

	if v := os.Getenv("AWS_SNS_ENDPOINT"); v != "" {
		cfg.Endpoint = aws.String(v)
	}

	session, err := session.NewSession(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, span := trace.StartSpan(context.Background(), "sns/test")
	defer span.End()

	in := &sns.PublishInput{
		TopicArn: aws.String("arn:aws:sns:eu-west-1:123456789012:Foo"),
		Message:  aws.String(`{"foo":"bar"}`),
	}

	c := awsoc.NewSNS(sns.New(session))
	if _, err := c.PublishWithContext(ctx, in); err != nil {
		log.Fatal(err)
	}

	fmt.Println("TraceID:", *in.MessageAttributes[b3.TraceIDKey].StringValue)
	fmt.Println("SpanID:", *in.MessageAttributes[b3.SpanIDKey].StringValue)
	fmt.Println("Span Sampled:", *in.MessageAttributes[b3.SpanSampledKey].StringValue)
	fmt.Println("Trace Topic Name:", *in.MessageAttributes[awsoc.TraceTopicName].StringValue)

	// Output:
	// TraceID: 616263646566676869676b6c6d6e6f71
	// SpanID: 6162636465666768
	// Span Sampled: 0
	// Trace Topic Name: Foo
}
