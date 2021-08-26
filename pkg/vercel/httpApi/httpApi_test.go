package httpApi_test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
	"gotest.tools/assert"
)

func testHttpRequest(t *testing.T, wg *sync.WaitGroup, api httpApi.API) {
	defer wg.Done()

	res, _ := api.Request(http.MethodGet, "/v8/projects", nil)

	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestHttpApiRateLimiter(t *testing.T) {
	api := httpApi.New("")

	var wg sync.WaitGroup

	for i := 0; i < 300; i++ {
		wg.Add(1)
		go testHttpRequest(t, &wg, api)
	}

	fmt.Println("Main: Waiting for workers to finish")
	wg.Wait()
	fmt.Println("Main: Completed")
}
