package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

type User struct {
	User     string `json:"user"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type Config struct {
	Users     []User
	Port      int
	SkipLogin string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateTicket() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func login(users []User, username string, password string) *User {
	for i := range users {
		if users[i].Username == username && users[i].Password == password {
			return &users[i]
		}
	}

	return nil
}

func main() {
	confFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln("please provide config.yml")
	}
	config := Config{}
	err = yaml.Unmarshal(confFile, &config)
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}
	var defaultUser *User
	if config.SkipLogin != "" {
		for i := range config.Users {
			if config.Users[i].User == config.SkipLogin {
				defaultUser = &config.Users[i]
			}
		}
	}

	tickets := map[string]*User{}

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		service := r.URL.Query().Get("service")
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")
		user := login(config.Users, username, password)
		if user == nil && defaultUser != nil {
			user = defaultUser
		}
		if user == nil {
			tmpl, err := template.New("login.tmpl").ParseFiles("login.tmpl")
			if err != nil {
				panic(err)
			}
			err = tmpl.Execute(w, service)
			if err != nil {
				panic(err)
			}
			return
		}
		serviceUrl, _ := url.Parse(service)
		q := serviceUrl.Query()
		ticket := generateTicket()
		tickets[ticket] = user
		q.Set("ticket", ticket)
		serviceUrl.RawQuery = q.Encode()
		http.Redirect(w, r, serviceUrl.String(), http.StatusSeeOther)
	})

	http.HandleFunc("/serviceValidate", func(w http.ResponseWriter, r *http.Request) {
		ticket := r.URL.Query().Get("ticket")
		w.Header().Set("Content-Type", "application/json")
		if tickets[ticket] != nil {
			json.NewEncoder(w).Encode(*tickets[ticket])
			tickets[ticket] = nil
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	err = http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
