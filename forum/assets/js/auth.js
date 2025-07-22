// For Login
document.getElementById('login-form').addEventListener('submit', function(e) {
    e.preventDefault();

    const formData = new FormData(this);
    const errorElement = document.getElementById('login-error');

    fetch('/login', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (response.ok) {
            window.location.reload(); // "Redirects" to be on the same page
        } else {
            return response.text().then(text => {
                throw new Error(text || 'Login failed');
            });
        }
    })
    .catch(error => {
        console.error('Error:', error);
        errorElement.textContent = error.message || 'An error occurred. Please try again.';
        errorElement.style.display = 'block';
    });
});

// For Registration
document.getElementById('register-form').addEventListener('submit', function(e) {
    e.preventDefault();

    const formData = new FormData(this);
    const errorElement = document.getElementById('register-error');

    fetch('/register', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (response.ok) {
            window.location.href = '/'; // Redirect to login page after successful registration
        } else {
            return response.text().then(text => {
                throw new Error(text);
            });
        }
    })
    .catch(error => {
        console.error('Error:', error);
        errorElement.textContent = error.message;
        errorElement.style.display = 'block';

    });
});
