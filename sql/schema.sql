CREATE TABLE IF NOT EXISTS blocklist (
  address INET
);

CREATE INDEX IF NOT EXISTS blocklist_address_gist_idx ON blocklist USING GIST (address inet_ops);
