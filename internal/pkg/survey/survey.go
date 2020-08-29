package survey

import (
	externalSurvey "github.com/AlecAivazis/survey/v2"
	"github.com/eujoy/erbuilder/internal/domain"
)

// Survey describes the survey package wrapper structure.
type Survey struct{}

// New creates and returns a survey instance.
func New() *Survey {
	return &Survey{}
}

// AskTableDetails is a wrapper function for survey.Ask to ask for table details.
func (s *Survey) AskTableDetails(questions []*externalSurvey.Question) (domain.TableAnswer, error) {
	var tableAnswer domain.TableAnswer
	// perform the questions provided.
	err := externalSurvey.Ask(questions, &tableAnswer)
	if err != nil {
		return domain.TableAnswer{}, err
	}

	return tableAnswer, nil
}

// AskColumnDetails is a wrapper function for survey.Ask to ask for table details.
func (s *Survey) AskColumnDetails(questions []*externalSurvey.Question) (domain.ColumnAnswer, error) {
	var columnAnswer domain.ColumnAnswer
	// perform the questions provided.
	err := externalSurvey.Ask(questions, &columnAnswer)
	if err != nil {
		return domain.ColumnAnswer{}, err
	}

	return columnAnswer, nil
}
