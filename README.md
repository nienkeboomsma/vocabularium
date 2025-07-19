# Vocabularium

## Features

A web interface for the [Collatinus](https://github.com/biblissima/collatinus) Latin lemmatiser. It extracts all the lemmas from uploaded texts and compiles frequency lists (by work, author or entire corpus) or glossaries. Words can be marked as 'known' and filtered out of all lists. Available in all languages supported by Collatinus (which are Basque, Catalan, Dutch, English, French, Galician, German, Italian, Portuguese and Spanish).

## Installation

1. Create a `.env` file:

```sh
cp .env.example .env
```

2. Make sure to replace the placeholder values.

3. Build the containers:

```sh
docker compose up --build -d
```

## Usage

1. Start the containers:

```sh
docker compose up
```

2. Open `http://localhost:4321` in your browser.

3. Stop the containers:

```sh
docker compose down
```
