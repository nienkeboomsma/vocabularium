# Project name tbd

## Installation

Create a `.env` file:

```sh
cp .env.example .env
```

Make sure to replace the placeholder values.

To build the containers:

```sh
docker compose up --build -d
```

Thereafter, use the following commands to start and stop the containers respectively:

```sh
docker compose up
docker compose down
```

## Usage

### Adding a work
Send a POST request to `http://localhost:4321/lemmatise` with the following `multipart/form-data` fields:

| Field    | Type                   |
| -------- | ---------------------- |
| `file`   | .txt file              |
| `author` | string                 |
| `title`  | string                 |
| `type`   | `"verse"` or `"prose"` |

With `curl` it would look like this:

```sh
curl -X POST \
  -F "file=@./amphitryo.txt" \
  -F "author=Plautus" \
  -F "title=Amphitryo" \
  -F "type=verse" \
  http://localhost:4321/lemmatise
```

### Browsing

Open `http://localhost:4321/works` in your browser.
