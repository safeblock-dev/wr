package wrtask

import (
	"context"

	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/syncgroup"
)

// TaskGroup manages a collection of concurrent tasks (actors).
// It ensures that when one task completes, all other tasks are interrupted.
// A TaskGroup with no tasks is a valid, empty group.
type TaskGroup struct {
	actors []actor
}

// actor represents a task with an execute function and an interrupt function.
// The interrupt function is called to preemptively stop the execution.
type actor struct {
	execute   func() error
	interrupt func(error)
}

// New creates a new, empty TaskGroup.
func New() TaskGroup {
	return TaskGroup{
		actors: nil,
	}
}

// Add appends a new task (actor) to the TaskGroup.
// The execute function is the task's main logic, and the interrupt function is called
// to stop the task. The interrupt function should cause the execute function to return quickly.
// The error returned by the first task to complete is passed to all interrupt functions and returned by Run.
func (g *TaskGroup) Add(execute func() error, interrupt func(error)) {
	if execute == nil || interrupt == nil {
		panic("execute and interrupt functions must not be nil")
	}

	g.actors = append(g.actors, actor{execute, interrupt})
}

// AddContext adds a task to the TaskGroup that accepts a context.
// The 'execute' function is the main task logic that should respect the provided context,
// particularly its cancellation signal. The 'interrupt' function is called when the task
// needs to be preemptively stopped, which should cause 'execute' to return quickly.
// The 'interrupt' function is also provided with the context and the error that caused the interruption.
func (g *TaskGroup) AddContext(execute func(ctx context.Context) error, interrupt func(context.Context, error)) {
	if execute == nil || interrupt == nil {
		panic("execute and interrupt functions must not be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	g.actors = append(g.actors, actor{
		func() error {
			return execute(ctx)
		},
		func(err error) {
			cancel()
			interrupt(ctx, err)
		},
	})
}

// Run executes all tasks in the group concurrently.
// It waits for the first task to complete, then interrupts all remaining tasks.
// Run blocks until all tasks have exited, and returns the error from the first task to complete.
// If the group is empty, Run returns nil immediately.
func (g *TaskGroup) Run() error {
	if len(g.actors) == 0 {
		return nil
	}

	errors := make(chan error, len(g.actors))
	defer close(errors)

	// Use syncgroup to manage task goroutines and handle panics.
	wg := syncgroup.New(syncgroup.PanicHandler(func(recovered panics.Recovered) {
		errors <- recovered.AsError()
	}))

	for _, a := range g.actors {
		wg.Go(func() {
			func(a actor) {
				errors <- a.execute()
			}(a)
		})
	}

	// Wait for the first task to complete and retrieve its error.
	err := <-errors

	// Interrupt all tasks with the error from the first completed task.
	for _, a := range g.actors {
		a.interrupt(err)
	}

	// Wait for all tasks to complete before returning.
	wg.Wait()

	return err
}

func (g *TaskGroup) Size() int {
	return len(g.actors)
}

func SkipInterrupt() func(error) {
	return func(error) {}
}
