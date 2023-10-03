package valueobject

import "testing"

func TestParamsContains(t *testing.T) {
	// 创建一组参数
	params := Params{
		NewParam("param1"),
		NewParam("param2"),
		NewParam("param3"),
	}

	// 测试参数是否包含特定名称的参数
	if !params.Contains("param1") {
		t.Errorf("Expected params to contain 'param1'")
	}

	if !params.Contains("param2") {
		t.Errorf("Expected params to contain 'param2'")
	}

	if !params.Contains("param3") {
		t.Errorf("Expected params to contain 'param3'")
	}

	if params.Contains("param4") {
		t.Errorf("Expected params not to contain 'param4'")
	}
}

func TestNewParam(t *testing.T) {
	// 创建一个新参数
	paramName := "testParam"
	param := NewParam(paramName)

	// 检查参数的名称是否正确
	if param.Name() != paramName {
		t.Errorf("Expected param name to be '%s', got '%s'", paramName, param.Name())
	}
}
