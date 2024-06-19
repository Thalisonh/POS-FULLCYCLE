package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	filename := "cotacao.txt"

	dolar := Do(ctx, "GET", "http://localhost:8080/cotacao")

	_, err := readFile(filename)
	if err != nil {
		err := createFile(filename)
		if err != nil {
			return
		}
	}

	writeFile(filename, dolar.Dolar)
}
func writeFile(name string, dolar string) {
	file, err := readFile(name)
	if err != nil {
		return
	}
	defer file.Close()

	fmt.Println(dolar)
	_, err = file.Write([]byte(dolar + "\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func readFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func createFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func Do(ctx context.Context, method, url string) Dolar {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	dolar := Dolar{}
	err = json.Unmarshal(body, &dolar)
	if err != nil {
		fmt.Println(err)
	}

	return dolar
}

type Dolar struct {
	Dolar string `json:"dolar"`
}
