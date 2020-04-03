package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type AthleteStat struct {
	Athlete struct {
		ID          string `json:"id"`
		UID         string `json:"uid"`
		GUID        string `json:"guid"`
		Type        string `json:"type"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		DisplayName string `json:"displayName"`
		ShortName   string `json:"shortName"`
		DebutYear   int    `json:"debutYear"`
		Links       []struct {
			Language   string   `json:"language"`
			Rel        []string `json:"rel"`
			Href       string   `json:"href"`
			Text       string   `json:"text"`
			ShortText  string   `json:"shortText"`
			IsExternal bool     `json:"isExternal"`
			IsPremium  bool     `json:"isPremium"`
		} `json:"links"`
		Headshot struct {
			Href string `json:"href"`
			Alt  string `json:"alt"`
		} `json:"headshot"`
		Position struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			DisplayName  string `json:"displayName"`
			Abbreviation string `json:"abbreviation"`
			Leaf         bool   `json:"leaf"`
			Parent       struct {
				Leaf bool `json:"leaf"`
			} `json:"parent"`
			Slug string `json:"slug"`
		} `json:"position"`
		Status struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Type         string `json:"type"`
			Abbreviation string `json:"abbreviation"`
		} `json:"status"`
		Age           int    `json:"age"`
		TeamName      string `json:"teamName"`
		TeamShortName string `json:"teamShortName"`
		Teams         []struct {
			Name         string `json:"name"`
			Abbreviation string `json:"abbreviation"`
		} `json:"teams"`
		Slug    string `json:"slug"`
		TeamID  string `json:"teamId"`
		TeamUID string `json:"teamUId"`
	} `json: "athlete"`
	Categories []struct {
		Name        string    `json:"name"`
		DisplayName string    `json:"displayName"`
		Totals      []string  `json:"totals"`
		Values      []float64 `json:"values"`
		Ranks       []string  `json:"ranks"`
	} `json:"categories"`
}

type ByAthleteResponse struct {
	Pagination struct {
		Count int    `json:"count"`
		Limit int    `json:"limit"`
		Page  int    `json:"page"`
		Pages int    `json:"pages"`
		First string `json:"first"`
		Next  string `json:"next"`
		Last  string `json:"last"`
	} `json:"pagination"`
	Athletes      []AthleteStat `json:"athletes"`
	CurrentSeason struct {
		Year        int    `json:"year"`
		DisplayName string `json:"displayName"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
		Type        struct {
			ID        string `json:"id"`
			Type      int    `json:"type"`
			Name      string `json:"name"`
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Week      struct {
				Number    int    `json:"number"`
				StartDate string `json:"startDate"`
				EndDate   string `json:"endDate"`
				Text      string `json:"text"`
			} `json:"week"`
		} `json:"type"`
	} `json:"currentSeason"`
	RequestedSeason struct {
		Year        int    `json:"year"`
		DisplayName string `json:"displayName"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
		Type        struct {
			ID        string `json:"id"`
			Type      int    `json:"type"`
			Name      string `json:"name"`
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Week      struct {
				Number    int    `json:"number"`
				StartDate string `json:"startDate"`
				EndDate   string `json:"endDate"`
				Text      string `json:"text"`
			} `json:"week"`
		} `json:"type"`
	} `json:"requestedSeason"`
	Glossary []struct {
		Abbreviation string `json:"abbreviation"`
		DisplayName  string `json:"displayName"`
	} `json:"glossary"`
	Categories []struct {
		Name         string   `json:"name"`
		DisplayName  string   `json:"displayName"`
		Labels       []string `json:"labels"`
		Names        []string `json:"names"`
		DisplayNames []string `json:"displayNames"`
		Descriptions []string `json:"descriptions"`
	} `json:"categories"`
}

func Update() {
	baseURL := "https://site.web.api.espn.com/apis/common/v3/sports/basketball/nba/statistics/byathlete"
	dataIn1, err := http.Get(baseURL)
	if err != nil {
		fmt.Println(err)
		log.Fatal("404: Error connecting to espn api")
	}
	body1, err := ioutil.ReadAll(dataIn1.Body)
	var response1 ByAthleteResponse
	if err != nil {
		log.Fatal("400: Error reading response from espn api")
	}
	json.Unmarshal(body1, &response1)
	dataIn1.Body.Close()
	athletes := response1.Athletes
	numPages := response1.Pagination.Pages
	var response ByAthleteResponse
	json.Unmarshal(body1, &response)
	for i := 1; i <= numPages; i++ {
		dataIn, err := http.Get(response.Pagination.Next)
		if err != nil {
			log.Fatal("400: Error connection to espn api")
		}
		body, err := ioutil.ReadAll(dataIn.Body)
		if err != nil {
			log.Fatal("400: Error reading response from espn api")
		}
		json.Unmarshal(body, &response)
		dataIn.Body.Close()
		temp := athletes
		athletes = append(temp, response.Athletes...)
	}
	dataOut, err := json.Marshal(athletes)
	if err != nil {
		log.Fatal("400: Error Marshalling nba athletes data")
	}
	ioutil.WriteFile("./data/nba-athletes-stats.json", dataOut, 0660)
}
