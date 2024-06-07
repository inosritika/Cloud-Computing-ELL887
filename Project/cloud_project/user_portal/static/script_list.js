function displayProducts(products) {
    const productList = document.getElementById('product-list');
    productList.innerHTML = ""; // Clear existing table rows
    console.log(products)
    products.forEach((product, index) => {
        const row = document.createElement("tr");

            // Serial number cell
            const serialNumberCell = document.createElement("td");
            serialNumberCell.textContent = index + 1;

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

            row.appendChild(serialNumberCell);
            row.appendChild(nameCell);
            row.appendChild(countryOfOriginCell);
            row.appendChild(typeOfDishCell);
            row.appendChild(priceCell);
            row.appendChild(vegetarianCell);
            row.appendChild(nonVegetarianCell);
            row.appendChild(veganCell);

            // Append the row to the table body
            productList.appendChild(row);
    });
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

function fetchFilterOptions() {
    fetch('/api/filterOptions')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch filter options');
            }
            return response.json();
        })
        .then(filterOptions => {
            console.log('Filter options:', filterOptions);
            // Call a function to display filter options in the UI
            displayFilterOptions(filterOptions);
        })
        .catch(error => {
            console.error('Error fetching filter options:', error);
        });
}

function displayFilterOptions(filterOptions) {
    const filterOptionsDiv = document.getElementById('filterOptions');
    
    // Clear existing filter options
    filterOptionsDiv.innerHTML = "";
    
    // // Country of Origin dropdown
    // const countrySelect = document.createElement('select');
    // countrySelect.id = 'countryOfOrigin'; // Add id attribute
    // countrySelect.name = 'CountryOfOrigin';
    // countrySelect.multiple = true;
    // const countryOptionDefault = document.createElement('option');
    // // countryOptionDefault.disabled = true;
    // // countryOptionDefault.selected = true;
    // // countryOptionDefault.hidden = true;
    // // countryOptionDefault.text;
    // countrySelect.appendChild(countryOptionDefault);
    // filterOptions.CountryOfOrigin.forEach(value => {
    //     const option = document.createElement('option');
    //     option.value = value;
    //     option.text = value;
    //     countrySelect.appendChild(option);
    // });
    // filterOptionsDiv.appendChild(countrySelect);
    
    // // Type of Dish dropdown
    // const typeSelect = document.createElement('select');
    // typeSelect.id = 'typeOfDish'; // Add id attribute
    // typeSelect.name = 'TypeOfDish';
    // typeSelect.multiple = true;
    // const typeOptionDefault = document.createElement('option');
    // // typeOptionDefault.disabled = true;
    // // typeOptionDefault.selected = true;
    // // typeOptionDefault.hidden = true;
    // // typeOptionDefault.text;
    // typeSelect.appendChild(typeOptionDefault);
    // filterOptions.TypeOfDish.forEach(value => {
    //     const option = document.createElement('option');
    //     option.value = value;
    //     option.text = value;
    //     typeSelect.appendChild(option);
    // });
    // filterOptionsDiv.appendChild(typeSelect);
    const countryDiv = document.createElement('div');
    countryDiv.innerHTML = '<h3>Country of Origin</h3>';
    filterOptions.CountryOfOrigin.forEach(value => {
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.name = 'countryOfOrigin';
        checkbox.value = value;
        const label = document.createElement('label');
        label.textContent = value;
        label.appendChild(checkbox);
        countryDiv.appendChild(label);
    });
    filterOptionsDiv.appendChild(countryDiv);

    // Type of Dish checkboxes
    const typeDiv = document.createElement('div');
    typeDiv.innerHTML = '<h3>Type of Dish</h3>';
    filterOptions.TypeOfDish.forEach(value => {
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.name = 'typeOfDish';
        checkbox.value = value;
        const label = document.createElement('label');
        label.textContent = value;
        label.appendChild(checkbox);
        typeDiv.appendChild(label);
    });
    filterOptionsDiv.appendChild(typeDiv);
    
    // Price range slider
    const priceRangeLabel = document.createElement('label');
    priceRangeLabel.for = 'priceRange';
    priceRangeLabel.textContent = 'Price Range:';
    filterOptionsDiv.appendChild(priceRangeLabel);

    const priceRangeInput = document.createElement('input');
    priceRangeInput.type = 'range';
    priceRangeInput.id = 'priceRange';
    priceRangeInput.name = 'priceRange';
    priceRangeInput.min = 0;
    priceRangeInput.max = filterOptions.maxPrice;
    priceRangeInput.value = filterOptions.maxPrice; // Set initial value to max price
    filterOptionsDiv.appendChild(priceRangeInput);

    // Checkbox for Vegetarian
    const vegetarianCheckbox = document.createElement('input');
    vegetarianCheckbox.type = 'checkbox';
    vegetarianCheckbox.id = 'isVegetarian';
    vegetarianCheckbox.name = 'isVegetarian';
    const vegetarianLabel = document.createElement('label');
    vegetarianLabel.for = 'isVegetarian';
    vegetarianLabel.textContent = 'Vegetarian';
    filterOptionsDiv.appendChild(vegetarianCheckbox);
    filterOptionsDiv.appendChild(vegetarianLabel);

    // Checkbox for Non-Vegetarian
    const nonVegetarianCheckbox = document.createElement('input');
    nonVegetarianCheckbox.type = 'checkbox';
    nonVegetarianCheckbox.id = 'isNonVegetarian';
    nonVegetarianCheckbox.name = 'isNonVegetarian';
    const nonVegetarianLabel = document.createElement('label');
    nonVegetarianLabel.for = 'isNonVegetarian';
    nonVegetarianLabel.textContent = 'Non-Vegetarian';
    filterOptionsDiv.appendChild(nonVegetarianCheckbox);
    filterOptionsDiv.appendChild(nonVegetarianLabel);

    // Checkbox for Vegan
    const veganCheckbox = document.createElement('input');
    veganCheckbox.type = 'checkbox';
    veganCheckbox.id = 'isVegan';
    veganCheckbox.name = 'isVegan';
    const veganLabel = document.createElement('label');
    veganLabel.for = 'isVegan';
    veganLabel.textContent = 'Vegan';
    filterOptionsDiv.appendChild(veganCheckbox);
    filterOptionsDiv.appendChild(veganLabel);
}

fetchProducts();

function applyFilters() {
    const selectedFilters = {};

    // // Get selected options for Country of Origin
    // const countrySelect = document.getElementById('countryOfOrigin');
    // selectedFilters['CountryOfOrigin'] = Array.from(countrySelect.selectedOptions).map(option => option.value);

    // // Get selected options for Type of Dish
    // const typeSelect = document.getElementById('typeOfDish');
    // selectedFilters['TypeOfDish'] = Array.from(typeSelect.selectedOptions).map(option => option.value);
    
    // Get selected options for Country of Origin
    const countryCheckboxes = document.querySelectorAll('input[name="countryOfOrigin"]:checked');
    selectedFilters.countryOfOrigin = Array.from(countryCheckboxes).map(checkbox => checkbox.value);

    // Get selected options for Type of Dish
    const typeCheckboxes = document.querySelectorAll('input[name="typeOfDish"]:checked');
    selectedFilters.typeOfDish = Array.from(typeCheckboxes).map(checkbox => checkbox.value);

    // Get selected price range
    const priceRangeInput = document.getElementById('priceRange');
    selectedFilters['Price'] = [priceRangeInput.value]; // Wrap in array as it's expected to be an array of strings

    // Get selected checkboxes for Vegetarian, Non-Vegetarian, Vegan
    selectedFilters['IsVegetarian'] = document.getElementById('isVegetarian').checked ? ['true'] : [];
    selectedFilters['IsNonVegetarian'] = document.getElementById('isNonVegetarian').checked ? ['true'] : [];
    selectedFilters['IsVegan'] = document.getElementById('isVegan').checked ? ['true'] : [];

    console.log('Selected filters:', selectedFilters);

    // Now you can send the selectedFilters object to the backend for processing
    // Example fetch call:
    fetch('/api/filterProducts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(selectedFilters)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to apply filters');
        }
        return response.json();
    })
    .then(filteredProducts => {
        // Process filtered products as needed
        console.log('Filtered products:', filteredProducts);
        displayProducts(filteredProducts);
    })
    .catch(error => {
        console.error('Error applying filters:', error);
    });
}

// Fetch filter options when the page loads
fetchFilterOptions();

// Add event listener for the "Apply Filters" button
document.getElementById('applyFiltersBtn').addEventListener('click', applyFilters);