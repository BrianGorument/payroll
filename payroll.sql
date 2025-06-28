
CREATE TYPE user_role AS ENUM ('employee', 'admin');
CREATE TYPE action AS ENUM ('CREATE', 'UPDATE', 'DELETE');

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY, -- Menggunakan SERIAL untuk auto-increment INT ID
  "username" varchar(50) UNIQUE NOT NULL,
  "password" varchar(255) NOT NULL,
  "salary" decimal(15,2),
  "role" user_role NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" INT NOT NULL, -- Menggunakan INT untuk foreign key
  "updated_by" INT NOT NULL -- Menggunakan INT untuk foreign key
);

COMMENT ON COLUMN "users"."salary" IS 'NULL untuk admin';


CREATE TABLE "payroll_periods" (
  "id" SERIAL PRIMARY KEY,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  "is_processed" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" INT NOT NULL,
  "updated_by" INT NOT NULL
);



CREATE TABLE "attendances" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "payroll_period_id" INT NOT NULL,
  "attendance_date" date NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" INT NOT NULL,
  "updated_by" INT NOT NULL
);


CREATE TABLE "overtimes" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "payroll_period_id" INT NOT NULL,
  "overtime_date" date NOT NULL,
  "hours" integer NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" INT NOT NULL,
  "updated_by" INT NOT NULL
);

COMMENT ON COLUMN "overtimes"."hours" IS 'Maksimum 3 jam per hari';


CREATE TABLE "reimbursements" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "payroll_period_id" INT NOT NULL,
  "amount" decimal(15,2) NOT NULL,
  "description" text,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" INT NOT NULL,
  "updated_by" INT NOT NULL
);


CREATE TABLE "payslips" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "payroll_period_id" INT NOT NULL,
  "base_salary" decimal(15,2) NOT NULL,
  "overtime_pay" decimal(15,2) NOT NULL DEFAULT 0,
  "reimbursement_pay" decimal(15,2) NOT NULL DEFAULT 0,
  "total_pay" decimal(15,2) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" INT NOT NULL,
  "updated_by" INT NOT NULL
);


CREATE TABLE "payroll_summaries" (
  "id" SERIAL PRIMARY KEY,
  "payroll_period_id" INT NOT NULL,
  "total_pay" decimal(15,2) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" INT NOT NULL,
  "updated_by" INT NOT NULL
);


CREATE TABLE "audit_logs" (
  "id" SERIAL PRIMARY KEY,
  "request_id" SERIAL NOT NULL, -- Sekarang auto-increment
  "table_name" varchar(50) NOT NULL,
  "record_id" INT NOT NULL,
  "action" action NOT NULL,
  "old_data" jsonb,
  "new_data" jsonb,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "created_by" INT NOT NULL,
  "ip_address" varchar(45) NOT NULL
);


-- Indexes
CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_username" ON "users" ("username");
CREATE INDEX IF NOT EXISTS "idx_payroll_periods_dates" ON "payroll_periods" ("start_date", "end_date");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_attendances_unique" ON "attendances" ("user_id", "payroll_period_id", "attendance_date");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_overtimes_unique" ON "overtimes" ("user_id", "payroll_period_id", "overtime_date");
CREATE INDEX IF NOT EXISTS "idx_reimbursements_period" ON "reimbursements" ("user_id", "payroll_period_id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_payslips_unique" ON "payslips" ("user_id", "payroll_period_id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_payroll_summaries_period" ON "payroll_summaries" ("payroll_period_id");
CREATE INDEX IF NOT EXISTS "idx_audit_logs_record" ON "audit_logs" ("table_name", "record_id");
CREATE INDEX IF NOT EXISTS "idx_audit_logs_created_at" ON "audit_logs" ("created_at");


-- Foreign Keys
ALTER TABLE "payroll_periods" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "payroll_periods" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "attendances" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "attendances" ADD FOREIGN KEY ("payroll_period_id") REFERENCES "payroll_periods" ("id");
ALTER TABLE "attendances" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "attendances" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "overtimes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "overtimes" ADD FOREIGN KEY ("payroll_period_id") REFERENCES "payroll_periods" ("id");
ALTER TABLE "overtimes" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "overtimes" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "reimbursements" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "reimbursements" ADD FOREIGN KEY ("payroll_period_id") REFERENCES "payroll_periods" ("id");
ALTER TABLE "reimbursements" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "reimbursements" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "payslips" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "payslips" ADD FOREIGN KEY ("payroll_period_id") REFERENCES "payroll_periods" ("id");
ALTER TABLE "payslips" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "payslips" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "payroll_summaries" ADD FOREIGN KEY ("payroll_period_id") REFERENCES "payroll_periods" ("id");
ALTER TABLE "payroll_summaries" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "payroll_summaries" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
