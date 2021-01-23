# Habitz webapp

## Technology

- SQLite - database
- Golang - backend
- HTML/JS - frontend w/ golang templates

## Database

### Users

| Column | Type   | comment                     |
| ------ | ------ | --------------------------- |
| `name` | string | Performer of habit creation |

### Habit Template


| Column    | Type   | comment                     |
| --------- | ------ | --------------------------- |
| `name`    | string | User name                   |
| `weekday` | string | Performer of habit creation |
| `habit`   | string | Habit to create             |

Composite key (`name`, `weekday`, `habit`)

### History

| Column        | Type     | comment                      |
| ------------- | -------- | ---------------------------- |
| `id`          | int      | Auto increasing id, PK       |
| `name`        | string   | User name                    |
| `weekday`     | string   | Performer of habit creation  |
| `date`        | datetime | Habit for this specific date |
| `habit`       | string   | Habit to create              |
| `complete`    | boolean  | Done or not                  |
| `complete_at` | datetime | When it was completed        |

## Webservice API

| Endpoint            | Verb   | comment                          |
| ------------------- | ------ | -------------------------------- |
| `/api/habitz/users` | GET    | Names of habit creators          |
| `/api/habitz/`      | POST   | Create a new habit template      |
| `/api/habitz/`      | DELETE | Remove a habit from the template |
| `/api/habitz/today` | GET    | Load todays habit history        |
| `/api/habitz/today` | UPDATE | Update todays habit history      |

### `/api/habitz/users`

**GET**

List users if there are any, just names. Mostly for form completion.

```json
["John Doh", "Mary Doh"]
```

### `/api/habitz/`

**POST**

If `John Doh` does not exist, a new user will be created. The habit will be added to the users template. If the weekday is today, also add it to todays history. 

```json
{
    "name": "John Doh",
    "habit": "Walk 1 mile a day",
    "weekdays": ["Monday","Friday"],
}
```

### `/api/habitz/`

**DELETE**

Remove the habit from the users template. If it's removed from todays template, remove it from todays history. 

```json
{
    "name": "John Doh",
    "habit": "Walk 1 mile",
    "weekdays": ["Monday"],
}
```

### `/api/habitz/today`

**GET**

Retrieve the habits to form today for all users. If nothing exists for today, create new entries in history with habitz from template and mark incomplete. 

```json
[
    {
        "user": "John Doh",
        "habitz": [
            {"id": 5, "habit": "Walk 1 mile", "complete": false},
            {"id": 3, "habit": "Drink 10 beers", "complete": false},
            {"id": 6, "habit": "Floss", "complete": false}
        ]
    },
    {
        "user": "Mary Doh",
        "habitz": [
            {"id": 5, "habit": "Play golf", "complete": false},
            {"id": 9, "habit": "Make chocolate balls", "complete": false},
            {"id": 12, "habit": "Write a poem", "complete": false}
        ]
    },
]
```

**UPDATE**

Incomplete updates are OK. 

```json
[
    {
        "user": "John Doh",
        "habitz": [
            {"id": 5, "complete": true}
        ]
    },
    {
        "user": "Mary Doh",
        "habitz": [
            {"id": 9, "complete": false}
        ]
    },
]
```

## Web App

| Endpoint | Verb | comment                       |
| -------- | ---- | ----------------------------- |
| `/`      | GET  | HTML for todays habit history |
| `/new`   | GET  | Form for new habit template   |

