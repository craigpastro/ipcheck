CREATE TABLE IF NOT EXISTS blocklist (
  address INET NOT NULL,
  filename TEXT NOT NULL,
  source_file_date TEXT
);

CREATE INDEX IF NOT EXISTS address_idx ON blocklist USING GIST (address inet_ops);
