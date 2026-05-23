import postgres from "postgres";

const databaseUrl =
  process.env.DATABASE_URL ??
  "postgres://lang_war:lang_war_password@localhost:5432/lang_war?sslmode=disable";

const sql = postgres(databaseUrl);

export default sql;
