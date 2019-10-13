package main

import (
	"fmt"
	"math/rand"
)

type CardId int

const N = 52

func maxIntIndex(ints []int) int {
	index := 0
	for i := 0; i < len(ints); i++ {
		if ints[i] > ints[index] {
			index = i
		}
	}

	return index
}

type Cards struct {
	kinds   [N]int
	values  [N]int
	symbols [N]string
}

type PlayedCard struct {
	playerName string
	card       CardId
}

type GameState int

const (
	GsDealing  = GameState(0)
	GsPlanning = GameState(1)
	GsPlaying  = GameState(2)
	GsFinished = GameState(3)
)

type Game struct {
	players            *Players
	state              GameState
	estimatedWinsCount int
	playedCards        []PlayedCard
	roundsLeft         uint
	cardsPerHand       uint
	totalEstimatedWins uint
	started            bool
	deck               []CardId
}

func NewGame(players []string) *Game {
	if len(players) < 2 {
		return nil
	}

	game := &Game{
		players: NewPlayers(players),
		deck:    make([]CardId, N),
		state:   GsDealing,
	}

	for i := 0; i < N; i++ {
		game.deck[i] = CardId(i)
	}

	game.init()

	return game
}

func (g *Game) init() {
	g.roundsLeft = N / uint(g.players.Len())
	g.roundsLeft = 3
	g.cardsPerHand = g.roundsLeft
	g.started = false
	g.players.active = g.players.dealer
}

func (g *Game) dealCards() {
	g.shuffleDeck()
	for i := 0; i < g.players.Len(); i++ {
		hand := g.deck[(uint(i) * g.cardsPerHand):(uint(i+1) * g.cardsPerHand)]
		g.players.At(i).hand = hand
	}
}

func (g *Game) shuffleDeck() {
	for i := range g.deck {
		j := rand.Intn(i + 1)
		g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
	}
}

func (g *Game) DealCards() {
	if g.started {
		g.roundsLeft--
		g.cardsPerHand = g.roundsLeft
		g.players.Next()
	} else {
		g.started = true
	}
	g.dealCards()

	g.totalEstimatedWins = 0
	g.estimatedWinsCount = 0

	g.players.active = g.players.dealer
	g.playedCards = make([]PlayedCard, 0)
	g.state = GsPlanning
}

func (g *Game) Plan(estimatedWins uint) bool {
	if estimatedWins > g.cardsPerHand {
		return false
	}

	lastPlayer := g.estimatedWinsCount == g.players.Len()-1
	sumsUpToHand := g.totalEstimatedWins+estimatedWins == g.cardsPerHand
	if lastPlayer && sumsUpToHand {
		return false
	}

	g.estimatedWinsCount++
	g.players.Current().estimatedWins = estimatedWins
	g.totalEstimatedWins += estimatedWins
	g.players.Next()

	if lastPlayer {
		g.state = GsPlaying
	}

	return true
}

func (g *Game) finishTurn(cardsData *Cards) {
	kind := cardsData.kinds[g.playedCards[0].card]
	values := make([]int, 0)
	valuesBy := make([]string, 0)

	// only consider same kind cards
	for i := 0; i < len(g.playedCards); i++ {
		played := g.playedCards[i]
		valuesBy = append(valuesBy, played.playerName)
		if kind == cardsData.kinds[played.card] {
			values = append(values, cardsData.values[played.card])
		} else {
			values = append(values, 0)
		}
	}

	maxi := maxIntIndex(values)
	winningPlayer := valuesBy[maxi]
	g.players.Win(winningPlayer)
	g.playedCards = make([]PlayedCard, 0)
}

func (g *Game) finishRound() {
	g.players.CalcRoundScores()
	if g.roundsLeft == 1 {
		g.state = GsFinished
	} else {
		g.state = GsDealing
	}
}

func (g *Game) PlayCardAt(cardsData *Cards, index int) bool {
	hand := g.players.Current().hand

	canPlayCard := false

	// player is first
	if len(g.playedCards) == 0 {
		canPlayCard = true
	} else {
		// check if has given kind
		allowedCardsIndexes := make([]int, 0)
		for i := 0; i < len(hand); i++ {
			if cardsData.kinds[hand[i]] == cardsData.kinds[g.playedCards[0].card] {
				allowedCardsIndexes = append(allowedCardsIndexes, i)
			}
		}

		// no given kind
		if len(allowedCardsIndexes) == 0 {
			canPlayCard = true
		} else {
			canPlayCard = intFind(allowedCardsIndexes, index) != -1
		}
	}

	if !canPlayCard {
		return false
	}

	card := hand[index]
	begin := hand[:index]
	end := hand[index+1:]
	hand = append(begin, end...)
	g.players.Current().hand = hand

	g.playedCards = append(g.playedCards, PlayedCard{g.players.Current().name, card})
	g.players.Next()

	if len(g.playedCards) == g.players.Len() {
		g.finishTurn(cardsData)
		if g.players.Current().HasEmptyHand() {
			g.finishRound()
		}
	}

	return true
}

//

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
			"♠️2", "♠️3", "♠️4", "♠️5", "♠️6", "♠️7", "♠️8", "♠️9", "♠️10", "♠️J", "♠️Q", "♠️K", "♠️A",
			"♥️2", "♥️3", "♥️4", "♥️5", "♥️6", "♥️7", "♥️8", "♥️9", "♥️10", "♥️J", "♥️Q", "♥️K", "♥️A",
			"♣️2", "♣️3", "♣️4", "♣️5", "♣️6", "♣️7", "♣️8", "♣️9", "♣️10", "♣️J", "♣️Q", "♣️K", "♣️A",
			"♦️2", "♦️3", "♦️4", "♦️5", "♦️6", "♦️7", "♦️8", "♦️9", "♦️10", "♦️J", "♦️Q", "♦️K", "♦️A",
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
		fmt.Println("===== dealing cards ====")
		if game.state == GsDealing {
			game.DealCards()

			fmt.Print("Rounds left: ", game.roundsLeft, "\n")
			for i := 0; i < game.players.Len(); i++ {
				dealer := " "
				if game.players.dealer == i {
					dealer = "*"
				}
				fmt.Print("P", i, ": ", fmt.Sprintf("%10s%s: ", game.players.At(i).name, dealer))
				for j := 0; j < len(game.players.At(i).hand); j++ {
					fmt.Print(
						fmt.Sprintf(
							"%4s",
							cards.symbols[game.players.At(i).hand[j]],
						),
						", ",
					)
				}
				fmt.Print("\n")
			}
		}

		fmt.Println("===== planning phase ====")
		for game.state == GsPlanning {
			fmt.Print(fmt.Sprintf("%10s", game.players.Current().name), "> ")
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

		fmt.Println("===== playing phase ====")
		for game.state == GsPlaying {
			fmt.Print(fmt.Sprintf("%10s", game.players.Current().name), ": ")
			for i := 0; i < len(game.players.Current().hand); i++ {
				fmt.Print(cards.symbols[game.players.Current().hand[i]], ", ")
			}
			fmt.Print("\n> ")
			var card int
			_, err := fmt.Scanf("%d", &card)
			if err == nil {
				fmt.Println("Selected: ", cards.symbols[game.players.Current().hand[card]])
				if !game.PlayCardAt(cards, card) {
					fmt.Println("Wrong card!")
					continue
				} else {
					fmt.Print("Table: ")
					for i := 0; i < len(game.playedCards); i++ {
						fmt.Print(cards.symbols[game.playedCards[i].card], ", ")
					}
					fmt.Print("\n")
				}
			} else {
				continue
			}
		}

		if game.state == GsFinished {
			fmt.Println("===== scores =====")
			for i := 0; i < game.players.Len(); i++ {
				fmt.Print(fmt.Sprintf("%10s", game.players.At(i).name), ": ", game.players.At(i).points, "\n")
			}
			break
		}
	}
}
