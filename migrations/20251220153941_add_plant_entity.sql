-- +goose Up
-- create "energy_plants" table
CREATE TABLE "energy_plants" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "plant_name" character varying(255) NOT NULL,
  "location" character varying(255) NULL,
  "capacity_mw" numeric NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_energy_plants_deleted_at" to table: "energy_plants"
CREATE INDEX "idx_energy_plants_deleted_at" ON "energy_plants" ("deleted_at");

-- +goose Down
-- reverse: create index "idx_energy_plants_deleted_at" to table: "energy_plants"
DROP INDEX "idx_energy_plants_deleted_at";
-- reverse: create "energy_plants" table
DROP TABLE "energy_plants";
