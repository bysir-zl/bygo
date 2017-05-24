package huyi_sms

import "testing"

func TestSend(t *testing.T) {
	s := NewSms(&Config{ApiId: " ", ApiKey: " "})
	err := s.Send("15828017237", "您的验证码是：123456。请不要把验证码泄露给其他人。")
	t.Error(err)
}
