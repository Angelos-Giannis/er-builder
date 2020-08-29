package mock

import (
	externalSurvey "github.com/AlecAivazis/survey/v2"
	"github.com/eujoy/erbuilder/internal/domain"
	"github.com/stretchr/testify/mock"
)

// Survey mock structure.
type Survey struct {
	mock.Mock
}

// AskTableDetails mock implementation.
func (s *Survey) AskTableDetails(questions []*externalSurvey.Question) (domain.TableAnswer, error) {
	args := s.MethodCalled("AskTableDetails")
	return args.Get(0).(domain.TableAnswer), args.Error(1)
}

// AskColumnDetails mock implementation.
func (s *Survey) AskColumnDetails(questions []*externalSurvey.Question) (domain.ColumnAnswer, error) {
	args := s.MethodCalled("AskColumnDetails")
	return args.Get(0).(domain.ColumnAnswer), args.Error(1)
}
