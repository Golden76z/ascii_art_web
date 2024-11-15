package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//A Parse.File is a local representation of a file that is saved to the Parse cloud.

// The main function will execute all my handlers
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)
	http.HandleFunc("/result", resultHandler)
	http.HandleFunc("/download", downloadHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server listening on port 8080...")

	//Port qui va être utilisé pour le site internet
	http.ListenAndServe(":8080", nil)
}

// Handler qui gère le téléchargement de mon fichier texte
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("file")
	http.ServeFile(w, r, fileName)
}

// homeHandler: Serves the main page with text input and style selection
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		//Si le chemin spécifié est different de celui qui existe alors on renvoie une erreur 404
		page, _ := template.ParseFiles("template/error404.html")
		_ = page.Execute(w, nil)
		return
	} else {
		//Dans le cas contraire on stocke les metadonnées dans une variable
		page, err := template.ParseFiles("template/index.html")
		if err != nil {
			//Message d'erreur si le stockage fail
			page, _ := template.ParseFiles("template/error404.html")
			_ = page.Execute(w, nil)
			return
		}
		//Execution de la page html
		err = page.Execute(w, nil)
		if err != nil {
			//Gestion d'erreur
			http.Error(w, "serveur error", http.StatusNotFound)
			log.Printf("error template %v", err)
			return
		}
	}
}

// asciiArtHandler: Handles POST requests for generating ASCII art
func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	page, err := template.ParseFiles("template/ascii.html")
	if err != nil {
		page, _ := template.ParseFiles("template/error404.html")
		_ = page.Execute(w, nil)
		return
	}
	err = page.Execute(w, nil)
	if err != nil {
		http.Error(w, "serveur error", http.StatusNotFound)
		log.Printf("error template %v", err)
		return
	}
}

// Handler qui s'occupe de ma page de resultat
func resultHandler(w http.ResponseWriter, r *http.Request) {
	page, err := template.ParseFiles("template/result.html")
	if err != nil {
		page, _ := template.ParseFiles("template/error404.html")
		_ = page.Execute(w, nil)
		return
	}

	//Stockage des données reçus dans les formulaire dans des variables
	text := r.FormValue("text")
	banner := r.FormValue("banner")
	color := r.FormValue("colorhex")

	//Stockage du resultat de ma fonction ascii art dans une variable
	result, err := AsciiArt(text, banner)
	fmt.Println(result)

	//Gestion d'erreur lors de l'utilisation de la formule
	if err != nil {
		http.Error(w, "serveur error", http.StatusNotFound)
		log.Printf("error template %v", err)
		return
	}

	//Création d'un fichier temporaire qui va stocker le resultat de mon ascii art
	tempFile, err := os.CreateTemp("", "ascii-art-*.txt")
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Printf("Error creating temp file: %v", err)
		return
	}
	defer tempFile.Close()

	//Ecriture du resultat dans le fichier texte
	_, err = tempFile.WriteString(result)
	if err != nil {
		http.Error(w, "serveur error", http.StatusNotFound)
		log.Printf("error template %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = page.Execute(w, map[string]interface{}{
		//Convertir le string en html pour que la balise et la couleur soient appliquées
		"Result":       template.HTML(StringtoHtml(result, color)),
		"Color":        color,           // Add the color value to the template data
		"DownloadLink": tempFile.Name(), // Add the download link to the template data
	})
	if err != nil {
		http.Error(w, "serveur error", http.StatusNotFound)
		log.Printf("error template %v", err)
		return
	}
}

func AsciiArt(text string, banner string) (s string, err error) {
	result := ""
	tabstandard := []string{}
	file, err := os.Open(banner + ".txt")
	if err != nil {
		fmt.Println("Ouverture du fichier banner impossible")
		os.Exit(2)
	}

	//Lecture lignes par lignes
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {

		//Stockage ligne par ligne du fichier banner dans un tableau de string
		tabstandard = append(tabstandard, fileScanner.Text())
	}
	tabstandard = append(tabstandard, " ")
	tabascii := []string{}

	//Tableau qui va faire l'intermédiaire entre le string et le standard.txt
	for i := 32; i <= 126; i++ {
		tabascii = append(tabascii, string(rune(i)))
	}
	tabword := strings.Split(text, "\\n")
	for _, word := range tabword {

		//Appel de ma fonction qui converti mon string en ascii art
		result += StringToAscii(word, tabascii, tabstandard)
	}
	return result, err
}

func StringToAscii(word string, tab2 []string, tab3 []string) string {
	result := ""
	result2 := ""
	for i := 0; i <= 8; i++ {
		for j := 0; j < len(word); j++ {
			for k := 0; k < len(tab2); k++ {
				if string(word[j]) == tab2[k] {
					result = tab3[k*9+i+1]
					fmt.Print(result)
					result2 += result
				}
			}
		}
		if i != 8 {
			fmt.Println("")
			result2 += "\n"
		}
	}
	return result2
}

// Fonction qui va appliquer ma couleur à ma string en la convertissant au format html
func StringtoHtml(s string, color string) string {
	//Si aucune couleure n'est séléctionnée alors on met blanc de base
	if color == "" {
		color = "White"
	}
	//On englobe la couleur dans la première partie de la balise
	htmlcolor := "<p style=\"color: " + color + " ;\">"
	for _, i := range s {
		//On remplace les espaces par &nbsp; qui représente un espace en html
		if i == 32 {
			htmlcolor += "&nbsp;"
		} else if i == 10 {
			//On remplace les 10 (\n) par des balises <br> qui représentent des saut de ligne en html
			htmlcolor += "<br>"
		} else {
			htmlcolor += string(i)
		}
	}
	//On ferme la balise
	return htmlcolor + "</p>"
}
