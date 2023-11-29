package server

import (
	"html/template"
	"net/http"
)

type viewController struct{}

func newViewController() *viewController {
	return &viewController{}
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

 <!-- Display current pack sizes -->
    <div>
        <h2>Current Pack Sizes</h2>
        <ul id="packSizesList"></ul>
    </div>

    <!-- Form to edit pack sizes -->
    <div>
        <h2>Edit Pack Sizes</h2>
        <label for="newPackSizes">New Pack Sizes (comma-separated):</label>
        <input type="text" id="newPackSizes" placeholder="e.g., 23,31,53" required>
        <button onclick="updatePackSizes()">Update Pack Sizes</button>
    </div>


 	<h2>Calculate packs</h2>
    <label for="orderQuantity">Order Quantity:</label>
    <input type="number" id="orderQuantity" placeholder="e.g., 12001" required>
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
                        resultDiv.innerHTML += "<p>" + pack.packs_count + " pack(s) of " + pack.pack_size + " items</p>";
                    });
                }
            })
            .catch(error => console.error('Error:', error));
        }

 		// Function to display current pack sizes
        function displayPackSizes() {
            fetch('/api/v1/pack_sizes')
            .then(response => response.json())
            .then(data => {
                const packSizesList = document.getElementById('packSizesList');
                packSizesList.innerHTML = '';

                data.pack_sizes.forEach(size => {
                    const listItem = document.createElement('li');
                    listItem.textContent = size;
                    packSizesList.appendChild(listItem);
                });
            })
            .catch(error => console.error('Error:', error));
        }

        // Function to update pack sizes
        function updatePackSizes() {
            const newPackSizes = document.getElementById('newPackSizes').value.split(',').map(Number);

            fetch('/api/v1/pack_sizes', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ pack_sizes: newPackSizes }),
            })
            .then(response => response.json())
            .then(data => {
                // Display updated pack sizes
                displayPackSizes();
            })
            .catch(error => console.error('Error:', error));
        }

		// Initial display of pack sizes
        displayPackSizes();
    </script>
</body>
</html>
`))

func (c *viewController) viewHandler(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
