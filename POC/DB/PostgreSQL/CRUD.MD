## CREATE

```sql
INSERT INTO users (firs_tname, last_name, email)
VALUES ("Quandale", "Dingle", "quandale.dingle@example.com");
```

## READ

```sql
SELECT *
FROM users
WHERE first_name = "Quandale"
  AND last_name = "Dingle"
  AND email = "quandale.dingle@example.com";
```

## UPDATE

```sql
UPDATE users
SET first_name = "updated firs_tname",
    last_name = "updated last_name"
WHERE email = "quandale.dingle@example.com";
```

## DELETE

```sql
DELETE FROM users
WHERE email = "quandale.dingle@example.com";
```
