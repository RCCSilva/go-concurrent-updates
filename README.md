# Go Concurrent Updates

Code based on [Rafael Ponte's repository](https://github.com/rafaelpontezup/preventing-lost-update-racecondition).

## Commands

- Execute all tests

```bash
go test ./cmd
```

- Execute individual tests

```bash
go test ./cmd -run TestUpdatesWithLock
```