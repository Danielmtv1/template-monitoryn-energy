-- +goose Up
-- create "examples" table
CREATE TABLE "examples" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

-- +goose Down
-- reverse: create "examples" table
DROP TABLE "examples";
