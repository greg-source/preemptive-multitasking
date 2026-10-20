# Preemptive Multitasking in Go

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)
![Topic](https://img.shields.io/badge/topic-Go%20scheduler-blue)

Companion code for the video **[Preemptive Multitasking](https://www.youtube.com/watch?v=WPtOZAf5arA)**.

The video explains how the Go scheduler decides when to pause one goroutine and run another: how it used to work, how asynchronous preemption fixed its weak spots, and what you can still observe when it is turned off. Every example here is a small, runnable program from the video, so you can reproduce each experiment yourself.

## The idea in 30 seconds

Goroutines are scheduled by the Go runtime, not by the OS. For a long time that scheduling was **cooperative**: a goroutine could only be paused at certain safe points, mostly function calls. A goroutine stuck in a loop with no function calls never reached one, so it could hog its thread forever.

Since **Go 1.14** the runtime also does **asynchronous preemption**. It sends a signal to the thread running a long goroutine and interrupts it at a safe point, even in the middle of a tight loop. This repo shows what the old behavior looked like, and why it mattered, by switching the new behavior off.

## Examples

| Folder | What it shows | What you will see with async preemption **off** |
| --- | --- | --- |
| [`ex1_simple`](ex1_simple/main.go) | One goroutine spins in `for {}` with no function calls. `GOMAXPROCS(1)` leaves a single thread for everyone. | `Goroutine 1 starting` is printed, then the program hangs. Goroutine 2 and `main` never get to run again. |
| [`ex2_simple`](ex2_simple/main.go) | The same loop, but with a function call (`fmt.Sprint()`) inside it. | The function call gives the scheduler a safe point. Goroutine 2 prints every second and `main` exits after 3 seconds. |
| [`ex3_stw`](ex3_stw/main.go) | A CPU-bound loop with no function calls runs next to workers that allocate memory, and a trace is recorded. | The garbage collector needs a stop-the-world pause, but cannot stop the busy loop. Every other goroutine is held up until the loop finishes (about 2 seconds). |

## Running the examples

You need [Go 1.23 or newer](https://go.dev/dl/).

Async preemption is on by default, so the experiments need it switched off with `GODEBUG=asyncpreemptoff=1`.

```bash
GODEBUG=asyncpreemptoff=1 go run ./ex1_simple
GODEBUG=asyncpreemptoff=1 go run ./ex2_simple
GODEBUG=asyncpreemptoff=1 go run ./ex3_stw
```

Now run the same program **without** the variable and compare. This is the main point of the video.

```bash
go run ./ex1_simple
```

> `ex1_simple` never stops when async preemption is off. Press `Ctrl+C` to quit.

### Inspecting the trace (`ex3_stw`)

`ex3_stw` writes a `trace.out` file. Open it with the Go trace viewer and look at the goroutine and GC timelines, with and without the `GODEBUG` flag:

```bash
go tool trace trace.out
```

## Repository layout

```
.
├── ex1_simple/   # tight loop, no function calls
├── ex2_simple/   # tight loop with a function call
├── ex3_stw/      # stop-the-world pause caused by a non-preemptible loop
└── go.mod
```

## Feedback

If you found the video useful, a like and a star on this repo help a lot. Questions and suggestions are welcome in [Issues](../../issues).
