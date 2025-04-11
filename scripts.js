// Handle sign-in form submission
signInForm.addEventListener('submit', async (event) => {
    event.preventDefault();
  
    const email = document.getElementById('signInEmail').value;
    const password = document.getElementById('signInPassword').value;
  
    console.log(JSON.stringify({ email, password }));
  
    try {
      const response = await fetch('http://localhost:8080/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });
  
      if (response.ok) {
        // Automatically redirect to dashboard.html (no alert message)
        window.location.href = 'dashboard.html';
      } else {
        const errorText = await response.text();
        alert('Login failed: ' + errorText);
      }
    } catch (error) {
      alert('An error occurred: ' + error.message);
    }
  });
  
  // Handle sign-up form submission
  signUpForm.addEventListener('submit', async (event) => {
    event.preventDefault();
  
    const name = document.getElementById('signUpName').value;
    const email = document.getElementById('signUpEmail').value;
    const password = document.getElementById('signUpPassword').value;
  
    console.log('Registering with:', { name, email, password });
  
    try {
      const response = await fetch('http://localhost:8080/api/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, email, password }),
      });
  
      if (response.ok) {
        alert('Registration successful! Please login.');
        registerFormContainer.classList.remove('active');
        loginFormContainer.classList.add('active');
      } else {
        const errorText = await response.text();
        alert('Registration failed: ' + errorText);
      }
    } catch (error) {
      alert('An error occurred: ' + error.message);
    }
  });
  