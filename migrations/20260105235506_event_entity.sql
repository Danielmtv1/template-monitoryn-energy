-- +goose Up
-- create "events" table
CREATE TABLE "events" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "event_type" character varying(100) NOT NULL,
  "source" character varying(255) NULL,
  "data" text NULL,
  "metadata" text NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_created_at" to table: "events"
CREATE INDEX "idx_created_at" ON "events" ("created_at");
-- create index "idx_event_type" to table: "events"
CREATE INDEX "idx_event_type" ON "events" ("event_type");

-- +goose Down
-- reverse: create index "idx_event_type" to table: "events"
DROP INDEX "idx_event_type";
-- reverse: create index "idx_created_at" to table: "events"
DROP INDEX "idx_created_at";
-- reverse: create "events" table
DROP TABLE "events";
