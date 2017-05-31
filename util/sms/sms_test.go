package sms

import "testing"

func TestSend(t *testing.T) {
	s := NewHuyiSms(" "," ")
	err := s.Send("15828017237", "您的验证码是：123456。请不要把验证码泄露给其他人。")
	t.Error(err)
}

func TestDaiyiSend(t *testing.T) {
	s := NewDaiyi(" "," ")
	err := s.Send("15828017237", "您的验证码为@，请勿告诉他人。","江苏骠客娱乐")
	t.Error(err)
}