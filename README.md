# Popular concepts in GO for Distributed Systems

## Concurrent Task Scheduler

`go run scheduler.go`

### Problem Statement

Design and implement a Concurrent Task Scheduler that allows tasks to be scheduled for execution with the following requirements:

- The scheduler should accept tasks with a given priority (higher priority tasks should execute first).
- The scheduler should support multiple worker threads executing tasks concurrently.
- The system should allow dynamically adding tasks while it is running.
- The scheduler should ensure fair execution of tasks, preventing starvation of lower-priority tasks.
- Implement the core logic using synchronization primitives (e.g., mutexes, condition variables, or semaphores) to ensure thread safety.
