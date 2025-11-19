import dotenv from "dotenv";
import { createClient } from "redis";

dotenv.config();

async function main() {
    const client = createClient({
        socket: {
            host: process.env.REDIS_HOST,
            port: process.env.REDIS_PORT
        }
    });

    client.on("error", console.error);
    await client.connect();

    const key = "user:john.doe@example.com";

    await client.hSet(key, {
        email: "quandale.dingle@example.com",
        first_name: "Quandale",
        last_name: "Dingle"
    });

    const user = await client.hGetAll(key);
    console.log("User:", user);

    await client.quit();
}

main();
