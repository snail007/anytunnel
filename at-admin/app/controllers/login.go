package controllers

import (
	"anytunnel/at-admin/app/models"
	"anytunnel/at-admin/app/utils"
	"image/color"
	"image/png"
	"strings"

	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
)

var (
	cap = captcha.New()
)

func init() {
	bs, _ := utils.NewEncrypt().Base64DecodeBytes(fontData)
	cap.AddFontFromBytes(bs)
}

type LoginController struct {
	BaseController
}

func (this *LoginController) Index() {
	this.Layout = "layout/login.html"
	this.TplName = "login/login.html"
	this.Data["title"] = beego.AppConfig.String("sys.name") + "登录"
	this.Render()
}

//login
func (this *LoginController) Login() {
	if this.isLogin() {
		return
	}
	userModel := models.User{}
	name := strings.TrimSpace(this.GetString("username"))
	password := strings.TrimSpace(this.GetString("password"))
	captcha := strings.TrimSpace(this.GetString("captcha"))
	captchaSession := this.GetSession("captcha")
	if captchaSession == nil || captcha == "" || captchaSession != strings.ToLower(captcha) {
		this.SetSession("captcha", "")
		this.jsonError("验证码错误!")
	}
	this.SetSession("captcha", "")
	user, err := userModel.GetUserByName(name)
	if err != nil {
		this.jsonError(err)
		return
	}
	if len(user) == 0 {
		this.jsonError("账号错误!")
	}
	encrypt := new(utils.Encrypt)
	password = userModel.EncodePassword(password)

	if user["password"] != password {
		this.jsonError("账号或密码错误!")
	}
	//加载权限列表

	//保存 session
	this.SetSession("author", user)
	//保存 cookie
	identify := encrypt.Md5Encode(this.Ctx.Request.UserAgent() + this.getClientIp() + password)
	passportValue := encrypt.Base64Encode(name + "@" + identify)
	passport := beego.AppConfig.String("author.passport")
	//fmt.Println("set cookie " + passportValue)
	this.Ctx.SetCookie(passport, passportValue, 3600)

	this.jsonSuccess("登录成功", "", "/main/index.html")
}

//logout
func (this *LoginController) Logout() {
	passport := beego.AppConfig.String("author.passport")
	this.Ctx.SetCookie(passport, "")
	this.SetSession("author", "")
	this.Redirect("/login/index.html", 302)
	this.StopRun()
}

func (this *LoginController) Captcha() {
	//bs, _ := ioutil.ReadFile("static/index/comic.ttf")
	//ioutil.WriteFile("a.txt", utils.NewEncrypt().Base64EncodeBytes(bs), os.ModeAppend)
	cap.SetSize(80, 28)
	cap.SetDisturbance(captcha.NORMAL)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 100, 100, 100})
	img, str := cap.Create(4, captcha.ALL)
	this.SetSession("captcha", strings.ToLower(str))
	png.Encode(this.Ctx.ResponseWriter, img)
}
