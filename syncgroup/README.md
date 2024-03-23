### syncgroup

The **syncgroup** package provides an enhanced version of the standard sync.WaitGroup with additional features for managing concurrent tasks. Key features include:

- **Panic Handling:** Customizable panic handling mechanism to capture and handle panics that occur within goroutines.
- **Flexible Task Submission:** Ability to submit tasks to the group and automatically manage their completion.
- **Wait for Completion:** A Wait method that blocks until all submitted tasks have completed, ensuring proper synchronization.
- **Error Propagation:** Support for propagating errors from tasks to the caller, allowing for more robust error handling.

This package offers a more comprehensive approach to managing concurrency, making it easier to handle errors and panics in a controlled manner.