package models

type Question struct {
	ID           int      `json:"id"`
	QuestionText string   `json:"questions_text"`
	Options      []string `json:"options"`
	CorrectIndex int      `json:"correct_index"`
}
