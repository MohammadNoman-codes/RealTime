const app = document.getElementById('app');
app.innerHTML = `
        <div class="container">
            <div class="form-container">
                <h1>Sign In</h1>
                <input type="email" id="email" placeholder="Email" required />
                <input type="password" id="password" placeholder="Password" required />
                <button id="loginBtn">Sign In</button>
                <p>Don't have an account? <button id="switchToSignUp">Sign Up</button></p>
            </div>
        </div>
    `;

// Function to render the login form
function renderLoginForm() {
    app.innerHTML = `
        <div class="container">
            <div class="form-container">
                <h1>Sign In</h1>
                <input type="email" id="email" placeholder="Email" required />
                <input type="password" id="password" placeholder="Password" required />
                <button id="loginBtn">Sign In</button>
                <p>Don't have an account? <button id="switchToSignUp">Sign Up</button></p>
            </div>
        </div>
    `;

    // Attach event listeners
    document.getElementById('loginBtn').addEventListener('click', handleLogin);
    document.getElementById('switchToSignUp').addEventListener('click', renderSignupForm);
}

// Function to render the signup form
function renderSignupForm() {
    app.innerHTML = `
        <div class="container">
            <div class="form-container">
                <h1>Create Account</h1>
                <input type="email" id="email" placeholder="Email" required />
                <input type="text" id="username" placeholder="Username" required />
                <input type="password" id="password" placeholder="Password" required />
                <button id="signupBtn">Sign Up</button>
                <p>Already have an account? <button id="switchToSignIn">Sign In</button></p>
            </div>
        </div>
    `;

    // Attach event listeners
    document.getElementById('signupBtn').addEventListener('click', handleSignup);
    document.getElementById('switchToSignIn').addEventListener('click', renderLoginForm);
}

// Function to handle login
async function handleLogin() {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const response = await fetch('http://localhost:1703/signin', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
    });

    const result = await response.json();

    if (response.ok) {
        alert(result.message); // You can redirect or perform other actions
    } else {
        alert(`Error: ${result.message}`);
    }
}

// Function to handle signup
async function handleSignup() {
    const email = document.getElementById('email').value;
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const response = await fetch('/signup', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, username, password })
    });

    const result = await response.json();

    if (response.ok) {
        alert(result.message); // Redirect or perform other actions
    } else {
        alert(`Error: ${result.message}`);
    }
}

// Initialize the app with the login form
renderLoginForm();
