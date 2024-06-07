function displayAddedProduct(product) {
    const addedProductContainer = document.getElementById("addedProductContainer");
    const addedProduct = document.getElementById("addedProduct");

    addedProduct.textContent = `Dish Name: ${product.dishName}, Country of Origin: ${product.countryOfOrigin}, Type of Dish: ${product.typeofDish}, Price: ${product.price}, Vegetarian: ${product.isVegetarian}, Non-Vegetarian: ${product.isNonVegetarian}, Vegan: ${product.isVegan}`;
    addedProductContainer.style.display = "block";

    setTimeout(function() {
        // After 2 seconds, refresh the page
        location.reload();
    }, 2000); // 2000 milliseconds = 2 seconds
}

document.getElementById('addItemForm').addEventListener('submit', function(event) {
    event.preventDefault(); // Prevent default form submission
    
    // Construct FormData object
    const formData = new FormData(this);

    // Construct newProduct object
    const newProduct = {
        dishName: formData.get("dishName"),
        countryOfOrigin: formData.get("countryOfOrigin"),
        typeOfDish: formData.get("typeOfDish"),
        price: parseInt(formData.get("dishPrice")), // Capture and parse price as integer
        isVegetarian: formData.get("isVegetarian") === "on",
        isNonVegetarian: formData.get("isNonVegetarian") === "on",
        isVegan: formData.get("isVegan") === "on",
    };

    // Log newProduct object (for testing)
    console.log("New Product:", newProduct);

    // Send newProduct object to backend (e.g., using fetch or AJAX)
    fetch('/api/addProduct', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(newProduct)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Product added:', data);
    })
    .catch(error => {
        console.error('Error:', error);
    });

    displayAddedProduct(newProduct);
});