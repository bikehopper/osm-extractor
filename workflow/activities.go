package osm_extractor_workflow

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.temporal.io/sdk/activity"
)

type Extracts struct {
	Extracts []Extract `json:"Extracts"`
}

type Extract struct {
	Output      string  `json:"output"`
	Directory   string  `json:"directory"`
	Description string  `json:"description"`
	Polygon     polygon `json:"polygon"`
}

type polygon struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
}

func ExtractOsmCutoutsActivity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("extracting OSM cutouts")
	outputDir := getEnv("OUTPUT_DIR", "/mnt/output")
	pbfPath := getEnv("PBF_PATH", "/mnt/input/latest.osm.pbf")
	configPath := getEnv("CONFIG_PATH", "../config.json")

	cmd := exec.Command("osmium", "extract", "--overwrite", "-d", outputDir, "-c", configPath, pbfPath)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Info("error getting osmium standerr pipe")
		return err
	}
	err = cmd.Start()
	if err != nil {
		logger.Info("Osmium extract cmd.Start() failed with %s\n", err)
		return err
	}

	stderrin := bufio.NewScanner(stderr)
	for stderrin.Scan() {
		logger.Info(stderrin.Text())
	}
	err = cmd.Wait()
	if err != nil {
		logger.Info("Osmium failed", err)
		return err
	}

	return nil
}

func UploadOsmCutoutsActivity(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("uploading extracts to bucket")

	outputDir := getEnv("OUTPUT_DIR", "/mnt/output")
	bucket := getEnv("BUCKET", "osm-extracts")
	s3Region := getEnv("S3_REGION", "us-east-1")
	s3EndpointUrl := getEnv("S3_ENDPOINT_URL", "")
	scheduledDate := activity.GetInfo(ctx).ScheduledTime

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = s3Region
		o.BaseEndpoint = aws.String(s3EndpointUrl)
		o.UsePathStyle = true
	})
	uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024
		u.Concurrency = 5
	})

	extracts, err := readExtractFile("config.json")
	if err != nil {
		logger.Error("error reading config.json")
		return err
	}

	for i := 0; i < len(extracts.Extracts); i++ {
		extractedFile, err := os.Open(filepath.Join(outputDir, extracts.Extracts[i].Output))
		if err != nil {
			logger.Error("failed to open file: %v", err)
			return err
		}
		defer extractedFile.Close()

		extractedFileStats, err := extractedFile.Stat()
		if err != nil {
			logger.Error("failed to stat file: %v", err)
			return err
		}

		destPathDated := filepath.Join(extracts.Extracts[i].Directory, getDatedFileName(extracts.Extracts[i].Output, scheduledDate))
		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket:             aws.String(bucket),
			Key:                aws.String(destPathDated),
			Body:               io.Reader(extractedFile),
			ContentDisposition: aws.String("application/octet-stream"),
			ContentLength:      aws.Int64(extractedFileStats.Size()),
		})
		if err != nil {
			logger.Error("error while uploading file: %v", err)
			return err
		}

		fmt.Println("File uploaded successfully.")
	}

	return nil
}

func CopyOsmCutouts(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Creating 'latest' copies of extracts")
	scheduledDate := activity.GetInfo(ctx).ScheduledTime

	bucket := getEnv("BUCKET", "osm-extracts")
	s3Region := getEnv("S3_REGION", "us-east-1")
	s3EndpointUrl := getEnv("S3_ENDPOINT_URL", "")

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = s3Region
		o.BaseEndpoint = aws.String(s3EndpointUrl)
		o.UsePathStyle = true
	})

	extracts, err := readExtractFile("config.json")
	if err != nil {
		logger.Error("error reading config.json")
		return err
	}

	for i := 0; i < len(extracts.Extracts); i++ {
		destPathDated := filepath.Join(extracts.Extracts[i].Directory, getDatedFileName(extracts.Extracts[i].Output, scheduledDate))
		destPathLatest := filepath.Join(extracts.Extracts[i].Directory, getLatestFileName(extracts.Extracts[i].Output))
		_, err = s3Client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:             aws.String(bucket),
			CopySource:         aws.String(filepath.Join(bucket, destPathDated)),
			Key:                aws.String(destPathLatest),
			ContentDisposition: aws.String("application/octet-stream"),
		})
		if err != nil {
			logger.Error("error while copying file in bucket")
			return err
		}

		fmt.Println("File moved successfully.")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func readExtractFile(filePath string) (Extracts, error) {
	jsonFile, err := os.Open(filePath)

	if err != nil {
		return Extracts{}, err
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var extracts Extracts

	err = json.Unmarshal(byteValue, &extracts)
	if err != nil {
		return Extracts{}, err
	}

	return extracts, nil
}

func getDatedFileName(fileName string, date time.Time) string {
	return date.Format("2006-01-02") + "-" + fileName
}

func getLatestFileName(fileName string) string {
	return "latest-" + fileName
}
