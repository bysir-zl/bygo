package wx_thrid_party

import "testing"

func TestSaveVerifyTicket(t *testing.T) {
    err:=SaveVerifyTicket("tttt")
    if err != nil {
        t.Error(err)
    }
    t.Log("ok")
}

func TestGetLastVerifyTicket(t *testing.T) {
    x,ok:=GetLastVerifyTicket()
    if !ok{
        t.Error(ok)
    }

    t.Log(x)
}

// 68403 ns/op
func BenchmarkGetLastVerifyTicket(b *testing.B) {
    for i := 0; i < b.N; i++ {
        x,_:=GetLastVerifyTicket()
        if x!="tttt"{
            panic(x)
        }
    }
}
// 3 690 215 ns/op
func BenchmarkSaveVerifyTicket(b *testing.B) {
    for i := 0; i < b.N; i++ {
        err:=SaveVerifyTicket("tttt")
        if err != nil {
            panic(err)
        }
    }
}