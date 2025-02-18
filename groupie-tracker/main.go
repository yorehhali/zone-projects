package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Page struct {
	Title   string
	Artists *[]Artist
}

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Relation struct {
	Id            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type ProfilePage struct {
	ArtistId           int
	Artist             Artist
	Locations          []string
	DatesLocations     map[string][]string
	BeautifulLocations map[string]string
}

var artists []Artist

var mainPage = Page{
	Title:   "Groupie Tracker",
	Artists: &artists,
}

func main() {
	artists = GetArtists()
	fs := http.FileServer(http.Dir("templates"))
	fmt.Println("Starting server on localhost 9090")
	http.HandleFunc("/", HandlerIndex)
	http.HandleFunc("/homepage", HandlerHomepage)
	http.HandleFunc("/profile", HandlerProfile)
	http.HandleFunc("/profiledates", HandlerProfiledates)
	http.Handle("/templates/", http.StripPrefix("/templates/", fs))
	http.ListenAndServe(":9090", nil)
}
func GetArtists() []Artist {
	var artistData []Artist
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatalln("Error fetching artists:", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err)
	}
	json.Unmarshal(body, &artistData)
	return artistData
}

func GetLocations(url string) Location {
	var artistLocation Location
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("Error fetching locations:", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err)
	}
	json.Unmarshal(body, &artistLocation)
	return artistLocation
}
func GetRelation(url string) Relation {
	var relationData Relation
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("Error fetching relations:", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err)
	}
	json.Unmarshal(body, &relationData)
	return relationData
}
func BeautifyLocation(location string) string {
	location = strings.Replace(location, "-", ", ", -1)
	location = strings.Replace(location, "_", "-", -1)
	return strings.ToUpper(location)
}
func BeautifyLocations(locations []string) map[string]string {
	formattedLocations := make(map[string]string)
	for _, location := range locations {
		formattedLocations[location] = BeautifyLocation(location)
	}
	return formattedLocations
}
func searchArtists(query string) []Artist {
	var filteredArtists []Artist
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			filteredArtists = append(filteredArtists, artist)
		}
	}
	return filteredArtists
}
func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	var filteredArtists []Artist
	if query != "" {
		filteredArtists = searchArtists(query)
	} else {
		filteredArtists = artists
	}
	mainPage := Page{
		Title:   "Groupie Tracker",
		Artists: &filteredArtists,
	}
	t, _ := template.ParseGlob("templates/*.html")
	t.ExecuteTemplate(w, "index.html", struct {
		Page
		SearchQuery string
	}{
		Page:        mainPage,
		SearchQuery: query,
	})
}
func HandlerHomepage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseGlob("templates/*.html")
	t.ExecuteTemplate(w, "homepage.html", mainPage)
}
func HandlerProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Error parsing form: %v", err)
			return
		}
		artistIdString := r.FormValue("id")
		artistId, _ := strconv.Atoi(artistIdString)
		profilePage := ProfilePage{
			ArtistId: artistId,
			Artist:   artists[artistId-1],
		}
		t, _ := template.ParseGlob("templates/*.html")
		t.ExecuteTemplate(w, "profile.html", profilePage)
	}
}
func HandlerProfiledates(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Error parsing form: %v", err)
			return
		}
		artistIdString := r.FormValue("id")
		artistId, _ := strconv.Atoi(artistIdString)
		datesLocations := GetRelation(artists[artistId-1].Relations).DatesLocation
		locations := GetLocations(artists[artistId-1].Locations).Locations
		beautifulLocations := BeautifyLocations(locations)
		profilePage := ProfilePage{
			ArtistId:           artistId,
			Artist:             artists[artistId-1],
			Locations:          locations,
			DatesLocations:     datesLocations,
			BeautifulLocations: beautifulLocations,
		}
		t, _ := template.ParseGlob("templates/*.html")
		t.ExecuteTemplate(w, "profiledates.html", profilePage)
	}
}
