package ldap

import (
	"testing"
)

func TestAuth(t *testing.T) {
	ret := Auth("10.11.50.65", 389, "aaronwang", "qiyi@13O819")
	if !ret {
		t.Fatalf("auth result : %v, want : true", ret)
	}
}
