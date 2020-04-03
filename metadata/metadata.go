package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type metadata struct {
	Leagues []string `json:"leagues"`
}

func Homepage() string {
	homepage, err := http.Get("http://www.espn.com/")
	if err != nil {
		fmt.Println(err)
		return "404: Error connecting to http://www.espn.com/"
	}
	defer homepage.Body.Close()
	body, err := ioutil.ReadAll(homepage.Body)
	if err != nil {
		return "400: Error reading response from http://www.espn.com/"
	}
	return string(body)
}

func GetSports() []string {
	file, err := ioutil.ReadFile("./metadata/data.json")
	if err != nil {
		fmt.Println(err)
		return []string{"400: error reading data.json file"}
	}
	data := metadata{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(err)
		return []string{"400: error unmarshalling data.json"}
	}
	fmt.Println(data)
	return data.Leagues
}
