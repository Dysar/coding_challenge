package server

import (
	"challenge/internal/model"
	"html/template"
	"net/http"
)

type PageVariables struct {
	Title         string
	OrderQuantity int
	Packs         []model.PackDetails
}

var indexTemplate = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Pack Calculator</title>
</head>
<body>
    <h1>Pack Calculator</h1>
    <label for="orderQuantity">Order Quantity:</label>
    <input type="number" id="orderQuantity" required>
    <button onclick="calculatePacks()">Calculate Packs</button>
    
    <div id="result"></div>

    <script>
        function calculatePacks() {
            const orderQuantity = document.getElementById('orderQuantity').value;

            fetch('/api/v1/calculate_packs', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ order_quantity: parseInt(orderQuantity) }),
            })
            .then(response => response.json())
            .then(data => {
                const resultDiv = document.getElementById('result');
                resultDiv.innerHTML = '';

                if (data.error) {
                    resultDiv.innerHTML = "<p>Error: " + data.error + "</p>";
                } else {
                    resultDiv.innerHTML = '<p>Packs Needed:</p>';
                    data.packs_needed.forEach(pack => {
                        resultDiv.innerHTML += "<p>" + pack.packs_count + " packs of " + pack.pack_size + " items</p>";
                    });
                }
            })
            .catch(error => console.error('Error:', error));
        }
    </script>
</body>
</html>
`))

func viewHandler(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
