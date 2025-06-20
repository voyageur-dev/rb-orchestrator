package models

type QuestionCount struct {
	ExamId string `json:"examId"`
	Count  int    `json:"count"`
}

type Question struct {
	ExamId      string   `json:"examId"`
	QuestionId  int      `json:"questionId"`
	Options     []Option `json:"options"`
	Question    string   `json:"question"`
	S3ImageUrls []string `json:"s3ImageUrls"`
}

type Option struct {
	IsCorrect   bool     `json:"isCorrect"`
	Text        string   `json:"text"`
	S3ImageUrls []string `json:"s3ImageUrls"`
}
