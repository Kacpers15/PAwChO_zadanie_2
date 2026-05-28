package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var mapaRegionow = map[string]map[string]string{
	"Mazowieckie":   {"Warszawa": "warszawa", "Płock": "plock"},
	"Małopolskie":   {"Kraków": "krakow", "Zakopane": "zakopane", "Nowy Sącz": "nowysacz"},
	"Dolnośląskie":  {"Wrocław": "wroclaw", "Kłodzko": "klodzko"},
	"Pomorskie":     {"Gdańsk": "gdansk", "Łeba": "leba", "Hel": "hel"},
	"Wielkopolskie": {"Poznań": "poznan", "Kalisz": "kalisz"},
	"Lubelskie":     {"Lublin": "lublin", "Terespol": "terespol"},
	"Podkarpackie":  {"Rzeszów": "rzeszow", "Krosno": "krosno"},
	"Podlaskie":     {"Białystok": "bialystok", "Suwałki": "suwalki"},
	"Opolskie":      {"Opole": "opole"},
	"Łódzkie":       {"Łódź": "lodz"},
}

var listaMiast = make(map[string]string)

func init() {
	for _, miasta := range mapaRegionow {
		for nazwaMiasta, idMiasta := range miasta {
			listaMiast[nazwaMiasta] = idMiasta
		}
	}
}

const portSerwera = "8080"
const autorProgramu = "Kacper Sumera"

type OdpowiedzImgw struct {
	Stacja             string `json:"stacja"`
	Temperatura        string `json:"temperatura"`
	PredkoscWiatru     string `json:"predkosc_wiatru"`
	WilgotnoscWzgledna string `json:"wilgotnosc_wzgledna"`
	SumaOpadu          string `json:"suma_opadu"`
	Cisnienie          string `json:"cisnienie"`
	DataPomiaru        string `json:"data_pomiaru"`
}

var szablonStrony *template.Template

func main() {
	var blad error
	szablonStrony, blad = template.ParseFiles("index.html")
	if blad != nil {
		log.Fatal("Nie można załadować pliku index.html: ", blad)
	}

	log.Printf("=== START APLIKACJI ===")
	log.Printf("Data uruchomienia: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	log.Printf("Autor programu: %s\n", autorProgramu)
	log.Printf("Aplikacja nasłuchuje na porcie TCP: %s\n", portSerwera)
	log.Printf("=======================\n")

	http.HandleFunc("/", obslugaStronyGlownej)
	http.HandleFunc("/weather", obslugaPogody)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Fatal(http.ListenAndServe(":"+portSerwera, nil))
}

func obslugaStronyGlownej(w http.ResponseWriter, r *http.Request) {
	regionyJSON, _ := json.Marshal(mapaRegionow)

	szablonStrony.Execute(w, map[string]interface{}{
		"RegionyJSON": template.JS(string(regionyJSON)),
		"Wybrane":     "",
	})
}

func obslugaPogody(w http.ResponseWriter, r *http.Request) {
	wybraneMiasto := r.URL.Query().Get("location")
	idStacji, istnieje := listaMiast[wybraneMiasto]

	if !istnieje {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	adresApi := fmt.Sprintf("https://danepubliczne.imgw.pl/api/data/synop/station/%s", idStacji)
	odpowiedz, blad := http.Get(adresApi)
	if blad != nil {
		http.Error(w, "Błąd podczas pobierania danych z IMGW", http.StatusInternalServerError)
		return
	}
	defer odpowiedz.Body.Close()

	var pogoda OdpowiedzImgw
	blad = json.NewDecoder(odpowiedz.Body).Decode(&pogoda)
	if blad != nil {
		http.Error(w, "Błąd parsowania danych pogodowych", http.StatusInternalServerError)
		return
	}

	regionyJSON, _ := json.Marshal(mapaRegionow)
	szablonStrony.Execute(w, map[string]interface{}{
		"RegionyJSON":  template.JS(string(regionyJSON)),
		"Wybrane":      wybraneMiasto,
		"DanePogodowe": pogoda,
	})
}
