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

## Rate Limiter Service

`go run rate_limiter.go`

### Problem Statement

Design and implement a rate limiter that controls the number of requests allowed per user within a given time window. The rate limiter should efficiently handle high-concurrency scenarios, ensuring fairness and preventing abuse.

- The rate limiter should allow N requests per user per time window (e.g., 100 requests per minute).
- Requests beyond the allowed limit should be rejected.
- Should support multiple users concurrently.
- The solution should be efficient in memory usage and scalable for high throughput.
- The rate limiter should allow different policies (e.g., fixed window, sliding window, token bucket).
