### wrpool

The **wrpool** package provides a flexible and efficient worker pool for managing and executing concurrent tasks in Go. It offers features such as dynamic worker scaling, task queuing, context support for cancellation and deadlines, error handling, and panic recovery. This package is designed to simplify concurrency management and improve resource utilization in Go applications.

- **Concurrent Task Execution:** Manages a pool of goroutines to execute tasks concurrently.
- **Task Queuing:** Supports queuing tasks when all workers are busy.
- **Dynamic Worker Scaling:** Can dynamically adjust the number of active workers based on workload.
- **Context Support:** Allows for context-based cancellation and deadline control.
- **Error Handling:** Provides mechanisms for handling errors from tasks.
- **Panic Recovery:** Can recover from panics in tasks to prevent crashes.
- **Configuration Flexibility:** Offers various options to configure the behavior of the worker pool, such as setting the maximum number of workers, specifying a custom error handler, or setting a panic recovery function.
- **Rate Limiting:** Controls the rate at which tasks are executed, ensuring that the system is not overwhelmed.
- **Resource Management:** Helps in managing resource utilization by limiting the number of concurrent tasks.
