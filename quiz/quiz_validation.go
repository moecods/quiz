package quiz

import (
	"moecods/quiz/utils"
)

type QuizValidation struct {
	quiz *Quiz
}

func (v *QuizValidation) QuizValidation() error {
	roles := utils.NewRule();

	err := roles.ValidateRequired("title", v.quiz.Title)
	if err != nil {
		return err
	}

	err = roles.ValidateMinLength("title", v.quiz.Title, 3)
	if err != nil {
		return err
	}

	return nil
}