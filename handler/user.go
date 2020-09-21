package handler

import (
	"fmt"
	"net/http"
	"netspace/db"
	"netspace/util"
	"time"
)

const (
	pwd_salt = "#g325g"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Write([]byte("注册内容"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println("表单解析失败:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(r.Form)
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("账号密码格式不正确"))
		return
	}
	encPasswd := util.Sha1([]byte(passwd + pwd_salt))
	suc := db.UserSignUp(username, encPasswd)
	if suc {
		w.Write([]byte("Success"))
	} else {
		w.Write([]byte("Failed"))
	}
}

// SignInHandler:登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPassword := util.Sha1([]byte(password + pwd_salt))
	fmt.Println(encPassword)
	// 1. 校验用户名称及密
	pwdChecked := db.UserSignIn(username, encPassword)
	if !pwdChecked {
		w.Write([]byte("密码错误 Failed"))
		return
	}
	// 2. 生成访问凭证
	token := GetToken(username)
	updateToken := db.UpdateToken(username, token)
	if !updateToken {
		w.Write([]byte("访问凭证 Faild"))
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
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(msg.JSONBytes())
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
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
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
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
