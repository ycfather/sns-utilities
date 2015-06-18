package ldap

import (
	"fmt"

	"github.com/mmitton/ldap"
)

const (
	DEFAULT_DOMAIN = "iqiyi"
)

func Auth(host string, port int, username, password string) bool {
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		// panic(err)
		return false
	}
	defer conn.Close()

	if err = conn.Bind(fmt.Sprintf("%s\\%s", DEFAULT_DOMAIN, username), password); err != nil {
		// panic(err)
		return false
	}

	return true
}
