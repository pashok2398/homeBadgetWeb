package main

import (
	"encoding/csv"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Helper functions for template
var funcs = template.FuncMap{
	"isNegative": func(value string) bool {
		val, err := strconv.ParseFloat(strings.ReplaceAll(value, "₪", ""), 64)
		return err == nil && val < 0
	},
	"isZeroOrPositive": func(value string) bool {
		val, err := strconv.ParseFloat(strings.ReplaceAll(value, "₪", ""), 64)
		return err == nil && val >= 0
	},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the categories from the query parameter "categories"
		categoriesParam := r.URL.Query().Get("categories")
		var targetCategories []string
		if categoriesParam != "" {
			// Split the categories by comma to create a slice
			targetCategories = strings.Split(categoriesParam, ",")
		}

		// Read the CSV file
		file, err := os.Open("data.csv")
		if err != nil {
			http.Error(w, "Unable to open CSV file", http.StatusInternalServerError)
			log.Printf("Error opening file: %v", err)
			return
		}
		defer file.Close()

		// Parse the CSV file
		reader := csv.NewReader(file)
		reader.Comma = ';' // Adjust the delimiter if needed
		records, err := reader.ReadAll()
		if err != nil {
			http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
			log.Printf("Error reading file: %v", err)
			return
		}

		// Determine whether to filter or display all records
		filteredRecords := make([][]string, 0)
		if targetCategories == nil { // Display all records except the header
			filteredRecords = records[1:]
		} else {
			// Filter records based on the target categories
			for _, record := range records[1:] { // Skip header row
				for _, category := range targetCategories {
					if record[0] == category { // Assuming category is in the first column
						filteredRecords = append(filteredRecords, record)
						break
					}
				}
			}
		}

		// Parse the template
		tmpl, err := template.New("csvTemplate").Funcs(funcs).Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Filtered CSV Data</title>
    <meta http-equiv="refresh" content="5">
    <link rel="shortcut icon" href="#">
    <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css">
    <style>
        body {
            background: linear-gradient(135deg, #e3f2fd, #bbdefb);
            font-family: 'Roboto', sans-serif;
            color: #333;
            display: flex;
            flex-wrap: wrap;
            justify-content: center;
            padding: 20px;
        }
        .container {
            max-width: 1200px;
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
        }
        .card {
            border-radius: 12px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            overflow: hidden;
            transition: transform 0.3s, box-shadow 0.3s;
            background-color: #ffffff;
            text-align: center;
        }
        .card:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
        }
        .card-header {
            padding: 12px;
            font-size: 1.2rem;
            font-weight: 600;
            color: #ffffff;
        }
        .card-body {
            padding: 20px;
            font-size: 1rem;
            border-radius: 0 0 12px 12px;
        }
        .card-body.positive {
            background-color: #4caf50; /* Green */
            color: white;
        }
        .card-body.negative {
            background-color: #e53935; /* Red */
            color: white;
        }
        .card-title {
            font-size: 1.5rem;
            margin: 10px 0;
            color: #333;
        }
        .card-text {
            font-size: 1.25rem;
            font-weight: bold;
            color: #555;
        }
    </style>
</head>
<body>
    <div class="container">
        {{if .}}
            {{range .}}
                {{if isZeroOrPositive (index . 4)}}
                    <div class="card">
                        <div class="card-header" style="background-color: #28a745;">
                            {{index . 0}}
                        </div>
                        <div class="card-body positive">
                            <p class="card-text">{{index . 4}}</p>
                        </div>
                    </div>
                {{else}}
                    <div class="card">
                        <div class="card-header" style="background-color: #dc3545;">
                            {{index . 0}}
                        </div>
                        <div class="card-body negative">
                            <p class="card-text">{{index . 4}}</p>
                        </div>
                    </div>
                {{end}}
            {{end}}
        {{else}}
            <p>No matching records found or no data available.</p>
        {{end}}
    </div>
</body>
</html>
        `)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			log.Printf("Error parsing template: %v", err)
			return
		}

		// Execute the template with the filtered data
		err = tmpl.Execute(w, filteredRecords)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
		}
	})

	// Start the web server
	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
