## CREATE

```js
db.users.insertOne({
  firs_tname: "Quandale",
  last_name: "Dingle",
  email: "quandale.dingle@example.com"
})
```

## READ

```js
db.users.findOne({
    first_name: "Quandale",
    last_name: "Dingle",
    email: "quandale.dingle@example.com"
})
```

## UPDATE

```js
db.Users.updateOne(
  { email: "quandale.dingle@example.com" },
  { $set: { first_name: "updated firs_tname", last_name: "updated last_name" } }
)
```

## DELETE

```js
db.Users.deleteOne({
    email: "quandale.dingle@example.com"
})
```
