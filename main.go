package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"

	consumer "github.com/harlow/kinesis-consumer"
	"github.com/spf13/pflag"
)

var (
	sourceStream    string
	awsEndPoint     string
	numberOfRecords PositiveNumber
)

//////////////////////////////////////////////////////////
// Positive Number Type, to do pflag validation
// Note: I think there must be a better way of doing this
//////////////////////////////////////////////////////////
type PositiveNumber struct {
	num int
}

func (n *PositiveNumber) String() string {
	return strconv.Itoa(n.num)
}

func (n *PositiveNumber) Set(nstring string) error {
	i, err := strconv.Atoi(nstring)
	if err != nil {
		return err
	}

	if i < 1 {
		return fmt.Errorf("%s needs to be greater than 0 or not set at all", nstring)
	} else {
		n.num = i
		return nil
	}
}

func (n *PositiveNumber) Type() string {
	return "int"
}

//////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////

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

	// fmt.Printf("Reading %d records from stream %s using endpoint %s\n", numberOfRecords.num, sourceStream, awsEndPoint)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: createConfig(),
	}))

	kin := kinesis.New(sess)

	c, err := consumer.New(sourceStream, consumer.WithClient(kin))
	if err != nil {
		log.Fatalf("consumer error: %v", err)
	}

	recordsRead := 0
	continueScanning := consumer.ScanStatus{
		StopScan:       false,
		SkipCheckpoint: false,
	}
	stopScanning := consumer.ScanStatus{
		StopScan:       true,
		SkipCheckpoint: false,
	}
	err = c.Scan(context.TODO(), func(r *consumer.Record) consumer.ScanStatus {

		if recordsRead >= numberOfRecords.num {
			return stopScanning
		}

		recordsRead++
		fmt.Println(string(r.Data))

		if recordsRead >= numberOfRecords.num {
			return stopScanning
		} else {
			return continueScanning
		}
	})

	if err != nil {
		log.Fatalf("scan error: %v", err)
	}
}

func init() {
	pflag.StringVarP(&sourceStream, "source", "s", "", "Source Stream")
	pflag.StringVarP(&awsEndPoint, "end-point", "e", "", "AWS End Point")
	pflag.VarP(&numberOfRecords, "number", "n", "Number of records to read")
}
