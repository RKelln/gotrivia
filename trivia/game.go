package trivia

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// Game represents the active trivia game.
// Only one game can be active at a time
type Game struct {
	Players []Player `json:"players"`
	SlideList
}

// MyGame represents the data sent back to Players, customized to that player.
type MyGame struct {
	SlideList
	Answers   []int     `json:"answers"`     // player's answers (which answer they answered, 0 = no answer)
	Results   []int     `json:"results"`     // player's answers results (correct = 1, incorrect = -1, n/a = 0)
	Completed []int     `json:"completed"`   // number of players who have completed answers
	Correct   []int     `json:"correct"`     // number of players who got the correct answer
	Rankings  []Ranking `json:"leaderboard"` // all player rankings
}

// GameStatus is used for general game information for all players.
// Used for administration and overview of the game status.
type GameStatus struct {
	SlideList
	Players   []Player  `json:"players"`
	Completed []int     `json:"completed"`   // number of players who have completed answers
	Correct   []int     `json:"correct"`     // number of players who got the correct answer
	Rankings  []Ranking `json:"leaderboard"` // all player rankings
}

// Ranking represents an row on the leaderboard
type Ranking struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
}

// Add player to the game
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

// FindPlayer returns player index and true, or -1 and false if not found.
func (g *Game) FindPlayer(name string) (int, bool) {
	for i, p := range g.Players {
		if p.Name == name {
			return i, true
		}
	}
	return -1, false
}

// AddAnswer returns nil if successful or an error.
// This is used to add/set the answer to a question for the named player.
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

// ForPlayer returns a MyGame struct and error if it fail.
// This is used to get the game information for a particular player.
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

// Status returns a GameStatus struct and error on failure.
// This is used to get the overall game status/information for all players.
func (g *Game) Status() *GameStatus {
	status := &GameStatus{}
	status.Players = g.Players
	status.Slides = g.Slides

	// other player results
	status.Correct, status.Completed, status.Rankings = g.Results()

	return status
}

// Results returns the statistics for:
// the correct answers, all answered questions, and the player rankings
// Player rankings are sorted from highest to lowest score, and do not handle ties in any particular way
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

// Save saves the game state to a json file
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

// GetGameJSON creates a game from a json file
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

// NewGame initializes a Game from a SlideList
// If the existing slides in the Game do not match the slides in the SlideList
// then reset slides and players.
func NewGame(g *Game, s *SlideList) error {

	if len(g.Slides) == 0 {
		g.Slides = s.Slides
	}

	if len(g.Slides) != len(s.Slides) {
		g.Slides = s.Slides
		g.Players = nil // clear existing players
	}

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
