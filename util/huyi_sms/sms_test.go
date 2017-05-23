package huyi_sms

import "testing"

func TestSend(t *testing.T) {
	s:=NewSms(&Config{ApiId:"C64911527",ApiKey:"e08477403832118b275aa0543688f63c"})
	s.Send("15828017237","123")
}
