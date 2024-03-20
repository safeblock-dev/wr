package wrgroup_test

import (
	"testing"

	"github.com/safeblock-dev/wr/panics"
	"github.com/safeblock-dev/wr/wrgroup"
	"github.com/stretchr/testify/require"
)

func TestWaitGroup(t *testing.T) {
	t.Parallel()

	t.Run("go and wait", func(t *testing.T) {
		t.Parallel()

		var counter int
		wg := wrgroup.New()

		wg.Go(func() {
			counter++
		})

		wg.Wait()
		require.Equal(t, 1, counter, "Counter should be incremented")
	})

	t.Run("with panic", func(t *testing.T) {
		t.Parallel()

		panicHandled := false
		panicHandler := func(_ panics.Recovered) {
			panicHandled = true
		}

		wg := wrgroup.New(wrgroup.PanicHandler(panicHandler))

		wg.Go(func() {
			panic("test panic")
		})

		wg.Wait()
		require.True(t, panicHandled, "Panic should be handled")
	})
}
