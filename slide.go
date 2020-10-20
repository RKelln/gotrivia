package main

import (
	"fmt"
	"io/ioutil"
)

import jsoniter "github.com/json-iterator/go"

const slides_path = "./slides.json"

type Slides struct {
	Slides []Slide `json:"slides"`
}

type Slide struct {
	Image         string   `json:"image"`
	Question      string   `json:"question,omitempty"`
	Answers       []string `json:"answers,omitempty"`
	CorrectAnswer int      `json:"correct,omitempty"`
}

func (s *Slides) answerKey() []int {
	answers := make([]int, len(s.Slides))
	for i, slide := range s.Slides {
		answers[i] = slide.CorrectAnswer
	}
	return answers
}

func GetRawSlideJSON(filepath string) (string, error) {
	str, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	fmt.Println("Successfully opened ", filepath)
	return string(str), nil
}

func GetSlideJSON(filepath string) (*Slides, error) {
	slides := &Slides{}
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
