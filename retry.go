package velix

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

const maxRetries = 3

var retryableStatus = map[int]bool{429: true, 503: true}

func withRetry(ctx context.Context, fn func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error
	backoff := 500 * time.Millisecond

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
		}

		resp, err = fn()
		if err != nil {
			return nil, err
		}
		if !retryableStatus[resp.StatusCode] {
			return resp, nil
		}
		// última tentativa: devolve a resposta (429/503) intacta para o
		// chamador classificar o erro em vez de fechar o body e forçar
		// falha de leitura em velix.go.
		if attempt == maxRetries {
			return resp, nil
		}
		if wait := retryAfterDuration(resp); wait > 0 {
			backoff = wait
		}
		resp.Body.Close()
	}
	return resp, err
}

// retryAfterDuration lê o header Retry-After (segundos ou HTTP-date) da
// resposta, retornando 0 se ausente ou inválido.
func retryAfterDuration(resp *http.Response) time.Duration {
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil {
		if secs < 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}
