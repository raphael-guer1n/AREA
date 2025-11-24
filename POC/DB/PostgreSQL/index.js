import dotenv from "dotenv";
import pkg from "pg";
import fs from "fs";

dotenv.config();
const { Client } = pkg;
const client = new Client();

async function main() {
    await client.connect();

    const schema = fs.readFileSync("./schema.sql").toString();
    await client.query(schema);

    await client.query(
        "INSERT INTO users (email, first_name, last_name) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
        ["quandale.dingle@example.com", "Quandale", "Dingle"]
    );

    const res = await client.query("SELECT * FROM users");
    console.log("Users:", res.rows);

    await client.end();
}

main();
