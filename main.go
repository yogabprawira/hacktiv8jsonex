package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const StatusFile = "status.json"
const IntervalRefresh = 2

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type WStatus struct {
	Status Status `json:"status"`
}

type Updater struct {
	r       *rand.Rand
	wStatus WStatus
}

func (u *Updater) Init() {
	s := rand.NewSource(time.Now().UnixNano())
	u.r = rand.New(s)
}

func (u *Updater) Update() {
	u.wStatus.Status.Water = u.r.Intn(14) + 1
	u.wStatus.Status.Wind = u.r.Intn(19) + 1
	byteVal, err := json.Marshal(&u.wStatus)
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(StatusFile, byteVal, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func index(w http.ResponseWriter, _ *http.Request) {
	tmplt, err := template.New("index").ParseFiles("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	byteVal, err := ioutil.ReadFile(StatusFile)
	if err != nil {
		log.Fatalln(err)
	}
	var wStatus WStatus
	err = json.Unmarshal(byteVal, &wStatus)
	if err != nil {
		log.Fatalln(err)
	}
	err = tmplt.Execute(w, wStatus)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	updater := &Updater{}
	updater.Init()

	go func() {
		for {
			select {
			case <-time.After(IntervalRefresh * time.Second):
				updater.Update()
			}
		}
	}()

	http.HandleFunc("/", index)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
