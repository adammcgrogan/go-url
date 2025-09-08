## Go URL Application

A simple web application built with Go to shorten long URLs. This project was created as a personal exercise to practice backend development, focusing on handling HTTP requests, managing redirects, and storing data.

### Features
- URL Shortening: Generate a short, unique URL from any long URL.
- Redirection: Seamlessly redirect users from the shortened URL to the original long URL.

## Todo
- [ ] History / List of created short URLs.
- [ ] Custom short URLs.
- [ ] All CRUD operations.

### Running Locally
Prerequisites: Go, Docker

Instructions:
1. Clone the repository: `git clone https://github.com/adammcgrogan/go-notes.git`, `cd go-notes`.
2. Set up the database: `docker-compose up -d`
3. Run the application: `go run .`
4. Open `http://localhost:8080` in your browser.
