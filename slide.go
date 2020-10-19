package main

import (
	"fmt"
	"io/ioutil"
)

import jsoniter "github.com/json-iterator/go"

type Slides struct {
	Slides []Slide `json:"slides"`
}

type Slide struct {
	Image         string   `json:"image"`
	Question      string   `json:"question,omitempty"`
	Answers       []string `json:"answers,omitempty"`
	CorrectAnswer int      `json:"correct,omitempty"`
}

const slides_path = "./slides.json"

func GetRawSlideJSON(filepath string) (string, error) {
	str, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	fmt.Println("Successfully opened ", filepath)
	return string(str), nil
}

func GetSlideJSON(filepath string) (Slides, error) {
	slides := Slides{}
	str, err := ioutil.ReadFile(filepath)
	if err != nil {
		return slides, err
	}
	fmt.Println("Successfully opened ", filepath)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(str, &slides); err != nil {
		return slides, err
	}
	fmt.Print(slides)
	return slides, nil
}
