package toolLdap

import (
	"github.com/jtblin/go-ldap-client"
	"log"
	"testing"
)

func Test_ldap(t *testing.T) {
	//username := "feixiaohui"
	//password := "AIOUniya520!"
	//
	//// 用来获取查询权限的 bind 用户.如果 ldap 禁止了匿名查询,那我们就需要先用这个帐户 bind 以下才能开始查询
	//// bind 的账号通常要使用完整的 DN 信息.例如 cn=manager,dc=example,dc=org
	//// 在 AD 上,则可以用诸如 mananger@example.org 的方式来 bind
	//bindusername := "LDAP"
	//bindpassword := "Hw8C5rjRBP!"
	//
	////l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "192.168.12.6", 636))
	////if err != nil {
	////	log.Fatal(err)
	////}
	////defer l.Close()
	//
	//
	//tlsConfig := &tls.Config{InsecureSkipVerify: true}
	//l, err := ldap.DialTLS("tcp", "192.168.12.6:636", tlsConfig)
	//if err != nil{
	//	panic(err)
	//}
	//fmt.Println(l)
	//
	//// No TLS, not recommended
	////l, err := ldap.Dial("tcp", "ldap.example.com:389")
	//
	//fmt.Println(username, password)
	//fmt.Println(bindusername, bindpassword)

	client := &ldap.LDAPClient{
		Base:         "dc=jms,dc=com",
		Host:         "192.168.12.6",
		Port:         636,
		UseSSL:       false,
		BindDN:       "uid=readonlysuer,ou=04YLAdminUser,dc=jms,dc=com",
		BindPassword: "readonlypassword",
		UserFilter:   "(uid=%s)",
		GroupFilter:  "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
	// It is the responsibility of the caller to close the connection
	defer client.Close()

	ok, user, err := client.Authenticate("LDAP", "Hw8C5rjRBP")
	if err != nil {
		log.Fatalf("Error authenticating user %s: %+v", "username", err)
	}
	if !ok {
		log.Fatalf("Authenticating failed for user %s", "username")
	}
	log.Printf("User: %+v", user)
}
