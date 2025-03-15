My solutions to https://protohackers.com/

Done:

- [x] 0 - Smoke Test
- [x] 1 - Prime Time
- [x] 2 - Means to an End
- [x] 3 - Budget Chat
- [x] 4 - Unusual Database Program
- [x] 5 - Mob in the Middle
- [x] 6 - Speed Daemon
- [ ] 7 - Line Reversal

# Running

To deploy, run

```bash
./deploy.sh 0_echo/0_echo.go
```

This will compile and scp the executable to your (my) server. You can run it there.

```bash
./0_echo
```

# Testing

Not very many tests since the assignment itself is one huge test suite.

Tests are only for tricky logic.

```bash
go test ./...
```
