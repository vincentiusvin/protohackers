My solutions to https://protohackers.com/

# Running

To deploy, run

```bash
./deploy.sh 0_echo/0_echo.go
```

This will compile and scp the executable to your server. You can run it there.

```bash
./0_echo
```

# Testing

Not very many tests since the assignment itself is one huge test suite.

Tests are only for tricky logic.

```bash
go test ./...
```
