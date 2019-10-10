package main

import (
	"fmt"
	"math/rand"
)

type CardId int

const N = 52

type Cards struct {
	kinds   [N]int
	values  [N]int
	symbols [N]string
}

func (c *Cards) kind(i CardId) int {
	return c.kinds[i]
}

func (c *Cards) symbol(i CardId) string {
	return c.symbols[i]
}

func (c *Cards) value(i CardId) int {
	return c.values[i]
}

type GameState int

const (
	GsNone     = GameState(0)
	GsPlanning = GameState(1)
	GsPlaying  = GameState(2)
	GsFinished = GameState(3)
)

type Game struct {
	state         GameState
	players       []string
	hands         map[string][]CardId
	activePlayer  int
	playerByOrder []int
	roundsLeft    uint
	cardsPerHand  uint
	estimatedWins []uint
	totalWins     uint
	started       bool
	deck          []CardId
}

func NewGame(players []string) *Game {
	if len(players) < 2 {
		return nil
	}

	game := &Game{
		players:       players,
		deck:          make([]CardId, N),
		estimatedWins: make([]uint, 0),
		playerByOrder: make([]int, len(players)),
	}
	for i := 0; i < len(players); i++ {
		game.playerByOrder[i] = i
	}
	for i := 0; i < N; i++ {
		game.deck[i] = CardId(i)
	}

	game.init()
	game.shufflePlayers()
	game.shuffleDeck()

	return game
}

func (g *Game) init() {
	g.roundsLeft = N / uint(len(g.players))
	g.started = false
	g.activePlayer = 1
}

func (g *Game) shufflePlayers() {

}

func (g *Game) shuffleDeck() {
	for i := range g.deck {
		j := rand.Intn(i + 1)
		g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
	}
}

func (g *Game) DealCards() {
	if g.roundsLeft == 0 {
		g.state = GsFinished
	}
	if g.started {
		g.roundsLeft--
	} else {
		g.started = true
	}
	g.totalWins = 0
	g.cardsPerHand = g.roundsLeft
	g.state = GsPlanning

}

func (g *Game) IsPlanningPhase() bool {
	return g.state == GsPlanning
}

func (g *Game) ActivePlayer() int {
	return 0
}

func (g *Game) Plan(estimatedWins uint) bool {
	if estimatedWins > g.cardsPerHand {
		return false
	}

	lastPlayer := g.activePlayer == len(g.players)-1
	sumsUpToHand := g.totalWins+estimatedWins == g.cardsPerHand
	if lastPlayer && sumsUpToHand {
		return false
	}

	g.totalWins += estimatedWins
	g.activePlayer++

	return true
}

func (g *Game) IsPlayingPhase() bool {
	return g.state == GsPlaying
}

func (g *Game) IsFinished() bool {
	return g.state == GsFinished
}

func (g *Game) PlayCard(card int) bool {
	return true
}

func main() {
	cards := &Cards{
		kinds: [N]int{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
			4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
		},
		values: [N]int{
			2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
			2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
			2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
			2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
		},
		symbols: [N]string{
			"♠️2", "♠️3", "♠️4", "♠️5", "♠️6", "♠7️", "♠️8", "♠️9", "♠️10", "♠️J", "♠️Q", "♠️K", "♠️A",
			"♥️2", "♥️3", "♥️4", "♥️5", "♥️6", "♠7", "♥️8", "♥️9", "♥️10", "♥️J", "♥️Q", "♥️K", "♥️A",
			"♣️2", "♣️3", "♣️4", "♣️5", "♣️6", "♣️7", "♣️8", "♣️9", "♣️10", "♣️J", "♣️Q", "♣️K", "♣️A",
			"♦️2", "♦️3", "♦️4", "♦5", "♦️6", "♦️7", "♦️8", "♦️9", "♦️10", "♦️J", "♦️Q", "♦️K", "♦️A",
		},
	}

	game := NewGame([]string{
		"bart",
		"bboruta",
		"czana",
		"kijek",
		"kovson",
		"marek",
	})

	for {
		game.DealCards()

		fmt.Print("Round left: ", game.roundsLeft, "\n")
		for i := 0; i < len(game.players); i++ {
			fmt.Print("P", i, ": ", game.players[i])
			for j := 0; j < len(game.hands[game.players[i]]); j++ {
				fmt.Print(cards.symbol(game.hands[game.players[i]][j]), ", ")
			}
			fmt.Print("\n")
		}

		for game.IsPlanningPhase() {
			fmt.Print("P: %d\n", game.ActivePlayer())
			var wins uint
			_, err := fmt.Scanf("%d", &wins)
			if err == nil {
				if !game.Plan(wins) {
					fmt.Println("Wrong estimate!")
					continue
				}
			} else {
				continue
			}
		}
		for game.IsPlayingPhase() {
			fmt.Print("P: %d\n", game.ActivePlayer())
			card := 1
			if !game.PlayCard(card) {
				fmt.Println("Wrong card!")
				continue
			}
		}
		if game.IsFinished() {
			fmt.Println("End")
			break
		}
	}
}
