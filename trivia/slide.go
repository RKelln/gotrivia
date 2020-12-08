package trivia

import (
	"fmt"
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

// SlideList represents a list of Slides
type SlideList struct {
	Slides []Slide `json:"slides"`
}

// Slide represents one image which may have a trivia question and answers
type Slide struct {
	Image         string   `json:"image"`
	Question      string   `json:"question,omitempty"`
	Answers       []string `json:"answers,omitempty"`
	CorrectAnswer int      `json:"correct,omitempty"`
}

// AnswerKey returns the correct answers for each trivia question
// 0 = no answer, 1+ = the correct answer
func (s *SlideList) AnswerKey() []int {
	answers := make([]int, len(s.Slides))
	for i, slide := range s.Slides {
		answers[i] = slide.CorrectAnswer
	}
	return answers
}

// GetRawSlideJSON returns the JSON as a string
func GetRawSlideJSON(filepath string) (string, error) {
	str, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	fmt.Println("Successfully opened ", filepath)
	return string(str), nil
}

// GetSlideJSON returns a SlideList built from the JSON
func GetSlideJSON(filepath string) (*SlideList, error) {
	slides := &SlideList{}
	str, err := ioutil.ReadFile(filepath)
	if err != nil {
		return slides, err
	}
	fmt.Println("Successfully opened ", filepath)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(str, slides); err != nil {
		return slides, err
	}
	return slides, nil
}
