package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	User string
	Port int
}

type Response struct {
	User string `json:"user"`
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

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		service := r.URL.Query().Get("service")
		serviceUrl, _ := url.Parse(service)
		q := serviceUrl.Query()
		q.Set("ticket", "123")
		serviceUrl.RawQuery = q.Encode()
		http.Redirect(w, r, serviceUrl.String(), http.StatusSeeOther)
	})

	http.HandleFunc("/serviceValidate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := Response{
			User: config.User,
		}

		json.NewEncoder(w).Encode(data)
	})

	err = http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
