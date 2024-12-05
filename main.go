package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"html/template"

	"github.com/gocarina/gocsv"
)

// Struct representing a row in the CSV
type Record struct {
	Category  string  `csv:"Category"`
	Spent     float64 `csv:"Spent"`
	Budget    float64 `csv:"Budget"`
	Fulfilled string  `csv:"Fulfilled"`
	Result    string  `csv:"Result"`
}

func main() {
	// Read the CSV file
	file, err := os.OpenFile("data.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	var records []Record
	if err := gocsv.UnmarshalFile(file, &records); err != nil {
		log.Fatalf("Error parsing CSV: %v", err)
	}

	// Parse the inline HTML template
	tmpl, err := template.New("csvTemplate").Parse(`
    <!DOCTYPE html>
    <html>
    <head>
        <title>CSV Data</title>
        <style>
            .card {
                border: 1px solid #ccc;
                padding: 10px;
                margin: 10px;
                border-radius: 5px;
                display: inline-block;
                width: 200px;
            }
            .negative {
                background-color: #f8d7da;
            }
            .positive {
                background-color: #d4edda;
            }
        </style>
    </head>
    <body>
        <h1>CSV Data</h1>
        <div>
            {{range .}}
            <div class="card {{if gt .Result 0}}positive{{else if lt .Result 0}}negative{{end}}">
                <p>Category: {{.Category}}</p>
                <p>Result: ₪{{.Result}}</p>
            </div>
            {{end}}
        </div>
    </body>
    </html>
    `)
	if err != nil {
		log.Fatalf("Error parsing inline template: %v", err)
	}

	// Set up the HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, records)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	// Start the server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
