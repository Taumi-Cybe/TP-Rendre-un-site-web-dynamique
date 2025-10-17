package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Image       string
	Price       float64
	Stock       int
	Discount    int
}

var products = []Product{
	{1, "PALACE PULL A CAPUCHE CHASSEUR", "Sweat vert style chasseur", "16A.webp", 130, 10, 15},
	{2, "PALACE PULL A CAPUCHE MARINE", "Sweat bleu marine streetwear", "18A.webp", 130, 8, 0},
	{3, "PALACE PANTALON CARGO PASSEPORT NOIR", "Cargo noir style passport", "33B.webp", 140, 5, 10},
	{4, "PALACE HOOD MOJITO", "Sweat mojito effet lavé", "21A.webp", 135, 6, 5},
	{5, "PALACE PANTALON BOSSY JEAN STONE", "Jean stone délavé", "34B.webp", 110, 15, 0},
	{6, "PALACE PANTALON CARGO GORE-TEX TKN NOIR", "Cargo Gore-Tex noir", "22A.webp", 140, 4, 20},
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/product/", productHandler)
	http.HandleFunc("/add", addHandler)

	log.Println("Serveur lancé sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Products": products,
	}
	templates.ExecuteTemplate(w, "index.html", data)
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/product/")
	idStr = path.Clean("/" + idStr)
	idStr = strings.TrimPrefix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	for _, p := range products {
		if p.ID == id {
			data := map[string]interface{}{
				"Product": p,
			}
			templates.ExecuteTemplate(w, "product.html", data)
			return
		}
	}
	http.NotFound(w, r)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templates.ExecuteTemplate(w, "add.html", nil)
	case http.MethodPost:
		r.ParseForm()
		name := r.FormValue("name")
		desc := r.FormValue("description")
		image := r.FormValue("image")
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		stock, _ := strconv.Atoi(r.FormValue("stock"))
		discount, _ := strconv.Atoi(r.FormValue("discount"))

		if image == "" {
			image = "default.webp"
		}

		newID := len(products) + 1
		newProduct := Product{
			ID:          newID,
			Name:        name,
			Description: desc,
			Image:       image,
			Price:       price,
			Stock:       stock,
			Discount:    discount,
		}
		products = append(products, newProduct)
		http.Redirect(w, r, "/product/"+strconv.Itoa(newID), http.StatusSeeOther)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}