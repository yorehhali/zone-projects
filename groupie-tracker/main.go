package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Page struct {
	Title    string
	Artists  *[]Artist
	RandomId int
}

type ProfilePage struct {
	ArtistId           int
	Artist             Artist
	Albums             []AlbumsApiData
	ArtistApi          SearchApiData
	RandomId           int
	Locations          []string
	DatesLocations     map[string][]string
	LatsLongs          map[string]map[string]float64
	BeautifulLocations map[string]string
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

type SearchApi struct {
	Data  []SearchApiData `json:"data"`
	Total int             `json:"total"`
	Next  string          `json:"next"`
}

type SearchApiData struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	PictureXl string `json:"picture_xl"`
	NbAlbum   int    `json:"nb_album"`
	NbFan     int    `json:"nb_fan"`
}

type AlbumsApiData struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover_medium"`
	Fans        int    `json:"fans"`
	ReleaseDate string `json:"release_date"`
	RecordType  string `json:"record_type"`
	Tracklist   string `json:"tracklist"`
}

type AlbumsApi struct {
	Data  []AlbumsApiData `json:"data"`
	Total int             `json:"total"`
}

type LatLongApi struct {
	Results []LatLongApiResult
}

type LatLongApiResult struct {
	Locations []LatLongApiResultLocation
}

type LatLongApiResultLocation struct {
	LatLng map[string]float64
}

var artist []Artist
var p = Page{
	Title:    "Groupie Tracker",
	Artists:  &artist,
	RandomId: GenRandomId(),
}

func main() {
	artist = GetArtists()
	fs := http.FileServer(http.Dir("templates"))
	router := http.NewServeMux()
	fmt.Println("Starting server on port 8080")
	router.HandleFunc("/", HandlerIndex)
	router.HandleFunc("/homepage", HandlerHomepage)
	router.HandleFunc("/profile", HandlerProfile)
	router.HandleFunc("/profiledates", HandlerProfiledates)
	router.Handle("/templates/", http.StripPrefix("/templates/", fs))
	http.ListenAndServe(":8080", router)
}

func GetArtists() []Artist {
	var artistData []Artist
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	json.Unmarshal([]byte(sb), &artistData)
	return artistData
}

func GetLocations(url string) Location {
	var artistLocation Location
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	json.Unmarshal([]byte(sb), &artistLocation)
	return artistLocation
}

func GetRelation(url string) Relation {
	var relationData Relation
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	json.Unmarshal([]byte(sb), &relationData)
	return relationData
}

func GetArtistApi(name string) SearchApiData {
	var dataApi SearchApi
	url := "https://api.deezer.com/search/artist/?q=" + strings.Replace(name, " ", "%20", 10) + "&index=0&limit=1"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	json.Unmarshal([]byte(sb), &dataApi)
	return dataApi.Data[0]
}
func GetArtistAlbums(artistId int) []AlbumsApiData {
	var dataApi AlbumsApi
	url := "https://api.deezer.com/artist/" + strconv.Itoa(artistId) + "/albums&limit=500"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	json.Unmarshal([]byte(sb), &dataApi)
	return dataApi.Data
}

func FilterAlbums(albums []AlbumsApiData) []AlbumsApiData {
	var newAlbums []AlbumsApiData
	for _, album := range albums {
		if album.RecordType == "album" {
			newAlbums = append(newAlbums, album)
		}
	}
	return RemoveDuplicatesAlbumsApi(newAlbums)
}

func RemoveDuplicatesAlbumsApi(albums []AlbumsApiData) []AlbumsApiData {
	var newAlbums []AlbumsApiData
	for _, album := range albums {
		if len(newAlbums) == 0 {
			newAlbums = append(newAlbums, album)
		} else {
			for i, newAlbum := range newAlbums {
				if strings.EqualFold(album.Title, newAlbum.Title) {
					break
				}
				if i == len(newAlbums)-1 {
					newAlbums = append(newAlbums, album)
				}
			}
		}
	}
	return newAlbums
}

func removeDuplicateStr(strSlice []string) []string {
	var list []string
	allKeys := make(map[string]bool)
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GenRandomId() int {
	artist = GetArtists()
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(len(artist))
}

func GetLatLongApi(location string) map[string]float64 {
	var dataApi LatLongApi

	url := "https://open.mapquestapi.com/geocoding/v1/address?key=37GzZAcEPu9TQdvGkZ3DREYAPaLVNZBC&location=" + strings.Replace(location, " ", "%20", 100) + "&thumbMaps=false&maxResults=1"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	json.Unmarshal([]byte(sb), &dataApi)
	return dataApi.Results[0].Locations[0].LatLng
}

func GetLatLong(cities []string, citiesMap map[string]string) map[string]map[string]float64 {
	m := make(map[string]map[string]float64)

	for _, city := range cities {
		m[city] = GetLatLongApi(citiesMap[city])
	}
	return m
}

func BeautifyLocation(location string) string {
	test := strings.Replace(location, "-", ", ", 100)
	test = strings.Replace(test, "_", "-", 100)
	return strings.ToUpper(test)
}

func BeautifyLocations(locations []string) map[string]string {
	newLocations := make(map[string]string)
	for _, location := range locations {
		newLocations[location] = BeautifyLocation(location)
	}
	return newLocations
}

func HandlerHomepage(w http.ResponseWriter, r *http.Request) {
	p.RandomId = GenRandomId()
	t, _ := template.ParseGlob("templates/*.html")
	t.ExecuteTemplate(w, "homepage.html", p)
}

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseGlob("templates/*.html")
	t.ExecuteTemplate(w, "index.html", p)
}

func HandlerProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		artistIdString := r.FormValue("id")
		artistId, _ := strconv.Atoi(artistIdString)
		artistApi := GetArtistApi(artist[artistId-1].Name)
		pProfile := ProfilePage{
			ArtistId:  artistId,
			Artist:    artist[artistId-1],
			Albums:    FilterAlbums(GetArtistAlbums(artistApi.Id)),
			ArtistApi: artistApi,
			RandomId:  GenRandomId(),
		}
		t, _ := template.ParseGlob("templates/*.html")
		t.ExecuteTemplate(w, "profile.html", pProfile)
	}
}

func HandlerProfiledates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		artistIdString := r.FormValue("id")
		artistId, _ := strconv.Atoi(artistIdString)
		artistApi := GetArtistApi(artist[artistId-1].Name)

		datesLocations := GetRelation(artist[artistId-1].Relations).DatesLocation
		locations := GetLocations(artist[artistId-1].Locations).Locations
		locations = removeDuplicateStr(locations)
		beautifulLocations := BeautifyLocations(locations)

		pProfile := ProfilePage{
			ArtistId:           artistId,
			Artist:             artist[artistId-1],
			Albums:             FilterAlbums(GetArtistAlbums(artistApi.Id)),
			ArtistApi:          artistApi,
			RandomId:           GenRandomId(),
			Locations:          locations,
			DatesLocations:     datesLocations,
			LatsLongs:          GetLatLong(locations, beautifulLocations),
			BeautifulLocations: beautifulLocations,
		}
		t, _ := template.ParseGlob("templates/*.html")
		t.ExecuteTemplate(w, "profiledates.html", pProfile)
	}
}