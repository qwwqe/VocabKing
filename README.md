# VocabKing API

All requests **must include** the header `Accept: application/json`.

All request bodies **must be** JSON and should follow this basic structure:

``` json
{
  "data": {
    ...
  }
}
```

All requests under `/api` **must include** a valid bearer token.

Example:

``` plain
Authorization: Bearer TOKEN
```

If the request is <span style="color: green">**successful**</span> the response will be:

``` json
{
  "result": "ok",
  "data": {
    ...
  }
}
```

If the request is <span style="color: red">**unsuccessful**</span> the response will be:

``` json
{
  "result": "error",
  "data": {
    "error": "",
    "ops": [
      ...
    ],
    "stack": [
      ...
    ]
  }
}
```

Note that it is not guaranteed that `data.ops` and `data.stack` will exist and
may be excluded in non-debug operation.

## Endpoints

| Verb | Path            | Description               |
|------|-----------------|---------------------------|
| POST | `/auth/login`   | Get login token           |
| POST | `/auth/refresh` | Refresh valid login token |
| POST | `/api/picture`  | Save picture              |
| POST | `/api/word`     | Save word                 |

### POST `/auth/login`

Use this endpoint to retrieve an authorization token you can use to access
endpoints under `/api`.

Tokens expire in **five minutes** and may be refreshed within **30 seconds** of expiry.

Failing to refresh your token, you will need to request for a new token using
this endpoint.

**Request data:**

``` json
{
  "data": {
    "username": "username",
    "password": "password"
  }
}
```

| Key        | Description      |
|------------|------------------|
| `username` | Account username |
| `password` | Account password |

**Response data:**

``` json
{
  "result": "ok",
  "data": {
    "expiry": 12345678901234567890,
    "token": "qwertyuiopasdfghjklzxcvbnm.qwertyuiopasdfghjklzxcvbnm"
  }
}
```

| Key      | Description                  |
|----------|------------------------------|
| `expiry` | Token expiry in milliseconds |
| `token`  | Authentication token         |

### POST `/auth/refresh`

<span style="color: red">**&mdash; Not Implemented &mdash;**</span>

Use this endpoint to refresh an existing authentication token's expiry.

Replace your token with the token given in the response data.

Not that you may only refresh a token within **30 seconds** of its expiry.
HTTP Bad Request will be returned otherwise.

**Request data:**

``` json
{
  "data": {
    "token": "qwertyuiopasdfghjklzxcvbnm.qwertyuiopasdfghjklzxcvbnm"
  }
}
```

| Key     | Description                          |
|---------|--------------------------------------|
| `token` | Authentication token to be refreshed |

**Response data:**

``` json
{
  "result": "ok",
  "data": {
    "expiry": 12345678901234567890,
    "token": "qwertyuiopasdfghjklzxcvbnm.qwertyuiopasdfghjklzxcvbnm"
  }
}
```

| Key      | Description                  |
|----------|------------------------------|
| `expiry` | Token expiry in milliseconds |
| `token`  | Authentication token         |

### POST `/api/picture`

<span style="color: red">**&mdash; Not Implemented &mdash;**</span>

### POST `/api/word`

<span style="color: red">**&mdash; Not Implemented &mdash;**</span>
