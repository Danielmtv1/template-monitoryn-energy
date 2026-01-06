-- +goose Up
-- modify "events" table
ALTER TABLE "events" ADD COLUMN "plant_source_id" uuid NOT NULL, ADD CONSTRAINT "fk_events_plant_source" FOREIGN KEY ("plant_source_id") REFERENCES "energy_plants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- create index "idx_plant_source_id" to table: "events"
CREATE INDEX "idx_plant_source_id" ON "events" ("plant_source_id");

-- +goose Down
-- reverse: create index "idx_plant_source_id" to table: "events"
DROP INDEX "idx_plant_source_id";
-- reverse: modify "events" table
ALTER TABLE "events" DROP CONSTRAINT "fk_events_plant_source", DROP COLUMN "plant_source_id";
