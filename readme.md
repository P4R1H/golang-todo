# Learning backend development with Go.

Key insights:
1) Can't use wildcards in path, need to extract each information ourselves.
ex
```
GET /lists/{listname}/tasks/{tasknum} X
GET /lists/
-> path = strings.trimPrefix(r.URL.path, '/lists')
-> parts = strings.split(path, '/')
Now access through parts[i] to gain information and serve endpoint
```