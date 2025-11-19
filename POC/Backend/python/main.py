from fastapi import FastAPI, HTTPException
import users

app = FastAPI()
userDB : users.UserDB = users.UserDB()

@app.get("/user/{userId}")
def get_user(userId: int):
    try:
        user = userDB.get_users(userId)
        return {"user": user}
    except:
        raise HTTPException(status_code=404, detail="User not found")

@app.post("/user")
def add_user(user: users.UserCreate):
    try:
        newUserId = userDB.add_user(user)
        return {"userId": newUserId}
    except:
        raise HTTPException(status_code=400, detail="Bad Request")

@app.delete("/user/{userId}")
def delete_user(userId: int):
    try:
        userDB.delete_user(userId)
        return {"userId": userId}
    except:
        raise HTTPException(status_code=404, detail="User not found")


@app.put("/user/{userId}")
def update_user(userId: int, user: users.UserCreate):
    try:
        userDB.update_user(userId, user)
        return {"userId": userId}
    except KeyError:
        raise HTTPException(status_code=404, detail="User not found")
