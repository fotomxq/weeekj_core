package CoreFilter

import "testing"

func TestCheckNiceName(t *testing.T) {
	var b bool
	b = CheckNiceName("abc很久很久（）")
	if !b {
		t.Error("check nice name 1")
	}
}

func TestCheckEmail(t *testing.T) {
	var b bool
	b = CheckEmail("booking@blueonline.com.au")
	if !b {
		t.Error("check email failed booking@blueonline.com.au")
	}
	b = CheckEmail("booking@blueonline.com")
	if !b {
		t.Error("check email failed booking@blueonline.com")
	}
	b = CheckEmail("booking@qq.com")
	if !b {
		t.Error("check email failed booking@qq.com")
	}
	b = CheckEmail("booking_ok@163.com")
	if !b {
		t.Error("check email failed booking_ok@163.com")
	}
	b = CheckEmail("gaozemin0509@gmail.com")
	if !b {
		t.Error("check email failed gaozemin0509@gmail.com")
	}
}

func TestCheckMark(t *testing.T) {
	var b bool
	b = CheckMark("all")
	if !b {
		t.Error("check nice name 1")
	}
}

func TestCheckMarkPage(t *testing.T) {
	var b bool
	b = CheckMarkPage("/index")
	if !b {
		t.Error("check nice name 1")
	}
}
