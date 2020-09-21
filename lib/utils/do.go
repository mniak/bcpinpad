package utils

import "context"

func RunWithContext(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	ch := make(chan interface{})
	defer close(ch)
	che := make(chan error)
	defer close(che)
	go func() {
		result, err := fn()
		if err != nil {
			che <- err
		} else {
			ch <- result
		}
	}()
	select {
	case result := <-ch:
		return result, nil
	case err := <-che:
		return byte(0), err
	case <-ctx.Done():
		return byte(0), ctx.Err()
	}
}
