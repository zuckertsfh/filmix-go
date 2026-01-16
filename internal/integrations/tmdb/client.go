package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURL   = "https://api.themoviedb.org/3"
	ImageBase = "https://image.tmdb.org/t/p/original"
)

type Client struct {
	ApiKey string
	Client *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey: apiKey,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

type TMDBMovieResponse struct {
	Results []TMDBMovie `json:"results"`
}

type TMDBMovie struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	ReleaseDate  string  `json:"release_date"`
	VoteAverage  float64 `json:"vote_average"`
	Popularity   float64 `json:"popularity"`
	GenreIDs     []int   `json:"genre_ids"`
}

type TMDBGenreResponse struct {
	Genres []TMDBGenre `json:"genres"`
}

type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetGenres() (map[int]string, error) {
	url := fmt.Sprintf("%s/genre/movie/list?api_key=%s&language=en-US", BaseURL, c.ApiKey)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch genres: %s", resp.Status)
	}

	var genreResp TMDBGenreResponse
	if err := json.NewDecoder(resp.Body).Decode(&genreResp); err != nil {
		return nil, err
	}

	genres := make(map[int]string)
	for _, g := range genreResp.Genres {
		genres[g.ID] = g.Name
	}
	return genres, nil
}

func (c *Client) FetchNowPlayingRaw() ([]TMDBMovie, error) {
	url := fmt.Sprintf("%s/movie/now_playing?api_key=%s&language=en-US&page=1", BaseURL, c.ApiKey)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch movies: %s", resp.Status)
	}

	var movieResp TMDBMovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&movieResp); err != nil {
		return nil, err
	}

	return movieResp.Results, nil
}
