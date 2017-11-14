package artisan

import "testing"

func TestJson2Go(t *testing.T) {
	err := Json2Go()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
