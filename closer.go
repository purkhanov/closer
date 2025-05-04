package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Func func(ctx context.Context) error

type closer struct {
	mu    sync.Mutex
	funcs []Func
}

var (
	instance *closer
	once     sync.Once
)

// GetInstance возвращает единственный экземпляр Closer
func GetInstance() *closer {
	once.Do(func() {
		instance = &closer{
			funcs: make([]Func, 0),
		}
	})

	return instance
}

func (c *closer) Add(f Func) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var msgs = make([]string, 0, len(c.funcs))
	var complete = make(chan struct{}, 1)

	go func() {
		for _, f := range c.funcs {
			if err := f(ctx); err != nil {
				msgs = append(msgs, fmt.Sprintf("[!] %v", err))
			}
		}

		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown canceled: %w", ctx.Err())
	}

	if len(msgs) > 0 {
		return fmt.Errorf("shutdown finshed with error(s): \n%s", strings.Join(msgs, "\n"))
	}

	return nil
}
