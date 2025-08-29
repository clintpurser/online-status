# Module online-status 

This module provides a sensor that checks the online status of `https://app.viam.com`. It returns `1` if the site is reachable and returns a `200 OK` status code, and `0` otherwise. The check includes a 2-second timeout.

## Model cdp:online-status:online-status

This model represents an online status sensor. It continuously monitors the reachability of `https://app.viam.com`.

### Readings

The `Readings` function returns a map with a single key `online` and an integer value:

| Name     | Type | Description                               |
|----------|------|-------------------------------------------|
| `online` | int  | `1` if online (200 OK), `0` if offline/error |

#### Example Readings

```json
{
  "online": 1
}
```
