function displayProducts(products) {
    const productList = document.getElementById('product-list');
    productList.innerHTML = ""; // Clear existing table rows
    console.log(products)
    products.forEach((product) => {
        const row = document.createElement("tr");

            // Serial number cell
            const serialNumberCell = document.createElement("td");
            serialNumberCell.textContent = product.ID;

            // Dish name cell
            const nameCell = document.createElement("td");
            nameCell.textContent = product.DishName;

            // Country of origin cell
            const countryOfOriginCell = document.createElement("td");
            countryOfOriginCell.textContent = product.CountryOfOrigin;

            // Type of dish cell
            const typeOfDishCell = document.createElement("td");
            typeOfDishCell.textContent = product.TypeOfDish;

            // Price cell
            const priceCell = document.createElement("td");
            priceCell.textContent = product.Price;

            // Vegetarian cell
            const vegetarianCell = document.createElement("td");
            vegetarianCell.textContent = product.IsVegetarian ? "Yes" : "No";

            // Non Vegetarian cell
            const nonVegetarianCell = document.createElement("td");
            nonVegetarianCell.textContent = product.IsNonVegetarian ? "Yes" : "No";

            // Vegan cell
            const veganCell = document.createElement("td");
            veganCell.textContent = product.IsVegan ? "Yes" : "No";

            const actionCell = document.createElement("td");
            const approveButton = document.createElement("button");
            approveButton.textContent = "Approve";
            approveButton.addEventListener("click", () => {
                // Send an Approve request 
                approveProducts(product.ID);
            });
            const rejectButton = document.createElement("button");
            rejectButton.textContent = "Reject";
            rejectButton.addEventListener("click", () => {
                // Send a Reject request 
                rejectProducts(product.ID);
            });
            actionCell.appendChild(approveButton);
            actionCell.appendChild(rejectButton);

            row.appendChild(serialNumberCell);
            row.appendChild(nameCell);
            row.appendChild(countryOfOriginCell);
            row.appendChild(typeOfDishCell);
            row.appendChild(priceCell);
            row.appendChild(vegetarianCell);
            row.appendChild(nonVegetarianCell);
            row.appendChild(veganCell);
            row.appendChild(actionCell)

            // Append the row to the table body
            productList.appendChild(row);
    });
}

function approveProducts(id){
    const requestData = {
        ID: id
    };

    // Send a POST request to your backend API to add the approved product
    fetch('/api/statusApprove', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to approve product');
        }
        // Handle success response (if needed)
        console.log('Product approved successfully');
    })
    .catch(error => {
        console.error('Error approving product:', error);
    });

    fetchProducts();
}

function rejectProducts(id) {
    console.log(id);
    const requestData = {
        ID: id
    };

    fetch('/api/statusReject', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(requestData)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to reject product');
        }
    })
    .catch(error => {
        console.error('Error rejecting product:', error);
    });

    fetchProducts();
}

function fetchProducts() {
    fetch('/api/listProducts') // Assuming this endpoint is configured on your backend server
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch products');
            }
            return response.json();
        })
        .then(products => {
            console.log(products)
            displayProducts(products);
        })
        .catch(error => {
            console.error('Error fetching products:', error);
        });
}

document.getElementById('logout-button').addEventListener('click', function() {
    console.log("Logging out")
    fetch('/api/logout', {
        method: 'POST'
    })
    .then(response => {
        if (response.ok) {
            // Redirect to the login page
            window.location.href = '/index.html'; // Adjust the URL to your login page
        } else {
            throw new Error('Logout failed');
        }
    })
    .catch(error => {
        console.log('Error:', error);
        console.error('Error:', error);
    });
});

fetchProducts();

