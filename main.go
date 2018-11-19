package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"

	consumer "github.com/harlow/kinesis-consumer"
	"github.com/spf13/pflag"
)

var (
	sourceStream    string
	awsEndPoint     string
	numberOfRecords int
)

func createConfig() aws.Config {
	if awsEndPoint != "" {
		return aws.Config{
			Endpoint: aws.String(awsEndPoint),
		}
	} else {
		return aws.Config{}
	}
}

func main() {
	pflag.Parse()

	if pflag.NFlag() == 0 {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		pflag.PrintDefaults()
		os.Exit(1)

	}

	fmt.Printf("Reading from stream %s using endpoint %s\n", sourceStream, awsEndPoint)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: createConfig(),
	}))

	kin := kinesis.New(sess)

	c, err := consumer.New(sourceStream, consumer.WithClient(kin))
	if err != nil {
		log.Fatalf("consumer error: %v", err)
	}

	recordsRead := 0
	err = c.Scan(context.TODO(), func(r *consumer.Record) consumer.ScanStatus {
		fmt.Println(string(r.Data))
		recordsRead++
		var stopScanning bool = false
		if recordsRead == numberOfRecords {
			stopScanning = true
		}

		return consumer.ScanStatus{
			StopScan:       stopScanning,
			SkipCheckpoint: false,
		}
	})

	if err != nil {
		log.Fatalf("scan error: %v", err)
	}
}

func init() {
	pflag.StringVarP(&sourceStream, "source", "s", "", "Source Stream")
	pflag.StringVarP(&awsEndPoint, "end-point", "e", "", "AWS End Point")
	pflag.IntVarP(&numberOfRecords, "number", "n", -1, "Number of records to read")
}
