package main

import (
	"errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"log"
)

const URL = "http://www.omdbapi.com/?t="

type Movie struct {
	Title    string
	Country  string
	Poster   string
	Response string
	Error    string
}

func searchMovie(keywords string) (*Movie, error) {
	resp, err := http.Get(URL + keywords + "&apikey=" + os.Getenv("OMDB_API_KEY"))

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var movie Movie
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		resp.Body.Close()
		return nil, err
	}
	if movie.Response != "True" {
		return nil, errors.New(movie.Error)
	}
	if movie.Poster == "" {
		log.Fatalln("ポスター画像がありません")
	}

	return &movie, err
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("映画の名前を入力してください")
	}

	keywords := url.QueryEscape(strings.Join(os.Args[1:], " "))
	movie, err := searchMovie(keywords)

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(movie.Poster)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatal("search query failed: %", resp.Status)
	}

	file, err := os.Create(movie.Title + ".jpg")
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		log.Fatal(err)
	}
}