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
	Answers []int  `json:"answers,omitempty"`
}

type MyGame struct {
	SlideList
	Answers  []int     `json:"answers"`
	Results  []int     `json:"results"`
	Correct  []int     `json:"correct"`
	Rankings []Ranking `json:"leaderboard"`
}

type Ranking struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
}

func (g *Game) AddPlayer(p Player) error {
	if _, found := FindPlayer(g.Players, p.Name); found {
		return fmt.Errorf("Player %v already exists in game", p.Name)
	}

	// ensure player has answers
	if len(p.Answers) != len(g.Slides) {
		p.Answers = make([]int, len(g.Slides))
	}

	g.Players = append(g.Players, p)

	return nil
}

func (g *Game) AddAnswer(name string, slide int, answer int) error {
	i, found := FindPlayer(g.Players, name)
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
	myGame := &MyGame{}
	var err error

	i, found := FindPlayer(g.Players, name)
	if !found {
		return myGame, fmt.Errorf("Player %v doesn't exists in game", name)
	}

	myGame.Slides = g.Slides
	myGame.Answers = g.Players[i].Answers
	myGame.Results, err = g.Players[i].Results(g.AnswerKey())
	myGame.Correct, myGame.Rankings = g.Results()

	return myGame, err
}

func (p *Player) Results(correct []int) ([]int, error) {
	results := make([]int, len(correct))
	if len(correct) != len(p.Answers) {
		return results, fmt.Errorf("Player %v answers to not match answer key", p.Name)
	}
	for i := range correct {
		if correct[i] > 0 {
			if correct[i] == p.Answers[i] {
				results[i] = 1
			} else {
				results[i] = -1
			}
		}
	}
	return results, nil
}

func (g *Game) Results() ([]int, []Ranking) {
	answers := g.AnswerKey()
	answer_results := make([]int, len(answers))
	player_results := make([][]int, len(g.Players))
	rankings := make([]Ranking, 0, len(g.Players))

	for i := range g.Players {
		player_results[i], _ = g.Players[i].Results(answers)
	}
	for i, player_correct := range player_results {
		sum := 0
		for j := range player_correct {
			if player_correct[j] > 0 {
				sum += player_correct[j]
				answer_results[j] += 1
			}
		}
		rankings = append(rankings, Ranking{Name: g.Players[i].Name, Points: sum})
	}

	sort.SliceStable(rankings, func(i, j int) bool {
		return rankings[i].Points < rankings[j].Points
	})

	return answer_results, rankings
}

func (g *Game) Save(filepath string) error {
	fmt.Println("Saving game:", filepath)

	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(g)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
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
func NewGame(game *Game, slides *SlideList) error {
	// TODO: check len of slides
	if len(game.Slides) == 0 {
		game.Slides = slides.Slides
	}

	// remove all players?
	if len(game.Players) > 0 {
		fmt.Print("Existing players: ")
		names := []string{}
		for _, p := range game.Players {
			names = append(names, p.Name)
		}
		fmt.Println(strings.Join(names, ", "))
	}

	return nil
}

func FindPlayer(players []Player, name string) (int, bool) {
	for i, p := range players {
		if p.Name == name {
			return i, true
		}
	}
	return -1, false
}
