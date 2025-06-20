package services

import (
	"context"
	"encoding/json"
	"function/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	lambdaSDK "github.com/aws/aws-sdk-go-v2/service/lambda"
)

type AskService struct {
	LambdaClient  *lambdaSDK.Client
	AskServiceArn string
}

const (
	getAnalysisPath    = "GET /rb/ask/{examId}/{questionId}"
	createAnalysisPath = "POST /rb/ask"
)

func (svc *AskService) GetAnalysis(examId string, questionId string) (map[string]interface{}, error) {
	params, _ := json.Marshal(map[string]string{
		"examId":     examId,
		"questionId": questionId,
	})
	payload, err := json.Marshal(map[string]string{
		"routeKey":       getAnalysisPath,
		"pathParameters": string(params),
	})
	if err != nil {
		return nil, err
	}

	resp, err := svc.LambdaClient.Invoke(context.TODO(), &lambdaSDK.InvokeInput{
		FunctionName:   aws.String(svc.AskServiceArn),
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

	var analysis map[string]interface{}
	if err := json.Unmarshal([]byte(respPayloadMap["body"].(string)), &analysis); err != nil {
		return nil, err
	}

	return analysis, nil
}

func (svc *AskService) CreateAnalysis(question models.Question, model string) error {
	return nil
}
