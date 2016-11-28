package main

import (
	"testing"
	"fmt"
	"net"
	"os"
	"encoding/json"
	"strings"
	"strconv"
	"time"
	"math"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8081"
	CONN_TYPE = "tcp"
	LOADSAVE  = "LoadSave"
	NUM_GAMES = 1
	NUM_PLAYERS = 2
	WIDTH = 23
	HEIGHT = 18
	MAX_X = 704
	MAX_Y = 544
	BLOCK_SIZE = 32
)

var game [NUM_GAMES]Game

func main() {
	for i:=0; i < NUM_GAMES; i++ {
		game[i].make()
	}

	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error opening server")
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	go moveAndAttack()
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connections")
			os.Exit(1)
		}
		go handleRequest(connection)
	}
}
func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	text := make([]byte, 0, 4096)
	
	size, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error Reading: ", err.Error())
	}
	
	text = append(text, buf[:size]...)
	s := string(text)
	
	if (strings.Contains(s, ":")) {
		splits := strings.Split(s, ":")
		id, _ := strconv.Atoi(splits[0])
		
		if !isIdUsed(id) {
			for i:=0; i < NUM_GAMES; i++ {
				if game[i].needsPlayer() {
					game[i].addPlayer(id)
					break
				}
			}
		}
	
		if strings.Contains(s, ":game") {
			b, _ := json.Marshal(getPlayerById(id))
			conn.Write([]byte(b))
		} else if strings.Contains(s, ":buytower:") {
			towerid, _ := strconv.Atoi(splits[2])
			y, _ := strconv.Atoi(splits[3])
			x, _ := strconv.Atoi(splits[4])
			getPlayerById(id).buyTower(towerid, x, y, 0, 0)
		} else if strings.Contains(s, ":selltower:") {
			y, _ := strconv.Atoi(splits[2])
			x, _ := strconv.Atoi(splits[3])
			getPlayerById(id).sellTower(x, y)
		} else if strings.Contains(s, ":buyenemy:") {
			enemyid, _ := strconv.Atoi(splits[2])
			getPlayerById(id).buyEnemy(getPlayerEnemyById(id), enemyid)
		} else {
			conn.Write([]byte("hi"))
		}
	} else {
		conn.Write([]byte("hi"))
	}
	
	conn.Close()
}
func moveAndAttack() {
	var UPDATES, UPDATEUNIT int64
	UPDATES = 60
	UPDATEUNIT = 1000000
	executionStamp := time.Now().UnixNano() / UPDATEUNIT
	for { 
		now := time.Now().UnixNano() / UPDATEUNIT
		difference := now - executionStamp
		interval := 1000 / UPDATES
		if (difference > interval) {
			// DO WORK
			for i:=0;i<NUM_GAMES;i++ {
				game[i].Player[0].moveEnemies()
				game[i].Player[0].towersAttack()
				game[i].Player[1].moveEnemies()
				game[i].Player[1].towersAttack()
			}
			
			executionStamp = time.Now().UnixNano() / UPDATEUNIT
		}
	}
}

type Game struct {
	Player[NUM_PLAYERS] *Player
	Id[NUM_PLAYERS] int
}
func isIdUsed(id int) bool {
	for i:=0;i<NUM_GAMES;i++ {
		if game[i].Id[0] == id {
			return true
		} else if game[i].Id[1] == id {
			return true
		}
	}
	return false
}
func (g * Game) addPlayer(id int) {
	if g.Id[0] == 0 {
		g.Id[0] = id
	} else if g.Id[1] == 0 {
		g.Id[1] = id
	}
}
func (g * Game) needsPlayer() bool{
	if g.Id[0] == 0 {
		return true
	} else if g.Id[1] == 0 {
		return true
	}
	return false
}
func (g * Game) make() {
	var p1, p2 Player
	p1.create(0,20,100000,[]int{})
	p2.create(0, 20, 100000, []int{})
	g.Player[0] = &p1
	g.Player[1] = &p2
	g.Id[0] = 0
	g.Id[1] = 0
}
func getPlayerById(id int) * Player {
	for i:=0;i<NUM_GAMES;i++ {
		if game[i].Id[0] == id {
			return game[i].Player[0]
		} else if game[i].Id[1] == id {
			return game[i].Player[1]
		}
	}
	return nil
}
func getPlayerEnemyById(id int) * Player {
	for i:=0;i<NUM_GAMES;i++ {
		if game[i].Id[0] == id {
			return game[i].Player[1]
		} else if game[i].Id[1] == id {
			return game[i].Player[0]
		}
	}
	return nil
}

type Enemy struct {
	X int `json:"x"`
	Dx int
	Y int `json:"y"`
	Dy int
	Id int `json:"id"`
	Direction string
	Health int `json:"health"`
	Damage int `json:"damage"`
}
func (e * Enemy) create(id int) {
	e.X = 0
	e.Dx = 0
	e.Y = 0
	e.Dy = 1
	e.Direction = "DOWN"
	e.Id = id
	e.Health = 50
	e.Damage = 10
}
func (e * Enemy) damage(t Tower) {
	e.Health -= t.Damage
}
func (e Enemy) attack(t * Tower) {
	t.damage(e)
}
func (e Enemy) isDead() bool {
	return e.Health == 0
}

type Tower struct {
	Id int `json:"id"`
	Health int `json:"health"`
	Damage int `json:"damage"`
}
func (t * Tower) create(id int) {
	t.Id = id
	t.Health = (id-19)*100
	t.Damage = (id-19)*10
}
func (t * Tower) destroy() {
	t.Id = 0
	t.Health = 0
	t.Damage = 0
}
func (t * Tower) damage(e Enemy) {
	t.Health -= e.Damage
}
func (t Tower) attack(e * Enemy) {
	e.damage(t)
}
func (t Tower) isDead() bool {
	return (t.Health == 0)
}

type Board struct {
	X int `json:"x"`
	Y int `json:"y"`
	Building Tower `json:"building"`
}
func (b * Board) create(x int, y int) {
	b.X = x
	b.Y = y
}
func (b Board) towerExists() bool {
	return !b.Building.isDead()
}

type Player struct {
	Score int `json:"score"`
	Lives int `json:"lives"`
	Gold int `json:"gold"`
	Specials []int `json:"specials"`
	Field [HEIGHT][WIDTH]Board `json:"field"`
	Enemies []Enemy `json:"enemies"`
}
func (p * Player) towersAttack() {
	for y:=0;y<HEIGHT;y++ {
		for x:=0;x<WIDTH;x++ {
			tower := &p.Field[y][x]
			if tower.towerExists() {
				for i:=0;i<len(p.Enemies);i++ {
					enemy := p.Enemies[i]
					x1 := x
					x2 := enemy.X + (BLOCK_SIZE*11)  / BLOCK_SIZE
					y1 := y
					y2 := enemy.Y / BLOCK_SIZE
					if math.Sqrt(math.Pow((float64) (x2-x1),2)+math.Pow((float64) (y2-y1),2)) < 2 {
						if enemy.isDead() {
							p.killEnemy(enemy, true)
						} else {
							//tower.Building.attack(&p.Enemies[i])
						}
					}
				}
			}
		}
	}
}
func (p * Player) moveEnemies() {
	size := len(p.Enemies)
	
	if size != 0 {
		for i:= 0; i < size; i++ {
			e := &p.Enemies[i] 
			
			//Start moving down
			e.X += e.Dx
			e.Y += e.Dy
			
			// If at the end, finish the enemy
			if e.Y == MAX_Y && e.X == 0 {
				p.loseLife(p.Enemies[i])
				p.moveEnemies()
				break
			} else {
				// Get exact enemy location in-game
				eX := e.X + (BLOCK_SIZE*11) 
				eY := e.Y 
				
				// make sure enemy isnt off-screen
				if (eX < 0) {
					e.X -= e.Dx
					eX = 0
					e.Dy = 1
					e.Dx = 0
					e.Direction = "DOWN"
				}
				if (eX > MAX_X) {
					e.X -= e.Dx
					eX = MAX_X
					e.Dy = 1
					e.Dx = 0
					e.Direction = "DOWN"
				}
				if (eY < 0) {
					e.Y -= e.Dy
					eY = 0
					e.Dy = 0
					e.Dx = -1
					e.Direction = "LEFT"
				}
				if (eY > MAX_Y) {
					if (eX < MAX_X/2) {
						e.Y -= e.Dy
						eY = MAX_Y
						e.Dy = 0
						e.Dx = 1
						e.Direction = "RIGHT"
					} else {
						e.Y -= e.Dy
						eY = MAX_Y
						e.Dy = 0
						e.Dx = -1
						e.Direction = "LEFT"
					}
				}
				
				eMX := eX + BLOCK_SIZE 
				eMY := eY + BLOCK_SIZE 
				
				for y:=0;y<HEIGHT;y++ {
					for x:=0;x<WIDTH;x++ {
						if p.Field[y][x].towerExists() {
							// get exact tower location in-game
							tX := x*BLOCK_SIZE
							tY := y*BLOCK_SIZE
							tMX := tX + BLOCK_SIZE
							tMY := tY + BLOCK_SIZE
							
							if (eMY >= tY) && (eMY <= tMY) && (eMX-1 >= tX) && (eMX-1 <= tMX) && (e.Direction == "DOWN") {
								// Check down movement
								e.Y -= e.Dy
								e.Dy = -1
								e.Dx = 0
								e.Direction = "UP"
							} else if (eY >= tY) && (eY <= tMY) && (eMX >= tX) && (eMX <= tMX) && (e.Direction == "UP") {
								// check up movement
								e.Y -= e.Dy
								e.Dy = 0
								e.Dx = -1
								e.Direction = "LEFT"
							} else if (eX >= tX) && (eX <= tMX) && (eY >= tY) && (eY <= tMY) && (e.Direction == "LEFT") {
								// check left movement
								e.X -= e.Dx
								e.Dy = 0
								e.Dx = 1
								e.Direction = "RIGHT"
							} else if (eMX >= tX) && (eMX <= tMX) && (eY+1 >= tY) && (eY+1 <= tMY) && (e.Direction == "RIGHT") {
								// check right movement
								if eY >= MAX_Y-5 {
									e.X -= e.Dx
									e.Dy = -1
									e.Dx = 0
									e.Direction = "UP"
								} else {
									e.X -= e.Dx
									e.Dy = 1
									e.Dx = 0
									e.Direction = "DOWN"
								}
							}	
						}
					}
				}
			}
		}
	}
}
func (p * Player) create(pnts int, life int, gold int, specs []int) {
	p.Score = pnts
	p.Lives = life
	p.Gold = gold
	p.Specials = specs
	for y:=0;y<HEIGHT;y++ {
		for x:=0;x<WIDTH;x++ {
			p.Field[y][x].X = x*BLOCK_SIZE
			p.Field[y][x].Y = y*BLOCK_SIZE
		}
	}
}
func (p Player) won(p1 Player) bool {
	return p.Lives > p1.Lives
}
func (p Player) isDead() bool {
	return (p.Lives == 0)
}
func (p * Player) buyTower(id int, x int, y int, xt int, yt int) {
	towerCost := (id-19) * 50
	item := &p.Field[y][x]
	
	if !item.towerExists() {
		if p.Gold >= towerCost {
			item.create(xt, yt)
			item.Building.create(id)
		
			p.spendGold(towerCost)
		}
	}
}
func (p * Player) sellTower(x int, y int) {
	b := &p.Field[y][x]
	
	if b.towerExists() {
		towerCost := (b.Building.Id-19) * 25
		
		p.addGold(towerCost)
		
		b.Building.destroy()
	}
}
func (p * Player) addPoints(h int) {
	p.Score += h
}
func (p * Player) losePoints(h int) {
	p.Score -= h
	if p.Score < 0 { p.Score = 0 }
}
func (p * Player) loseLife(e Enemy) {
	p.Lives --
	
	p.losePoints(e.Id * 300)
	
	p.killEnemy(e, false)
	
	if p.Lives <= 0 { 
		print("Game WON!.........\n")
		os.Exit(1)
	}
}
func (p * Player) addGold(h int) {
	p.Gold += h
}
func (p * Player) spendGold(h int) {
	p.Gold -= h
	if p.Gold < 0 { p.Gold = 0 }
}
func (p * Player) addSpecial(h int) {
	p.Specials = append(p.Specials, h)
}
func (p * Player) useSpecial(h int) {
	p.Specials = append(p.Specials[:h], p.Specials[h+1:]...)
}
func (p * Player) buyEnemy(p1 * Player, id int) {
	enemyCost := id * 10
	
	if p.Gold >= enemyCost {
		var e Enemy
		e.create(id)
	
		p1.Enemies = append(p1.Enemies, e)
		
		p.spendGold(enemyCost)
		p.addPoints(enemyCost * 10)
	}
}
func (p * Player) killEnemy(e Enemy, pnts bool) {
	h := p.getEnemyId(e)

	if h < len(p.Enemies) && h >= 0 {
		if pnts {
			enemyCost := p.Enemies[h].Id * 20
	
			p.addGold(enemyCost)
			p.addPoints(enemyCost * 10)
		}
	
		p.Enemies = append(p.Enemies[:h], p.Enemies[h+1:]...)
	}
}
func (p * Player) getEnemyId(e Enemy) int {
	for i := 0; i < len(p.Enemies); i++ {
		if p.Enemies[i] == e {
			return i
		}
	}
	return 0
}

func assert(b bool, s string, t *testing.T) {
	if b {
		t.Error(s)
	}
}