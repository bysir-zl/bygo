package mp

import "testing"

func TestCreateMenu(t *testing.T) {
	err := DeleteMenu()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
