# Link Shortener

A fast and lightweight link shortener built in Go. Paste a long URL and get a short one back. Links expire after 2 hours and click counts are tracked in real time.

## Features

- Shorten any URL instantly
- Click tracking on every link
- Links expire automatically after 2 hours
- Clean dark mode web interface
- Stats page showing all active links
- Docker support for easy setup

## Getting Started

### Option 1 — Run with Docker (recommended)

Make sure you have [Docker](https://www.docker.com/products/docker-desktop/) installed, then run:

```bash
docker build -t url-shortener .
docker run -p 8080:8080 url-shortener
```

Open your browser and go to `http://localhost:8080`

### Option 2 — Run with Go

Make sure you have [Go](https://go.dev/dl/) installed, then run:

```bash
git clone https://github.com/xtekkis/url-shortener.git
cd url-shortener
go run main.go
```

Open your browser and go to `http://localhost:8080`

## Usage

1. Paste a long URL into the input box
2. Click **Shorten Link**
3. Copy and share your short link
4. Visit `/stats` to see all active links and their click counts

## Tech Stack

- **Go** — backend server and routing
- **HTML/CSS/JavaScript** — frontend interface
- **Docker** — containerisation