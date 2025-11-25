from pydantic import BaseModel

class UserCreate(BaseModel):
    name: str = ""
    email: str = ""
    password: str = ""

class User(UserCreate):
    userId: int = 0

class UserDB:
    users: dict[int, User]
    counter: int

    def __init__(self):
        self.users = {}
        self.counter = 0

        self.users[self.counter] = User(
            userId=self.counter,
            name="name",
            email="email",
            password="password",
        )

    def get_users(self, userId: int) -> User:
        return self.users[userId]

    def add_user(self, user: UserCreate) -> int:
        self.counter += 1
        new_user = User(userId=self.counter, **user.model_dump())
        self.users[self.counter] = new_user
        return self.counter

    def delete_user(self, userId: int) -> None:
        del self.users[userId]

    def update_user(self, userId: int, user: UserCreate) -> None:
        if userId not in self.users:
            raise KeyError
        updated = User(
            userId=userId,
            **user.model_dump(),
        )
        self.users[userId] = updated