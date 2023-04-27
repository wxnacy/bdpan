package bdpan

import (
	"fmt"
	"testing"

	"github.com/wxnacy/gotool"
)

func TestErr1(t *testing.T) {
	m := make(map[string]interface{}, 0)
	m["errno"] = -9
	errResp := &ErrorResponse{}
	gotool.MapConverForInterface(m, errResp)
	if errResp.Err() != ErrPathNotFound {
		t.Error(errResp.Err())
	}
}

func TestErr2(t *testing.T) {
	m := make(map[string]interface{}, 0)
	m["error"] = "1234"
	m["error_description"] = "desc"
	errResp := &ErrorResponse{}
	gotool.MapConverForInterface(m, errResp)
	err := fmt.Sprintf("%s[%s]", m["error"], m["error_description"])
	if errResp.Error() != err {
		t.Errorf("%v != %v", errResp.Err(), err)
	}
}

func TestErr3(t *testing.T) {
	m := make(map[string]interface{}, 0)
	m["error_code"] = 1234
	m["error_msg"] = "desc"
	errResp := &ErrorResponse{}
	gotool.MapConverForInterface(m, errResp)
	err := fmt.Sprintf("%d[%s]", m["error_code"], m["error_msg"])
	if errResp.Error() != err {
		t.Errorf("%v != %v", errResp.Err(), err)
	}
}

func TestResponseErr(t *testing.T) {
	m := make(map[string]interface{}, 0)
	m["errno"] = -9
	errResp := &Response{}
	gotool.MapConverForInterface(m, errResp)
	if errResp.Err() != ErrPathNotFound {
		t.Error(errResp.Err())
	}
}

func TestRespIsError(t *testing.T) {
	m := make(map[string]interface{}, 0)
	m["errno"] = -9
	errResp := &Response{}
	gotool.MapConverForInterface(m, errResp)
	if !errResp.IsError() {
		t.Errorf("resp %#v is error", errResp)
	}
}
