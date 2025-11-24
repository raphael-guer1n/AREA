import { MongoClient } from "mongodb";
import dotenv from "dotenv";

dotenv.config();

async function main() {
    const client = new MongoClient(process.env.MONGO_URI);
    await client.connect();

    const db = client.db(process.env.MONGO_DB);
    const users = db.collection("users");

    await users.updateOne(
        { email: "quandale.dingle@example.com" },
        { $set: { first_name: "Quandale", last_name: "Dingle" } },
        { upsert: true }
    );

    const allUsers = await users.find().toArray();
    console.log("Users:", allUsers);

    await client.close();
}

main();
