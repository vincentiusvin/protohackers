My solutions to https://protohackers.com/

Done:

- [x] 0 - Smoke Test
- [x] 1 - Prime Time
- [x] 2 - Means to an End
- [x] 3 - Budget Chat
- [x] 4 - Unusual Database Program
- [x] 5 - Mob in the Middle
- [x] 6 - Speed Daemon
- [x] 7 - Line Reversal
- [ ] 8 - Insecure Sockets Layer

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

# Notes

## 0

Simple echo, wrote three different ways to do it

## 1

This is actually harder than a bunch of the challanges suceeding it. You have to work with json over tcp and take care of big numbers.

## 2

Simple tcp. Session is stored based on connection.

## 3

Also simple tcp. Data is shared between multiple connections so you gotta do synchronziation.

## 4

First time working with UDP in go. Still simple outside of trying to find which method to use.

## 5

Still simple, just a TCP proxy. Needed to do a little bit of parsing.

## 6

Okay this is starting to get difficult. Lots of domain logic here. Gotta do parsing too.

I used something a-la functional parsing when dealing with the bytes. In hindsight this is not a very good idea.

Simple state-based parsing might be a better idea.

## 7

This one is more of an infrastructure challenge instead of a domain logic one.

Lowered the retransmission timer from 3 secs to 1 sec because the packet loss doesn't seem to be at 25%? Kept getting the 60s timeout.

Might be a skill issue on my part.
