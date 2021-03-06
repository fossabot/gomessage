package main

import (
	"context"
	"crypto/tls"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/jung-kurt/gofpdf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

var (
	influxURI         = flag.String("influxdb", "http://localhost:8086", "InfluxDB URI")
	influxToken       = flag.String("influxdb-authtoken", "admin:admin", "InfluxDB authentication token (optional)")
	s3URI             = flag.String("s3", "minio-s3.minio-s3.svc.cluster.local:9000", "S3 URI")
	s3accesskeyid     = flag.String("s3-accesskey-id", "admin", "S3 access key")
	s3accesskeysecret = flag.String("s3-accesskey-secret", "admin1234", "S3 access key")
	s3bucket          = flag.String("s3-bucket", "myreports", "S3 bucket name")
)

func init() {
	initLog()
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) // TODO: configurable
}

var client influxdb2.Client

func configureInfluxDB() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client = influxdb2.NewClientWithOptions(*influxURI, *influxToken,
		influxdb2.DefaultOptions().
			SetUseGZip(true).
			SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			}))
}

func configureS3() {
	ctx := context.Background()
	minioClient, err := minio.New(*s3URI, &minio.Options{
		Creds:  credentials.NewStaticV4(*s3accesskeyid, *s3accesskeysecret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Debugf("%#v\n", minioClient) // minioClient is now setup

	location := "site1"

	err = minioClient.MakeBucket(ctx, *s3bucket, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, *s3bucket)
		if errBucketExists == nil && exists {
			log.Infof("We already own %s\n", *s3bucket)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Infof("Successfully created %s\n", *s3bucket)
	}
}

func parseTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0)
}

func readDataMessage() []string {
	queryAPI := client.QueryAPI("")
	// get QueryTableResult

	var rows []string
	result, err := queryAPI.Query(context.Background(), `from(bucket:"testdata")|> range(start: -5m) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				log.Infof("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			log.Infof("value: %v\n", result.Record().Value())
			rows = append(rows, result.Record().String())
		}
		// check for an error
		if result.Err() != nil {
			log.Infof("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

	return rows
}

func createPDF(rows []string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 11)
	for _, row := range rows {
		pdf.Write(12, row)
	}
	err := pdf.OutputFileAndClose("/report.pdf")
	if err != nil {
		log.Fatalf("error creating PDF file: %s\n", err.Error())
	}
}

func uploadPDF() {
	objectName := "report-" + strings.TrimSuffix(time.Now().UTC().Local().Format(time.RFC3339), "Z") + ".pdf"
	filePath := "/report.pdf"
	contentType := "application/pdf"

	//TODO: Extract common client init logic
	ctx := context.Background()
	minioClient, err := minio.New(*s3URI, &minio.Options{
		Creds:  credentials.NewStaticV4(*s3accesskeyid, *s3accesskeysecret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Upload the zip file with FPutObject
	n, err := minioClient.FPutObject(ctx, *s3bucket, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Infof("Successfully uploaded %s of size %d\n", objectName, n.Size)
}

func cleanup() {
	// Close Client
	log.Infof("Closing client connection")
	defer client.Close()
}

func main() {
	flag.Parse()
	log.Infoln("Start reporter...")

	configureInfluxDB()
	configureS3()

	// Watch for CTRL+C / SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	for {
		time.Sleep(1 * time.Minute)
		data := readDataMessage()
		createPDF(data)
		uploadPDF()
	}
}
