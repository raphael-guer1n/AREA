## CREATE

```redis
HSET user:quandale.dingle@example.com \
  firs_tname "Quandale" \
  last_name "Dingle" \
  email "quandale.dingle@example.com"
```

## READ

```redis
HGETALL user:quandale.dingle@example.com
```

## UPDATE

```redis
HSET user:quandale.dingle@example.com \
  first_name "updated firs_tname" \
  last_name "updated last_name"
```

## DELETE

```redis
DEL user:quandale.dingle@example.com
```
