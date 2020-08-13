package checks

import (
	"testing"
)

func TestCheckCode(t *testing.T) {
	httpCheck := &URLCheck{
		RightCode: 200,
	}

	if httpCheck.checkCode(201) != false {
		t.Error()
	}

	if httpCheck.checkCode(200) != true {
		t.Error()
	}

	if httpCheck.checkCode(300) != false {
		t.Error()
	}
}
