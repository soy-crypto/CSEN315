package roundrobin

import (
	"fmt"
	"net/url"
	"sync"
)

func main() {

	//Test Roundrobin
	resources := []*url.URL{
		{Host: "127.0.0.1"},
		{Host: "127.0.0.2"},
		{Host: "127.0.0.3"},
		{Host: "127.0.0.4"},
		{Host: "127.0.0.5"},
		{Host: "127.0.0.6"},
		{Host: "127.0.0.7"},
		{Host: "127.0.0.8"},
		{Host: "127.0.0.9"},
		{Host: "127.0.0.10"},
	}

	for i := 1; i < len(resources)+1; i++ {
		fmt.Sprintf("RoundRobinSliceOfSize(%d)", i)
		rr, err := New(resources[:i]...)
		if err != nil {
			b.Fatal(err)
		}
		wg := &sync.WaitGroup{}
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				rr.Next()
			}()
		}
		wg.Wait()

	} // for

}
