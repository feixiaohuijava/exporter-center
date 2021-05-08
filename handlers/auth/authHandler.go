package auth

import (
	"crypto/tls"
	errors "errors"
	"exporter-center/config"
	"exporter-center/config/configStruct"
	"exporter-center/logs"
	_ "exporter-center/logs"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"github.com/goinggo/mapstructure"
	"net/http"
	"time"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 1

var MySecret = []byte("夏天夏天悄悄过去")

// GetToken 生成JWT
func GetToken(username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "my-project",                               // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// @Tags 获取token接口
// @Summary 获取token接口
// @Description 获取token接口
// @Accept application/json
// @Produce application/json
// @Param username body string false "AD账号"
// @Param password body string false "AD密码"
// @Router /auth/login [post]
func AuthHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user UserInfo
	var err error
	err = c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的参数",
		})
		return
	}

	//获取配置文件
	var adConfig configStruct.AdConfig
	adStruct := config.GetYamlConfig("config_ad", &adConfig)
	err = mapstructure.Decode(adStruct, &adConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "解析配置文件出错",
		})
		return
	}
	logs.Infoln("获取到的配置地址:", adConfig)

	err = LdapCertifi(user, adConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "校验出错!",
		})
		return
	}
	// 生成Token
	tokenString, _ := GetToken(user.Username)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": gin.H{"token": tokenString},
	})
	return
}

func LdapCertifi(user UserInfo, adConfig configStruct.AdConfig) error {

	// 此处就不进行捕获了
	//defer func() {
	//	if e := recover(); e != nil {
	//		logger.Errorln("ldap处理报错,原因:", e)
	//	}
	//}()

	var err error
	var conn *ldap.Conn
	var sr *ldap.SearchResult
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	conn, err = ldap.DialTLS("tcp", adConfig.Ad.Address, tlsConfig)
	defer conn.Close()
	if err != nil {
		logs.Errorln(err)
		return err
	}
	err = conn.Bind(adConfig.Ad.BindUsername, adConfig.Ad.BindPassword)
	if err != nil {
		logs.Errorln(err)
		return err
	}

	searchRequest := ldap.NewSearchRequest(
		// 这里是 basedn,我们将从这个节点开始搜索
		adConfig.Ad.Basedn,
		// 这里几个参数分别是 scope, derefAliases, sizeLimit, timeLimit,  typesOnly
		// 详情可以参考 RFC4511 中的定义,文末有链接
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		// 这里是 LDAP 查询的 Filter.这个例子例子,我们通过查询 uid=username 且 objectClass=organizationalPerson.
		// username 即我们需要认证的用户名
		fmt.Sprintf("(&(objectClass=organizationalPerson)(sAMAccountName=%s))", user.Username),
		// 这里是查询返回的属性,以数组形式提供.如果为空则会返回所有的属性
		//[]string{"dn"},
		[]string{},
		nil,
	)
	sr, err = conn.Search(searchRequest)
	if err != nil {
		logs.Errorln(err)
		return err
	}

	logs.Infoln("获取的返回值", sr.Entries)

	if len(sr.Entries) != 1 {
		logs.Errorln("用户不存在或者值重复!")
		return errors.New("用户不存在或者值重复!")
	}
	userdn := sr.Entries[0].DN
	logs.Infoln(sr.Entries[0])
	logs.Infoln(sr.Entries[0].Attributes[0].Name)
	logs.Infoln(userdn)

	err = conn.Bind(userdn, user.Password)
	if err != nil {
		logs.Errorln("用户名密码不匹配:", err)
		return err
	}
	return nil
}
