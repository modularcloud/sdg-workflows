package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

var facts = []string{
	"Octopuses have three hearts and blue blood.",
	"Honey never spoils; edible honey has been found in ancient Egyptian tombs.",
	"Bananas are berries, but strawberries are not.",
	"A group of flamingos is called a flamboyance.",
	"The Eiffel Tower can grow more than 15 cm taller in summer due to thermal expansion.",
	"Sharks existed before trees did.",
	"Wombat poop is cube-shaped.",
}

var quotes = []string{
	"The only way to do great work is to love what you do. — Steve Jobs",
	"In the middle of difficulty lies opportunity. — Albert Einstein",
	"Whether you think you can or you think you can't, you're right. — Henry Ford",
	"Simplicity is the ultimate sophistication. — Leonardo da Vinci",
	"The best way to predict the future is to invent it. — Alan Kay",
}

var jokes = []string{
	"Why do programmers prefer dark mode? Because light attracts bugs.",
	"There are 10 kinds of people in the world: those who understand binary and those who don't.",
	"I would tell you a UDP joke, but you might not get it.",
	"Why did the developer go broke? Because he used up all his cache.",
	"How many programmers does it take to change a light bulb? None — that's a hardware problem.",
}

var eightBall = []string{
	"It is certain.",
	"Without a doubt.",
	"Yes, definitely.",
	"You may rely on it.",
	"As I see it, yes.",
	"Most likely.",
	"Outlook good.",
	"Yes.",
	"Signs point to yes.",
	"Reply hazy, try again.",
	"Ask again later.",
	"Better not tell you now.",
	"Cannot predict now.",
	"Concentrate and ask again.",
	"Don't count on it.",
	"My reply is no.",
	"My sources say no.",
	"Outlook not so good.",
	"Very doubtful.",
}

type state struct {
	Recent map[string][]int `json:"recent"`
}

func statePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".random-cli", "state.json")
}

func loadState() *state {
	s := &state{Recent: map[string][]int{}}
	path := statePath()
	if path == "" {
		return s
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	_ = json.Unmarshal(data, s)
	if s.Recent == nil {
		s.Recent = map[string][]int{}
	}
	return s
}

func saveState(s *state) {
	path := statePath()
	if path == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.Marshal(s)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o644)
}

// pickNoRepeat picks an index from items, avoiding recently-used indices for
// the given category. Returns the chosen string.
func pickNoRepeat(s *state, category string, items []string) string {
	recent := s.Recent[category]
	avoid := map[int]bool{}
	// Keep memory under half the list so there's always something to pick.
	maxMem := len(items) / 2
	if maxMem < 1 {
		maxMem = 0
	}
	if len(recent) > maxMem {
		recent = recent[:maxMem]
	}
	for _, i := range recent {
		avoid[i] = true
	}

	candidates := make([]int, 0, len(items))
	for i := range items {
		if !avoid[i] {
			candidates = append(candidates, i)
		}
	}
	if len(candidates) == 0 {
		candidates = make([]int, len(items))
		for i := range items {
			candidates[i] = i
		}
	}

	idx := candidates[rand.Intn(len(candidates))]
	s.Recent[category] = append([]int{idx}, recent...)
	if len(s.Recent[category]) > maxMem {
		s.Recent[category] = s.Recent[category][:maxMem]
	}
	return items[idx]
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: random-cli <fact|quote|joke|random|8ball [question]>")
	os.Exit(2)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	s := loadState()
	defer saveState(s)

	switch os.Args[1] {
	case "fact":
		fmt.Println(pickNoRepeat(s, "fact", facts))
	case "quote":
		fmt.Println(pickNoRepeat(s, "quote", quotes))
	case "joke":
		fmt.Println(pickNoRepeat(s, "joke", jokes))
	case "random":
		categories := []struct {
			name  string
			items []string
		}{
			{"fact", facts},
			{"quote", quotes},
			{"joke", jokes},
		}
		c := categories[rand.Intn(len(categories))]
		fmt.Println(pickNoRepeat(s, c.name, c.items))
	case "8ball":
		// 8ball ignores no-repeat memory — feels more authentic.
		fmt.Println(eightBall[rand.Intn(len(eightBall))])
	default:
		usage()
	}
}
