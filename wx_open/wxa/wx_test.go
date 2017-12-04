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
	x, err := GetAppCode("eE7v52HriTk2V-6njRedn7PdkxjO7wslaBJb0y20ev27aTPFuWoUMVst90hNlLn7if5yn460uV2rCC5_IQSSKURwNgc7Eb1_0L-E2lyUn5_EKtjqU1MOcWCCC0WCq8QYZDWhAHDOSL", "?orderid=1", "pages/start/index", 0, false, Color{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(x)
}

func TestGetMonthRetain(t *testing.T) {
	x, err := GetMonthRetain("3CWYQOE52KfMWPMxO3nhyv8iYZydDRc2pzKOcp7LeJpzT0I6vZQIYiI76i2CAj2Cr4srV0VzmHHkf1lwCESm6O1UA8d1cugN-Kf1mcPLH1gMsjg5dxQ2FEa7Kzcl9XDxRVScAHDYCZ","201711")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(x)
}

func TestMonthVisitRsp(t *testing.T) {
	x, err := GetMonthVisit("3CWYQOE52KfMWPMxO3nhyv8iYZydDRc2pzKOcp7LeJpzT0I6vZQIYiI76i2CAj2Cr4srV0VzmHHkf1lwCESm6O1UA8d1cugN-Kf1mcPLH1gMsjg5dxQ2FEa7Kzcl9XDxRVScAHDYCZ","201711")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(x)
}
