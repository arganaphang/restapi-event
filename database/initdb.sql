CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS "transactions";

CREATE TABLE IF NOT EXISTS "transactions"(
  "transaction_id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "id" BIGINT NOT NULL,
  "customer" VARCHAR NOT NULL,
  "quantity" SMALLINT NOT NULL,
  "price" DECIMAL(10, 2) NOT NULL,
  "timestamp" TIMESTAMP NOT NULL
);