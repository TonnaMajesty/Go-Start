package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
)

var (
	Blogger   = find("1227368500")
	Weibo = find("H3GIgngon")
)

type Result string
type Find func(ctx context.Context, query string) (Result, error)

func find(kind string) Find {
	return func(_ context.Context, query string) (Result, error) {
		return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
	}
}

func main() {
	SinaWeibo := func(ctx context.Context, query string) ([]Result, error) {
		g, ctx := errgroup.WithContext(ctx)

		finds := []Find{Blogger, Weibo}
		results := make([]Result, len(finds))
		for i, find := range finds {
			// errgroup 小心闭包的问题
			i, find := i, find // https://golang.org/doc/faq#closures_and_goroutines
			g.Go(func() error {
				result, err := find(ctx, query)
				if err == nil {
					results[i] = result
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return nil, err
		}
		return results, nil
	}

	results, err := SinaWeibo(context.Background(), "https://weibo.com/1227368500/H3GIgngon")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}

}