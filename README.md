# VocabKing API

| Header             | Required   | Description                                                                                  | Example                         |
|--------------------|------------|----------------------------------------------------------------------------------------------|---------------------------------|
| `Accept`           | Yes        | Explicitly indicate to the server that you want JSON format responses. Use example verbatim. | `Accept: application/json`      |
| `Authorization`    | For `/api` | Bearer token for authenticated endpoints                                                     | `Authorization: Bearer abc.xyz` |
| `X-Client-Name`    | No         | Name of the client used to communicate with the server                                       | `X-Client-Name: Android-App`    |
| `X-Client-Version` | No         | Version of the client used to communicate with the server                                    | `X-Client-Version: v1.0.0`      |

All request bodies **must be** in JSON format and should follow this basic structure:

``` json
{
  "data": {
    ...
  }
}
```

`data` can be safely omitted if it is empty or has no required fields.

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

Note that it is not guaranteed that `data.ops` and `data.stack` will exist and may be excluded in non-debug operation.

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
    "username": "vocabking_0723",
    "password": "V0C4BK1NGKN0W5M4NYW0RD54NDH45M4NYFR1END5"
  }
}
```

| Key        | Required | Type   | Description      |
|------------|----------|--------|------------------|
| `username` | Yes      | String | Account username |
| `password` | Yes      | String | Account password |

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

| Key      | Type   | Description             |
|----------|--------|-------------------------|
| `expiry` | Number | Token expiry in seconds |
| `token`  | String | Authentication token    |

### POST `/auth/refresh`

<span style="color: red">**&mdash; Not Implemented &mdash;**</span>

Use this endpoint to refresh an existing authentication token's expiry.

Replace your token with the token given in the response data.

Note that you may only refresh a token within **30 seconds** of its expiry.
HTTP Bad Request will be returned otherwise.

**Request data:**

``` json
{
  "data": {
    "token": "qwertyuiopasdfghjklzxcvbnm.qwertyuiopasdfghjklzxcvbnm"
  }
}
```

| Key     | Required | Type   | Description                          |
|---------|----------|--------|--------------------------------------|
| `token` | Yes      | String | Authentication token to be refreshed |

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

| Key      | Type   | Description                  |
|----------|--------|------------------------------|
| `expiry` | Number | Token expiry in milliseconds |
| `token`  | String | Authentication token         |

### POST `/api/picture`

<span style="color: red">**&mdash; Not Implemented &mdash;**</span>

### POST `/api/word`

<span style="color: red">**&mdash; Not Implemented &mdash;**</span>
