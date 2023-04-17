package retry

import (
	"time"
)

type Retrier interface {
	PassedIterationCount() int
	SetDelay(duration time.Duration)
	Stop()
}

type retrier struct {
	iteration int
	delay     *time.Duration
	stopped   bool

	execute func() error
	inspect func(retrier Retrier, err error)
}

// PassedIterationCount returns number of passed iterations.
func (r *retrier) PassedIterationCount() int {
	return r.iteration
}

// SetDelay sets delay for next iteration.
func (r *retrier) SetDelay(duration time.Duration) {
	r.delay = &duration
}

// Stop preventing further executions.
func (r *retrier) Stop() {
	r.stopped = true
}

type Option func(r *retrier)

func Execute(fn func() error) Option {
	return func(r *retrier) {
		r.execute = fn
	}
}

func Inspect(fn func(r Retrier, err error)) Option {
	return func(r *retrier) {
		r.inspect = fn
	}
}

func Retry(options ...Option) error {
	r := &retrier{}
	for _, option := range options {
		option(r)
	}

	if r.execute == nil {
		return nil
	}

	for {
		r.iteration++
		r.delay = nil

		err := r.execute()
		if err == nil {
			return nil
		}

		if r.inspect != nil {
			r.inspect(r, err)
		}

		if r.stopped {
			return err
		}

		if r.delay != nil {
			<-time.After(*r.delay)
		}
	}
}
