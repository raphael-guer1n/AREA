from pydantic import BaseModel

class User(BaseModel):
    userId: int = 0
    name: str = ""
    email: str = ""
    password: str = ""

class UserDB:
    users : dict[int, User]
    counter : int

    def __init__(self):
        self.users = {}
        self.counter = 1

        self.users[self.counter] = User(
            userId=self.counter,
            name="name",
            email="email",
            password="password",
        )

    def get_users(self, userId : int) -> User:
        return self.users[userId]

    def add_user(self, user : User) -> int:
        self.counter += 1
        self.users[self.counter] = user
        return self.counter

    def delete_user(self, userId : int) -> None:
        del self.users[userId]

    def update_user(self, userId : int, user : User) -> None:
        self.users[userId] = user