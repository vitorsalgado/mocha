package mocha

import (
	"github.com/vitorsalgado/mocha/matchers"
	"log"
	"net/http"
	"testing"
)

func TestMocha(t *testing.T) {
	t.Run("should mock request", func(t *testing.T) {
		m := New()
		defer m.Close()

		m.Mock(NewBuilder().
			Header("test", matchers.EqualTo("hello")).
			Res().
			Build())

		req, _ := http.NewRequest(http.MethodGet, m.Server.URL, nil)
		req.Header.Add("test", "hello")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != 201 {
			t.Errorf("expected status code 201. received %d", res.StatusCode)
		}
	})
}
