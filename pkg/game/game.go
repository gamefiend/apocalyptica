// game manages data structures for games and playbooks arranged in format listed in docs/gamebook-playbook.md
package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gamefiend/apocalyptica/pkg/moves"
	"github.com/gamefiend/apocalyptica/pkg/playbook"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	gamepath   = fmt.Sprintf("%s/data/game", basepath)
)

//Gamebooks interface defines methods to initialize data structs from file structure.
type Gamebooks interface {
	NewGame() (g Games)
	ListGamesAvailable() (l string)
	LoadMovesFromGame(dir string) (m Move)
	LoadMovesFromPlaybook(path string) (m Move)
	GetCurrentGame(f string) (l string)
	SetCurrentGame(f string) bool
}

//Games type is where `games.json` structure is unmarshalled into.
type Games []struct {
	Name        string              `json:"Name"`
	Author      string              `json:"Author"`
	Description string              `json:"Description"`
	Playbooks   map[string]playbook `json:"Playboooks"`
	isSelected  bool                // NewGame sets default to false
}

//playbook type holds moves from a games' playbooks.
type playbook struct {
	Name string `json:"Name"`
	File string `json:"Moves"`
}

func (g *Games) NewGame() (gg Games) {
	var gg Games
	files, e := ioutil.ReadDir(gamepath)
	for _, x := range files {
		if x == "game.json" {
			f, e := ioutil.ReadFile(x)
			if e != nil {
				log.Fatalf("Cannot open file %s : %v\n", x, e)
			}
			r := bytes.NewReader(f)
			if e := json.NewDecoder(r).Decode(&g); e != nil {
				log.Fatalf("Cannot decode json: %v\n", e)
			}
		}
	}
	return

}

func (g *Games) ListGamesAvailable() (l string) {
	gl := make([]string, 0, len(g))
	for _, game := range g {
		output := fmt.Sprintf("**%s by %s**\n*%s*", game.Name, game.Author, game.Description)
		gl = append(gl, output)
	}
	return strings.Join(gl, "\n")
}

func (g *Games) LoadMovesFromGame(dir string) (m Move) {
	var m Move
	for _, p := range g {
		if p.isSelected {
			gpath := fmt.Sprintf("%s/%s", gamepath, dir)
			for _, x := range p.Playbooks {
				ppath := fmt.Sprintf("%s/%s", gpath, x.File)
				var pb Move
				r := ioutil.NewReader(ppath)
				if e := json.NewDecoder(r).Decode(&pb); e != nil {
					log.Fatalf("Cannot read playbook %s : %v\n", ppath, e)
				}
				m = append(m, pb)
			}

		}
	}
	return
}

func (g *Games) LoadMovesFromPlaybook(path string) (m Move) {
	var m Move
	f, e := ioutil.ReadFile(filepath)
	if e != nil {
		log.Fatalf("Cannot read file %s : %v\n", filepath, e)
	}
	r := bytes.NewReader(f)
	if e := json.NewDecoder(r).Decode(&m); e != nil {
		log.Fatalf("Problem unmarshalling json : %v", e)
	}
	return
}

func (g *Games) GetCurrentGame() string {
	for _, x := range g {
		if x.isSelected {
			return x.Name
		}
	}
	//if we have arrived here, we obviously have nothing selected...
	return "NoGameSelected"
}

func (g *Games) SetCurrentGame(f string) bool {
	for _, x := range g {
		if x.Name == f {
			if x.isSelected {
				return true
			}
			x.isSelected = true
			return true
		}
		x.isSelected = false
	}
	return false
}
