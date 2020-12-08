package trivia

import (
	"fmt"
)

// Player represents the  player and their answers for the Game.
type Player struct {
	Name    string `json:"name"`
	Answers []int  `json:"answers,omitempty"` // player's answers (which answer they answered, 0 = no answer)
}

// Results returns the results for a player in this format:
// [answerindex] = -1, 0, or 1
// where -1 = incorrect, 0 = no question or unanswered, 1 = correctly answered
// Thus you can sum the players score by adding all the values in the returned array.
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
