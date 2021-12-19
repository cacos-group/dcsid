package dcsid

import (
	"testing"
)

func TestGenid_NextId(t *testing.T) {

	dcsid := New(&Config{
		DSN:    "admin:admin@tcp(127.0.01:3306)/cacos?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&charset=utf8mb4&loc=Local",
		BizTag: "test1",
		Step:   1000,
	})

	go func() {
		for i := 0; i < 100000; i++ {
			_, err := dcsid.NextId()
			if err != nil {
				t.Error(err)
				return
			}
			//fmt.Println(id)
			//time.Sleep(time.Millisecond)
		}
	}()

	for i := 0; i < 100000; i++ {
		_, err := dcsid.NextId()
		if err != nil {
			t.Error(err)
			return
		}
		//fmt.Println(id)
		//time.Sleep(time.Millisecond)
	}

}
