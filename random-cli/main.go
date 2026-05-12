package main

import (
	"fmt"
	"math/rand"
	"os"
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

func pick(items []string) string {
	return items[rand.Intn(len(items))]
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: random-cli <fact|quote|joke>")
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	switch os.Args[1] {
	case "fact":
		fmt.Println(pick(facts))
	case "quote":
		fmt.Println(pick(quotes))
	case "joke":
		fmt.Println(pick(jokes))
	default:
		usage()
	}
}
