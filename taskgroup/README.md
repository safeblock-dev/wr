### taskgroup

The **taskgroup** package offers the following features:

- **TaskGroup**: A structure that allows you to run multiple functions (actors) concurrently, with the ability to interrupt all actors when one of them returns.
- **Add**: A method to add actors to the TaskGroup. Each actor has an execute function and an interrupt function, ensuring they can be preemptively stopped.
- **AddContext**: A method to add actors with context support, allowing for more fine-grained control over task cancellation and interruption.
- **Run**: A method to execute all actors in the TaskGroup concurrently and handle their interruption and error propagation.

For more details, you can visit the taskgroup [examples](example/taskgroup/main.go).