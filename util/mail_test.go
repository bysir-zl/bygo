package util

import "testing"

func TestNewMail(t *testing.T) {
	mail := NewMail("120735581@qq.com", "peerdmnoqirqbiaa", "smtp.qq.com")

	err:=mail.Send("sbsb","sbsbsb",[]string{"yangzefeng@kuaifazs.com"})
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkMail(b *testing.B) {
	mail := NewMail("120735581@qq.com", "peerdmnoqirqbiaa", "smtp.qq.com")

	for i := 0; i < b.N; i++ {
		 mail.Send("sbsb","sbsbsb",[]string{"yangzefeng@kuaifazs.com"})

	}
}
