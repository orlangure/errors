# Annotated errors

```
❯ go run ./example
The original error is created outside of our app:
internal function call failed; external package call failed; EOF
map[foo:bar user:joseph] [main/main.go:30:30 main/main.go:37:37]

The original error is created by us, without any reason:
action on an id failed; operation on a resource failed; this is expected
map[id:42 resource:container root:cause] [main/main.go:49:49 main/main.go:56:56 main/main.go:62:62]

```
