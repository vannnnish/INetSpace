package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"netspace/db"
	"netspace/util"
	"time"
)

const (
	pwd_salt = "#g325g"
)

func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

func DoSignupHandler(c *gin.Context) {

	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")
	if len(username) < 3 || len(passwd) < 5 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号密码格式不正确",
		})
		return
	}
	encPasswd := util.Sha1([]byte(passwd + pwd_salt))
	suc := db.UserSignUp(username, encPasswd)
	if suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Success",
			"code": 0,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "Fail",
			"code": -1,
		})
		return
	}
}

// SignInHandler:登录接口
func SignInHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

func DoSignInHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	encPassword := util.Sha1([]byte(password + pwd_salt))
	fmt.Println(encPassword)
	// 1. 校验用户名称及密
	pwdChecked := db.UserSignIn(username, encPassword)
	if !pwdChecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "密码错误 Failed",
			"code": -1,
		})
		return
	}
	// 2. 生成访问凭证
	token := GetToken(username)
	updateToken := db.UpdateToken(username, token)
	if !updateToken {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "密码错误 Failed",
			"code": -1,
		})
		return
	}
	// 3. 返回登录成功信息
	//w.Write([]byte("htt://" + r.Host + "/static/view/home.html"))
	msg := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}

	c.Data(http.StatusOK, "application/json", msg.JSONBytes())
}

func UserInfoHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	/*	token := r.Form.Get("token")
		// 2. 验证token是否有效
		isValidToken := IsTokenValid(token)
		if !isValidToken {
			w.WriteHeader(http.StatusForbidden)
			return
		}*/
	// 3. 查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.JSON(http.StatusOK, resp.JSONBytes())
}

func GetToken(username string) string {
	// md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

func IsTokenValid(token string) bool {
	//TODO: 判断token时效性
	//TODO: 从数据库表tbl_user_token查询username对应token信息
	//TODO: 对比两个token是否一致
	return true
}
