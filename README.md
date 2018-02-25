# Simple Online Leaderboard

An online leaderboard running on Google App Engine.

## Usage

To post a new score, you can issue an HTTP POST request.

```
$ curl -x POST "https://[PROJECT_ID].appspot.com/api/scores?name=Lambert&score=42&password=12345"
```

Get all scores with this request.

```
$ curl "https://[PROJECT_ID].appspot.com/api/scores?skip=0&take=3"
```

The response will look similar to this. Scores will be sorted from highest to lowest.

```json
{
  "data": [
    {
      "name": "Blinn",
      "score": "53"
    },
    {
      "name": "Lambert",
      "score": "42"
    },
    {
      "name": "Phong",
      "score": "24"
    }
  ]
}
```

The `skip` parameter tells how many records to skip, and `take` the maximum number of records to return.
