# URL Shortener

A modern, mobile-friendly URL shortener built with Go (backend) and Tailwind CSS (frontend).

## Features
- Shorten long URLs to easy-to-share short links
- Responsive, clean UI inspired by modern SaaS design
- Copy-to-clipboard button for short URLs
- Info popup explaining the app
- In-memory storage (easy to add a database for persistence)
- CORS enabled for easy deployment with static frontends

## Getting Started

### 1. Run Locally

#### Backend (Go)
1. Make sure you have Go installed (`go version`)
2. Clone this repo and `cd` into the project directory
3. Run:
   ```sh
   go run main.go
   ```
4. The server will start on `http://localhost:8080`

#### Frontend
- Open `index.html` in your browser (it will work locally if served by the Go backend)
- Or, deploy the HTML to a static host (Vercel, Netlify, etc.) and point the API URL in the JS to your backend

### 2. Deploy to Production

#### Backend
- Deploy `main.go` to a platform that supports Go web servers (Railway, Render, Fly.io, etc.)
- Make sure the backend is accessible at a public URL (e.g., `https://your-app.up.railway.app`)
- CORS is enabled by default

#### Frontend
- Deploy `index.html` to any static host (Vercel, Netlify, etc.)
- Update the API URL in the JS to point to your backend (see `API_BASE` in the script)

## Customization
- To persist URLs, connect the Go backend to a database (e.g., SQLite, PostgreSQL, Redis)
- To change the UI, edit `index.html` and use Tailwind utility classes

## Credits
- Built by [Your Name]
- UI inspired by SafetyKit and modern SaaS design

---

Feel free to fork, modify, and use for your own projects! 