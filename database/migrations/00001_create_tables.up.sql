CREATE TABLE IF NOT EXISTS author (
  id UUID PRIMARY KEY,
  name text NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  modified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS work (
    id UUID PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES author(id),
    title TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(author_id, title)
);

CREATE TABLE IF NOT EXISTS word (
    id UUID PRIMARY KEY,
    lemma_raw TEXT NOT NULL,
    lemma_rich TEXT NOT NULL,
    translation TEXT NOT NULL,
    lasla_frequency INT,
    known BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (lemma_raw, lemma_rich)
);

CREATE TABLE IF NOT EXISTS work_word (
    id UUID PRIMARY KEY,
    work_id UUID NOT NULL REFERENCES work(id),
    word_id UUID NOT NULL REFERENCES word(id),
    word_index INT NOT NULL,
    sentence_index INT NOT NULL,
    original_form TEXT NOT NULL,
    tag TEXT NOT NULL,
    morph_analysis TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (work_id, word_index)
);
