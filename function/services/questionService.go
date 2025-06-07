package services

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	lambdaSDK "github.com/aws/aws-sdk-go-v2/service/lambda"
	"strings"
)

type QuestionService struct {
	LambdaClient       *lambdaSDK.Client
	QuestionServiceArn string
}

const (
	getQuestionCountPath = "POST /rb/questions/count"
)

func (svc *QuestionService) GetQuestionCounts(examIds []string) (map[string]int, error) {
	payload, err := json.Marshal(map[string]string{
		"routeKey": getQuestionCountPath,
		"body":     strings.Join(examIds, ","),
	})
	if err != nil {
		return nil, err
	}

	resp, err := svc.LambdaClient.Invoke(context.TODO(), &lambdaSDK.InvokeInput{
		FunctionName:   aws.String(svc.QuestionServiceArn),
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

	var counts map[string]int
	if err := json.Unmarshal([]byte(respPayloadMap["body"].(string)), &counts); err != nil {
		return nil, err
	}

	return counts, nil
}
