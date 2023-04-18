package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

const port = ":8080"

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Notfound(w, r)
		return
	}
	var s []string
	renderTemplate(w, "home")
	ascii := r.FormValue("textarea")
	font := r.FormValue("font")
	if font == "Standard" {
		dat, err1 := os.Open("ascii.txt")
		if err1 != nil {
			return
		}
		reader := bufio.NewScanner(dat)
		for reader.Scan() {
			s = append(s, reader.Text())
		}
	} else if font == "Shadow" {
		dat, err1 := os.Open("shadow.txt")
		if err1 != nil {
			return
		}
		reader := bufio.NewScanner(dat)
		for reader.Scan() {
			s = append(s, reader.Text())
		}
	} else if font == "Thinkertoy" {
		dat, err1 := os.Open("thinkertoy.txt")
		if err1 != nil {
			return
		}
		reader := bufio.NewScanner(dat)
		for reader.Scan() {
			s = append(s, reader.Text())
		}
	}
	fmt.Fprintf(w, "<pre id=\"my-ascii\">")
	for _, c := range ascii {
		if (c < 32 || c >= 127) && c != '\n' && c != '\r' {
			fmt.Fprintf(w, "this caracter is not available.")
			fmt.Fprintf(w, "</pre>")
			return
		}
	}

	replace := strings.ReplaceAll(ascii, "\\n", "\n")
	jump := strings.ReplaceAll(replace, string([]byte{0x0D, 0x0A}), "\n")
	content := strings.Split(jump, "\n")
	// Eclatement de l'argument pour chaque occurrence des caractères \ et n
	// Library
	file, err := os.Open("ascii.txt")
	if err != nil {
		fmt.Println("Le fichier n'existe pas.")
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	f, err := os.Create("download.txt")
	if err != nil {
		fmt.Println("le fichier n'as pas été créé")
	}

	// Lecture du tableau de lignes données dans la commande
	for _, element := range content {
		// Dans le cas où la cellule lue n'est pas vide
		if len(content) > 0 {
			line := []rune(element) // Eclatement du contenu de la cellule en tableau de rune
			// Boucle pour les 8 lignes de l'ascii art
			for a := 0; a < 8; a++ {
				// LEcture du tableau de rune
				for i := 0; i < len(line); i++ {
					group := (int(line[i]) - 32) * 9 // Définition de la première ligne dédiée à l'ascii art correspondant à la rune
					adress := group + a + 1          // Définition de l'adresse de la ligne de l'ascii art correspondant à la rune, à la ligne imprimé + décalage

					fmt.Fprintf(w, s[adress]) // Impression de la ligne récupétrée
					f.WriteString(s[adress])

				}
				fmt.Fprintf(w, "\n") // Impression d'un retour à la ligne pour passer aux lignes suivantes
				f.WriteString(string(rune('\n')))
			}
		} else {
			fmt.Fprintf(w, "\n") // Impression d'une ligne seule si la cellule ne contient pas de texte
			f.WriteString(string(rune('\n')))
		}
	}
	fmt.Fprintf(w, "</pre>")
}

func renderTemplate(w http.ResponseWriter, s string) {
	t, err := template.ParseFiles("./templates/" + s + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/download.txt", download)

	fmt.Println("SERVER AWAITS: http://localhost:8080/")
	http.ListenAndServe(port, nil)
}

func Notfound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Error 404, this page doesn't exist")
}

func download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-disposition", "attachment; filename=download.txt")
	w.Header().Set("Content-type", ".txt; charset=UTF-8")
	http.ServeFile(w, r, "./download.txt")
}
