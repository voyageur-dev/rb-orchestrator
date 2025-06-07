package services

import (
	"context"
	"encoding/json"
	"function/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	lambdaSDK "github.com/aws/aws-sdk-go-v2/service/lambda"
)

type MetadataService struct {
	LambdaClient       *lambdaSDK.Client
	MetadataServiceArn string
}

const (
	getMetadataPath    = "GET /rb/metadata"
	updateMetadataPath = "PUT /rb/metadata"
)

func (svc *MetadataService) GetMetadata() ([]models.Metadata, error) {
	payload, err := json.Marshal(map[string]string{
		"routeKey": getMetadataPath,
	})
	if err != nil {
		return nil, err
	}

	resp, err := svc.LambdaClient.Invoke(context.TODO(), &lambdaSDK.InvokeInput{
		FunctionName:   aws.String(svc.MetadataServiceArn),
		InvocationType: "RequestResponse",
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	var respPayloadMap map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &respPayloadMap); err != nil {
		return nil, err
	}

	var metadata []models.Metadata
	if err := json.Unmarshal([]byte(respPayloadMap["body"].(string)), &metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func (svc *MetadataService) UpdateMetadata(metadata []models.Metadata) error {
	metadataBody, _ := json.Marshal(metadata)

	payload, err := json.Marshal(map[string]string{
		"routeKey": updateMetadataPath,
		"body":     string(metadataBody),
	})
	if err != nil {
		return err
	}

	_, err = svc.LambdaClient.Invoke(context.TODO(), &lambdaSDK.InvokeInput{
		FunctionName:   aws.String(svc.MetadataServiceArn),
		InvocationType: "RequestResponse",
		Payload:        payload,
	})
	if err != nil {
		return err
	}

	return nil
}
