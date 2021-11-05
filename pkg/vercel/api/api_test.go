package api_test

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/api"
	"gotest.tools/assert"
)

func TestApiRequest(t *testing.T) {
	client := api.New(os.Getenv("VERCEL_TOKEN"))

	res, _ := client.Request(http.MethodGet, "/v8/projects", nil, nil)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestApiRateLimiter(t *testing.T) {
	client := api.New(os.Getenv("VERCEL_TOKEN"))

	var wg sync.WaitGroup

	for i := 0; i < 300; i++ {
		wg.Add(1)
		go testHttpRequest(t, &wg, client)
	}

	fmt.Println("Main: Waiting for workers to finish")
	wg.Wait()
	fmt.Println("Main: Completed")
}

func testHttpRequest(t *testing.T, wg *sync.WaitGroup, client *api.Api) {
	defer wg.Done()

	res, _ := client.Request(http.MethodGet, "/v8/projects", nil, nil)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}
