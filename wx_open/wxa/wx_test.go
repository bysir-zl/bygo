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
	x, err := GetAppCode("fMvkYCehn8D_OaKbpMrFYES4nRwG8s4Azx0L8d2Yj9-0algp0sJ-9UCtBRFUY7I-W320Ddozre8uXTece2S4LarLUSUCpglJyB6KfS1Tg23Tcbii5KSPC1K2cv_yqN5YSZOdADDCJY", "123", "pages/index/index", 0, false, Color{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(x)
}
