package controller

import (
	"gin/common"
	"gin/dto"
	"gin/model"
	"gin/response"
	"gin/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type testUser struct {
	Name    string   `json:"name"`
	Telephone  string `json:"telephone"`
	Password string `json:"password"`
}

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	// 使用map获取请求参数
	//var requestMap = make(map[string]string)
	//json.NewDecoder(ctx.Request.Body).Decode(&requestMap)

	// 使用结构体接收
	//var requestUser = model.testUser{}
	var requestUser testUser
	//json.NewDecoder(ctx.Request.Body).Decode(&requestUser)
	ctx.Bind(&requestUser) //支持所有的类型的解析，将ctx接受到的ctx.Request.Body，解析到user结构体中
	//ctx.BindJSON(&requestUser)
	//获取参数
	//name := ctx.PostForm("name")
	//telephone := ctx.PostForm("telephone")
	//password := ctx.PostForm("password")
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	//fmt.Println(name)
	//fmt.Println(telephone)
	//fmt.Println(password)
	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil, "手机号须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil, "密码至少为6位")
		return
	}
	//随机传一个姓名
	if len(name) == 0 {
		name = util.RandomString(8)
	}

	log.Println(name,telephone,password)
	//判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil, "用户已经存在了")
		return
	}
	//创建用户
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		response.Response(ctx, http.StatusInternalServerError, 500,nil, "密码加密错误")
		return
	}
	newUser := model.User{
		Name: name,
		Telephone: telephone,
		Password: string(hashPassword),
	}
	DB.Create(&newUser)

	//发放token给前端
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 500,nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	//返回结果
	//response.Success(ctx,nil,"注册成功")
	response.Success(ctx, gin.H{"token": token},"注册成功")

}

func Login(ctx *gin.Context){
	DB := common.GetDB()
	//获取参数
	var requestUser testUser
	ctx.BindJSON(&requestUser)
	telephone := requestUser.Telephone
	password := requestUser.Password
	//telephone := ctx.PostForm("telephone")
	//password := ctx.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil, "手机号须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil, "密码至少为6位")
		return
	}
	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422,nil, "该用户不存在了")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) ; err != nil {
		response.Fail(ctx,nil,"密码错误")
		return
	}
	//发放token给前端
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 500,nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	//返回结果
	response.Success(ctx, gin.H{"token": token},"登录成功")
}

func Info(ctx *gin.Context){
	user,_ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

func isTelephoneExist(db * gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}