package wxa

import "testing"

func TestGetCategory(t *testing.T) {
	x, err := GetCategory("eSQMOFH4Y7BNiObdnBbpoXs5_N8FQ3FQuYePgQ5B5E1Gqjj5llvPjQbyciyN0h6LlhgfSPmBbaU-r9FTy1z6FCG-kxKknrLy28K7020MhNT87ERPFTc1q9fHdaFIhMsJIEGbADDQMO")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", x)
}

func TestGetAppCode(t *testing.T) {
	x, err := GetAppCode("JI7Rr1w0mW_DcoBMU07WYbLzIWpa9uuItJg-ojMvDoT7Eh7FFUfZ_0V9QUXwLWrFrVGu4uRU2_e83a82RYXwO4xUk1OeS0Q3cvq2XZBjcnYaxlzPATslyWq7SYFJVGZqGPKjAMDSBW", "123", "pages/index/index", 0, false, Color{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(x)
}
