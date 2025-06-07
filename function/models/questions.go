package models

type QuestionCount struct {
	ExamId string `json:"examId"`
	Count  int    `json:"count"`
}
