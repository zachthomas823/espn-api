package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/espn-api/metadata"
	"github.com/espn-api/update"
)

func main() {
	update.Update()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, metadata.Homepage())
	})
	http.HandleFunc("/sports", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, metadata.GetSports())
	})
	http.HandleFunc("/nba/players", func(w http.ResponseWriter, r *http.Request) {
		keys := r.URL.Query()
		data, err := ioutil.ReadFile("./data/nba-athletes-stats.json")
		if err != nil {
			fmt.Fprint(w, "400: Cannot find nba player data")
		}
		var athletes []update.AthleteStat
		err = json.Unmarshal(data, &athletes)
		if err != nil {
			fmt.Fprint(w, err)
		}
		if keys.Get("firstname") != "" {
			firstName := strings.ToLower(keys.Get("firstname"))
			newAthletes := []update.AthleteStat{}
			for _, v := range athletes {
				if strings.ToLower(v.Athlete.FirstName) == firstName {
					fmt.Println("found " + firstName)
					newAthletes = append(newAthletes, v)
				}
			}
			athletes = newAthletes
		}
		if keys.Get("lastname") != "" {
			firstName := strings.ToLower(keys.Get("lastname"))
			newAthletes := []update.AthleteStat{}
			for _, v := range athletes {
				if strings.ToLower(v.Athlete.LastName) == firstName {
					newAthletes = append(newAthletes, v)
				}
			}
			athletes = newAthletes
		}
		response, err := json.Marshal(athletes)
		if err != nil {
			fmt.Fprint(w, "400: Error marshalling into json")
		} else {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(response))
		}
	})
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
