package main

type Player struct {
	name          string
	points        uint
	estimatedWins uint
	wins          uint
	hand          []CardId
}

func (p *Player) CalcRoundScore() {
	if p.wins == p.estimatedWins {
		if p.estimatedWins == 0 {
			p.points += 20
		} else {
			p.points += p.estimatedWins + 10
		}
	}
	p.estimatedWins = 0
	p.wins = 0
}

func (p *Player) Win() {
	p.wins++
}

func (p *Player) HasEmptyHand() bool {
	return len(p.hand) == 0
}

type Players struct {
	active  int
	dealer  int
	players []*Player
}

func NewPlayers(players []string) *Players {
	pp := &Players{
		players: make([]*Player, len(players)),
		active:  0,
		dealer:  0,
	}
	for i := 0; i < pp.Len(); i++ {
		p := &Player{
			name:          players[i],
			points:        0,
			estimatedWins: 0,
			wins:          0,
		}
		pp.players[i] = p
	}
	return pp
}

func (p *Players) Win(name string) {
	for i := 0; i < p.Len(); i++ {
		if p.players[i].name == name {
			p.active = i
			p.players[i].Win()
		}
	}
}

func (p *Players) CalcRoundScores() {
	for i := 0; i < p.Len(); i++ {
		p.players[i].CalcRoundScore()
	}
	p.dealer++
	if p.dealer >= p.Len() {
		p.dealer = 0
	}
}

func (p *Players) Len() int {
	return len(p.players)
}

func (p *Players) Next() {
	p.active++
	if p.active >= p.Len() {
		p.active = 0
	}
}

func (p *Players) Current() *Player {
	return p.players[p.active]
}

func (p *Players) At(index int) *Player {
	return p.players[index]
}
