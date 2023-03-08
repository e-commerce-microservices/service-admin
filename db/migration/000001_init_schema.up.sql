CREATE TYPE report_status_enum AS ENUM ('waiting', 'handled');

CREATE TABLE "report"  (
    "id" serial8 PRIMARY KEY,
    "product_id" serial8 NOT NULL,
    "description" text NOT NULL,
    "status" report_status_enum DEFAULT 'waiting'
);