# preemptive-multitasking

Use `GODEBUG=asyncpreemptoff=1` to disable asynchronous preemption.

```code
GODEBUG=asyncpreemptoff=1 go run ex1_simple/main.go
```