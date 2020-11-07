package trivia

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type Game struct {
	Players []Player `json:"players"`
	SlideList
}

type Player struct {
	Name    string `json:"name"`
	Answers []int  `json:"answers,omitempty"` // player's answers (which answer they answered, 0 = no answer)
}

type MyGame struct {
	SlideList
	Answers   []int     `json:"answers"`     // player's answers (which answer they answered, 0 = no answer)
	Results   []int     `json:"results"`     // player's answers results (correct = 1, incorrect = -1, n/a = 0)
	Completed []int     `json:"completed"`   // number of players who have completed answers
	Correct   []int     `json:"correct"`     // number of players who got the correct answer
	Rankings  []Ranking `json:"leaderboard"` // player rankings
}

type GameStatus struct {
	SlideList
	Players   []Player  `json:"players"`
	Completed []int     `json:"completed"`   // number of players who have completed answers
	Correct   []int     `json:"correct"`     // number of players who got the correct answer
	Rankings  []Ranking `json:"leaderboard"` // player rankings
}

type Ranking struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
}

func (g *Game) AddPlayer(p Player) error {
	if _, found := g.FindPlayer(p.Name); found {
		return fmt.Errorf("Player %v already exists in game", p.Name)
	}

	// ensure player has answers
	if len(p.Answers) != len(g.Slides) {
		p.Answers = make([]int, len(g.Slides))
	}

	g.Players = append(g.Players, p)

	return nil
}

func (g *Game) FindPlayer(name string) (int, bool) {
	for i, p := range g.Players {
		if p.Name == name {
			return i, true
		}
	}
	return -1, false
}

func (g *Game) AddAnswer(name string, slide int, answer int) error {
	i, found := g.FindPlayer(name)
	if !found {
		return fmt.Errorf("Player %v doesn't exists in game", name)
	}

	if slide < 0 || slide >= len(g.Players[i].Answers) {
		return fmt.Errorf("Invalid slide %v", slide)
	}

	if answer <= 0 || answer > len(g.Slides[slide].Answers) {
		return fmt.Errorf("Invalid answer %v for slide %v", answer, slide)
	}

	g.Players[i].Answers[slide] = answer

	return nil
}

func (g *Game) ForPlayer(name string) (*MyGame, error) {
	var err error
	myGame := &MyGame{}

	i, found := g.FindPlayer(name)
	if !found {
		return myGame, fmt.Errorf("Player %v doesn't exists in game", name)
	}

	answerKey := g.AnswerKey()

	// my answers
	myGame.Answers = g.Players[i].Answers

	// my slides (with unanswered correct answer info removed)
	myGame.Slides = make([]Slide, len(g.Slides))
	copy(myGame.Slides, g.Slides)
	// hide correct answers for unanswered questiosn
	for i := range myGame.Slides {
		if myGame.Answers[i] == 0 && answerKey[i] > 0 {
			myGame.Slides[i].CorrectAnswer = 0
		}
	}

	// my results
	myGame.Results, err = g.Players[i].Results(answerKey)
	if err != nil {
		return myGame, err
	}

	// other player results
	myGame.Correct, myGame.Completed, myGame.Rankings = g.Results()

	return myGame, err
}

func (g *Game) Status() (*GameStatus, error) {
	var err error
	status := &GameStatus{}
	status.Players = g.Players
	status.Slides = g.Slides

	// other player results
	status.Correct, status.Completed, status.Rankings = g.Results()

	return status, err
}

func (p *Player) Results(correct []int) ([]int, error) {
	results := make([]int, len(correct))
	if len(correct) != len(p.Answers) {
		return results, fmt.Errorf("Player %v answers to not match answer key", p.Name)
	}

	for i := range correct {
		// if player has an answer and there is a correct answer
		if p.Answers[i] > 0 && correct[i] > 0 {
			if correct[i] == p.Answers[i] {
				results[i] = 1
			} else {
				results[i] = -1
			}
		}
	}

	return results, nil
}

func (g *Game) Results() ([]int, []int, []Ranking) {
	answers := g.AnswerKey()
	all_correct := make([]int, len(answers))
	all_answered := make([]int, len(answers))
	all_player_results := make([][]int, len(g.Players))
	rankings := make([]Ranking, 0, len(g.Players))

	for i := range g.Players {
		all_player_results[i], _ = g.Players[i].Results(answers)
	}
	for i, player_results := range all_player_results {
		sum_correct := 0
		for j := range player_results {
			if player_results[j] > 0 {
				sum_correct += player_results[j]
				all_correct[j] += 1
			}
			if player_results[j] != 0 {
				all_answered[j] += 1
			}
		}
		rankings = append(rankings, Ranking{Name: g.Players[i].Name, Points: sum_correct})
	}

	// highest to lowest rankings
	sort.SliceStable(rankings, func(i, j int) bool {
		return rankings[i].Points > rankings[j].Points
	})

	return all_correct, all_answered, rankings
}

func (g *Game) Save(filepath string) error {
	fmt.Println("Saving game:", filepath)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	jsonString, err := json.MarshalIndent(g, "", " ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath, jsonString, 0644); err != nil {
		return err
	}

	return nil
}

func GetGameJSON(filepath string) (*Game, error) {
	game := &Game{}

	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return game, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return game, err
	}

	fmt.Println("Successfully read ", filepath)

	if len(b) > 0 {
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		if err := json.Unmarshal(b, game); err != nil {
			return game, err
		}
	}

	return game, nil
}

// adds slides to game
func NewGame(g *Game, s *SlideList) error {

	if len(g.Slides) == 0 {
		g.Slides = s.Slides
	}

	if len(g.Slides) != len(s.Slides) {
		g.Slides = s.Slides
		g.Players = nil // clear existing players
	}

	// remove all players?
	if len(g.Players) > 0 {
		fmt.Print("Existing players: ")
		names := []string{}
		for _, p := range g.Players {
			names = append(names, p.Name)
		}
		fmt.Println(strings.Join(names, ", "))
	}

	return nil
}
