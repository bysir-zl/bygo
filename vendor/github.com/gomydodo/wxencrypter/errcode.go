package wxencrypter

import (
	"errors"
)

type WechatErrorCode int

const (
	ValidateSignatureErrorCode WechatErrorCode = -40001 - iota
	ParseXmlErrorCode
	ComputeSignatureErrorCode
	IllegalAesKeyCode
	ValidateAppidErrorCode
	EncryptAESErrorCode
	DecryptAESErrorCode
	IllegalBufferCode
	EncodeBase64ErrorCode
	DecodeBase64ErrorCode
	GenReturnXmlErrorCode
)

var (
	ValidateSignatureError = errors.New("ValidateSignatureError")
	ParseXmlError          = errors.New("ParseXmlError")
	ComputeSignatureError  = errors.New("ComputeSignatureError")
	IllegalAesKey          = errors.New("IllegalAesKey")
	ValidateAppidError     = errors.New("ValidateAppidError")
	EncryptAESError        = errors.New("EncryptAESError")
	DecryptAESError        = errors.New("DecryptAESError")
	IllegalBuffer          = errors.New("IllegalBuffer")
	EncodeBase64Error      = errors.New("EncodeBase64Error")
	DecodeBase64Error      = errors.New("DecodeBase64Error")
	GenReturnXmlError      = errors.New("GenReturnXmlError")
)

func errorToCode(err1 error) (code WechatErrorCode, err error) {
	switch err1 {
	case ValidateSignatureError:
		code = ValidateSignatureErrorCode
	case ParseXmlError:
		code = ParseXmlErrorCode
	case ComputeSignatureError:
		code = ComputeSignatureErrorCode
	case IllegalAesKey:
		code = IllegalAesKeyCode
	case ValidateAppidError:
		code = ValidateAppidErrorCode
	case EncryptAESError:
		code = EncryptAESErrorCode
	case DecryptAESError:
		code = DecryptAESErrorCode
	case IllegalBuffer:
		code = IllegalBufferCode
	case EncodeBase64Error:
		code = EncodeBase64ErrorCode
	case DecodeBase64Error:
		code = DecodeBase64ErrorCode
	case GenReturnXmlError:
		code = GenReturnXmlErrorCode
	default:
		err = errors.New("Not found Code")
	}
	return
}
