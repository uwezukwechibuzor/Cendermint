CREATE TABLE blocks (
  height INT PRIMARY KEY UNIQUE NOT NULL,
  block_hash TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL,
);