package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

func main() {

	var second int
	var minute int

	var rootCmd = &cobra.Command{
		Use:   "transmission",
		Short: "CLI for download transmission picture",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("second", second)
			fmt.Println("minute", minute)

			// 1秒钟2个
			var secondLimit = rate.NewLimiter(Per(second, 1*time.Second), 1)
			//60秒100个
			var minuteLimit = rate.NewLimiter(Per(minute, 1*time.Minute), 20)
			var limiter = MultiLimiter(secondLimit, minuteLimit)
			url := fmt.Sprintf("%s%s", "http://sdkshxj.scdl.cn:31200", "/sdxj-store/HZSH0010F18700513/1_255/2023_9/3_16_3_10.jpg")

			ctx := context.Background()

			fmt.Println("=========")
			for {
				if err := limiter.Wait(ctx); err != nil {
					fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
					go downloadPicture(url)
				}
			}
		},
	}

	rootCmd.Flags().IntVarP(&second, "second", "s", 5, "count per second")
	rootCmd.Flags().IntVarP(&minute, "minute", "m", 200, "count per minute")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
}

func downloadPicture(url string) ([]byte, error) {
	resp, err := (&http.Client{
		Timeout: 10 * time.Second,
	}).Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Println("request too many", time.Now().Format("2006-01-02 15:04:05"))
			// todo 限流重试
		}
		return nil, fmt.Errorf("downloadPicture unexpected status code: %d", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return content, nil

}
