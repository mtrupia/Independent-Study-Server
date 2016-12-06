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
	"Server/utils"
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
	go moveAll()
	go attackAll()
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
func moveAll() {
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
				go game[i].Player[0].moveEnemies()
				go game[i].Player[1].moveEnemies()
			}
			
			executionStamp = time.Now().UnixNano() / UPDATEUNIT
		}
	}
}
func attackAll() {
	var UPDATES, UPDATEUNIT int64
	UPDATES = 1
	UPDATEUNIT = 1000000
	executionStamp := time.Now().UnixNano() / UPDATEUNIT
	for { 
		now := time.Now().UnixNano() / UPDATEUNIT
		difference := now - executionStamp
		interval := 1000 / UPDATES
		if (difference > interval) {
			for i:=0;i<NUM_GAMES;i++ {
				go game[i].Player[0].towersAttack()
				go game[i].Player[1].towersAttack()
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
	DX int
	Y int `json:"y"`
	DY int
	Id int `json:"id"`
	Health int `json:"health"`
	Damage int `json:"damage"`
}
func (e * Enemy) create(id int) {
	e.X = 0
	e.DX = 0
	e.Y = 0
	e.DY = 0
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
func (e * Enemy) isDead() bool {
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
	Path	utils.Point
}
func (p * Player) towersAttack() {
	for y:=0;y<HEIGHT;y++ {
		for x:=0;x<WIDTH;x++ {
			tower := &p.Field[y][x]
			if tower.towerExists() {
				for i:=0;i<len(p.Enemies);i++ {
					enemy := &p.Enemies[i]
					x1 := x*BLOCK_SIZE + BLOCK_SIZE/2
					x2 := enemy.X + (BLOCK_SIZE*11) + BLOCK_SIZE/2
					y1 := y*BLOCK_SIZE + BLOCK_SIZE/2
					y2 := enemy.Y + BLOCK_SIZE/2
					if math.Sqrt(math.Pow((float64) (x2-x1),2)+math.Pow((float64) (y2-y1),2)) < BLOCK_SIZE*2 {
						tower.Building.attack(enemy)
						if enemy.isDead() {
							p.killEnemy(p.Enemies[i], true)
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
			
			//Start moving
			e.X += e.DX
			e.Y += e.DY
			
			// If at the end, finish the enemy, else, update dx, dy
			if e.Y == MAX_Y && e.X == 0 {
				p.loseLife(p.Enemies[i])
				break
			} else {
				if ( ((e.Y % 32) == 0 && e.DY != 0) || ((e.X % 32) == 0 && e.DX != 0 ) ) {
					p.getEnemyDirection(e)
				}
			}
			
			
		}
	}
}
func (p Player) getEnemyDirection(e * Enemy) {
	e.DX = 0
	e.DY = 0
	
	// get actual enemy location for board x and y
	enemyX := e.X + (BLOCK_SIZE*11)
	enemyY := e.Y
	enemyX /= 32
	enemyY /= 32
	
	coords := p.Path
	for coords.X-1 != enemyY || coords.Y-1 != enemyX {
		coords = *coords.Parent
	}
	coords = *coords.Parent
	
	//print("[",coords.Y-1,":",coords.X-1,"], [",enemyX,":",enemyY,"]\n")
	if enemyY < coords.X-1 {
		e.DY = 1
	} else if enemyY > coords.X-1 {
		e.DY = -1
	} else if enemyX < coords.Y-1 {
		e.DX = 1
	} else if enemyX > coords.Y-1 {
		e.DX = -1
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
	p.getPath()
}
func (p * Player) getPath(){
	var scene Scene
	scene.initScene(p)
	initAstar(&scene)
	p.Path = findPath(&scene)
}
func parsePath(p utils.Point) {
	// print path
	print("[",p.X-1,":",p.Y-1,"], ")
	if p.Parent != nil {
		parsePath(*(p.Parent))
	} else {
		print("\n")
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
			
			p.getPath()
		}
	}
}
func (p * Player) sellTower(x int, y int) {
	b := &p.Field[y][x]
	
	if b.towerExists() {
		towerCost := (b.Building.Id-19) * 25
		
		p.addGold(towerCost)
		
		b.Building.destroy()
		
		p.getPath()
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
		p1.getEnemyDirection(&e)
	
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