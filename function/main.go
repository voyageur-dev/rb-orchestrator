package main

import (
	"context"
	"fmt"
	"function/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	lambdaSDK "github.com/aws/aws-sdk-go-v2/service/lambda"
	"log"
	"net/http"
	"os"
)

const (
	updateMetadataPath = "PUT /rb/orchestrator/metadata"
)

var (
	questionServiceArn string
	metadataServiceArn string

	lambdaClient    *lambdaSDK.Client
	questionService *services.QuestionService
	metadataService *services.MetadataService
)

func init() {
	questionServiceArn = os.Getenv("QUESTION_SERVICE_ARN")
	metadataServiceArn = os.Getenv("METADATA_SERVICE_ARN")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	lambdaClient = lambdaSDK.NewFromConfig(cfg)
	questionService = &services.QuestionService{
		LambdaClient:       lambdaClient,
		QuestionServiceArn: questionServiceArn,
	}
	metadataService = &services.MetadataService{
		LambdaClient:       lambdaClient,
		MetadataServiceArn: metadataServiceArn,
	}
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	path := request.RouteKey

	return func() (events.APIGatewayV2HTTPResponse, error) {
		switch path {
		case updateMetadataPath:
			return updateMetadata()
		default:
			return events.APIGatewayV2HTTPResponse{
				Body:       "Path Not Found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	}()
}

func updateMetadata() (events.APIGatewayV2HTTPResponse, error) {
	fmt.Println("Update Metadata Begin")

	metadata, err := metadataService.GetMetadata()
	if err != nil {
		log.Println(fmt.Sprintf("Error getting metadata from rb-metadata-service: %v", err))
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error updating metadata",
		}, nil
	}

	var examIds []string
	for _, item := range metadata {
		examIds = append(examIds, item.ExamId)
	}

	counts, err := questionService.GetQuestionCounts(examIds)
	if err != nil {
		log.Println(fmt.Sprintf("Error getting count from rb-question-service: %v", err))
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error updating metadata",
		}, nil
	}

	updated := false
	for i := range metadata {
		item := &metadata[i]
		if count, exists := counts[item.ExamId]; exists && item.QuestionCount != count {
			fmt.Println(fmt.Sprintf("updating question count for %s from %d to %d", item.ExamId, item.QuestionCount, count))
			item.QuestionCount = count
			updated = true
		}
	}

	if updated {
		err = metadataService.UpdateMetadata(metadata)
		if err != nil {
			log.Println(fmt.Sprintf("Error updating metadata from rb-metadata-service: %v", err))
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       "Error updating metadata",
			}, nil
		}
	}

	fmt.Println("Update Metadata End")

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
