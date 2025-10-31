package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"strings"
	"time"
)

func main() {
	endpoint := "http://localhost:8080/"

	fmt.Println("Введите длинный URL")
	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\n")

	client := resty.New()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := client.R().SetContext(ctx).SetHeader("Content-Type", "text/plain").SetBody(long).Post(endpoint)

	if err != nil {
		panic(err)
	}

	fmt.Println("Статус-код ", resp.Status())
	fmt.Println(resp.String())
}
