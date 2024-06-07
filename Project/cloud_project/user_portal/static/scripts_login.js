document.getElementById('login-form').addEventListener('submit', function(event) {
    event.preventDefault(); // Prevent default form submission
 
    // Construct FormData object
    const formData = new FormData(this);
 
    // Construct newProduct object
    const newProduct = {
        email: formData.get('email'),
        password: formData.get('password')
    };
 
    // Log newProduct object (for testing)
    console.log("New Product:", newProduct);
 
    // Send newProduct object to backend (e.g., using fetch or AJAX)
    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(newProduct)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to login');
        }
        return response.json();
    })
    .then(data => {
        console.log('Login successful:', data);
        // Redirect the user to the login page or dashboard
        window.location.href = '/load.html'; // Adjust the URL as needed
    })
    .catch(error => {
        console.error('Error:', error);
        // Handle invalid response
        // This could be an alert or updating the page's content to show an error message
        alert('Login failed. Please check your credentials and try again.');
    });
});

document.getElementById('register-form').addEventListener('submit', function(event) {
    event.preventDefault(); // Prevent default form submission
 
    // Construct FormData object
    const formData = new FormData(this);
 
    // Construct newProduct object
    const newProduct = {
        username: formData.get('username'),
        email: formData.get('email'),
        password: formData.get('password')
    };
 
    // Log newProduct object (for testing)
    console.log("New Product:", newProduct);
 
    // Send newProduct object to backend (e.g., using fetch or AJAX)
    fetch('/api/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(newProduct)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Failed to register');
        }
        return response.json();
    })
    .then(data => {
        console.log('Register successful:', data);
        // Redirect the user to the login page or dashboard
        window.location.href = '/load.html'; // Adjust the URL as needed
    })
    .catch(error => {
        console.error('Error:', error);
        // Handle invalid response
        // This could be an alert or updating the page's content to show an error message
        alert('Register failed. Please check your credentials and try again.');
    });
});