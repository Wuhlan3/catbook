package controller

import (
	"gin/common"
	"gin/model"
	"gin/response"
	"gin/vo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type IPostController interface {
	PageList(ctx *gin.Context)
	RestController
}

type PostController struct {
	DB *gorm.DB
}

func NewPostController() PostController {
	db := common.GetDB()
	db.AutoMigrate(model.Cat{})
	db.AutoMigrate(model.Post{})
	db.AutoMigrate(model.Category{})
	db.AutoMigrate(model.Appear{})
	db.AutoMigrate(model.Location{})
	return PostController{DB: db}
}


func (p PostController) PageList(ctx *gin.Context) {
	return

}

func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	//插入location表
	var myLocation model.Location
	if err := p.DB.Where("place = ?", requestPost.Place).First(&myLocation).Error; err != nil{
		myLocation.Place = requestPost.Place
		if err := p.DB.Create(&myLocation).Error; err != nil {
			panic(err)
			return
		}
		p.DB.Where("place = ?", requestPost.Place).First(&myLocation)
	}
	//插入category表
	var myCategory model.Category
	if err := p.DB.Where("breed = ?", requestPost.Breed).First(&myCategory).Error; err != nil{
		myCategory.Breed = requestPost.Breed
		if err := p.DB.Create(&myCategory).Error; err != nil {
			panic(err)
			return
		}
		p.DB.Where("breed = ?", requestPost.Breed).First(&myCategory)
	}

	//插入cat表
	cat := model.Cat{
		Name:  		requestPost.Name,
		Color: 		requestPost.Color,
		CategoryId:		myCategory.ID,
	}

	var myCat model.Cat	//用于接收sql查询
	//如果不存在这只猫，就创建一个新的行
	if err := p.DB.Where("name = ?", requestPost.Name).First(&myCat).Error; err != nil{
		if err := p.DB.Create(&cat).Error; err != nil {
			panic(err)
			return
		}
		p.DB.Where("name = ?", requestPost.Name).First(&myCat)
	}

	//插入appear表
	appear := model.Appear{
		Cid:  		myCat.ID,
		LocationId: 	myLocation.ID,
		Specific:		requestPost.Specific,
	}
	var myAppear model.Appear
	if err := p.DB.Where("location_id = ?", myAppear.LocationId).First(&myAppear).Error; err != nil{
		if err := p.DB.Create(&appear).Error; err != nil {
			panic(err)
			return
		}
		p.DB.Where("location_id = ?", myAppear.LocationId).First(&myAppear)
	}
	//插入post表
	//获取登录用户
	user, _ := ctx.Get("user")
	//保存一个新的图片
	myPost := model.Post{
		Uid:		user.(model.User).ID,
		Cid:		myCat.ID,
		ImgUrl:		requestPost.ImgUrl,
	}

	if err := p.DB.Create(&myPost).Error; err != nil {
		panic(err)
		return
	}
	response.Success(ctx, nil, "上传成功")
}

func (p PostController) Update(ctx *gin.Context) {
	var updaterequest vo.CreateUpdateRequest
	if err := ctx.ShouldBind(&updaterequest); err != nil {
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	//得到猫的旧名字updaterequest.Name,验证是否存在
	var cat model.Cat
	p.DB.Where("Name = ?", updaterequest.Name).Find(&cat)
	if cat.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "要更新的猫不存在")
		return
	}
	var newcat model.Cat
	p.DB.Where("Name = ?", updaterequest.NewName).Find(&newcat)
	if newcat.ID != 0 && updaterequest.NewName != updaterequest.Name {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "请换个名字")
		return
	}
	var post model.Post
	p.DB.Where("Cid = ?", cat.ID).Find(&post)
	if post.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "该记录不存在")
		return
	}
	user, _ := ctx.Get("user")
	if user.(model.User).ID != post.Uid {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "这只猫不是您上传的哦")
		return
	}
	//找
	var category model.Category
	p.DB.Where("ID = ?", cat.CategoryId).Find(&category)
	if category.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "该记录不存在")
		return
	}
	var appear model.Appear
	p.DB.Where("Cid = ?", cat.ID).Find(&appear)
	if appear.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "该记录不存在")
		return
	}
	var location model.Location
	p.DB.Where("ID = ?", appear.LocationId).Find(&location)
	if location.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "该记录不存在")
		return
	}

	//改新名字
	if updaterequest.Name != updaterequest.NewName {
		p.DB.Model(&cat).Where("Name = ?", updaterequest.Name).Update("Name", updaterequest.NewName)
	}
	//改颜色
	p.DB.Model(&cat).Where("ID = ?", cat.ID).Update("Color", updaterequest.NewColor)
	//改breed
	p.DB.Model(&category).Where("ID = ?", cat.CategoryId).Update("Breed", updaterequest.NewBreed)
	//改Specific
	p.DB.Model(&appear).Where("Cid = ?", cat.ID).Update("Specific", updaterequest.NewSpecific)
	//改Place
	p.DB.Model(&location).Where("ID = ?", appear.LocationId).Update("Place", updaterequest.NewPlace)
	response.Success(ctx, nil, "更新成功")

	return
}

func (p PostController) Delete(ctx *gin.Context) {
	var deleterequest vo.DeleteRequest
	if err := ctx.ShouldBind(&deleterequest); err != nil {
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	//得到猫的名字deleterequest.CatName
	//验证CatName是否存在
	var cat model.Cat
	p.DB.Where("Name = ?", deleterequest.CatName).Find(&cat)
	if cat.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "要删除的猫名字不存在")
		return
	}
	//根据CAT表的ID删掉post中那条记录
	var post model.Post
	p.DB.Where("Cid = ?", cat.ID).Find(&post)
	if post.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "post表记录中不存在该猫名")
		return
	}
	user, _ := ctx.Get("user")
	if user.(model.User).ID != post.Uid {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "这只猫不是您上传的哦")
		return
	}
	p.DB.Unscoped().Delete(&post)
	//根据cat表的ID删掉Appear表中相同cai_id记录，再用这个Appear记录的location_id删掉location表那条记录
	var appear model.Appear
	p.DB.Where("Cid = ?", cat.ID).Find(&appear)
	if appear.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "appear表记录中不存在该猫名")
		return
	}
	p.DB.Unscoped().Delete(&appear)
	//删除cat表的记录
	p.DB.Unscoped().Delete(&cat)
	//返回成功
	response.Success(ctx, nil, "删除成功")
	return
}



type NeedToShow struct {
	UserName string `json:"user_name"`
	CatName  string `json:"cat_name"`
	Color    string `json:"color"`
	Breed    string `json:"breed"`
	Place    string `json:"place"`
	Specific string `json:"specific"`
	ImgUrl   string `json:"img_url"`
}

func (p PostController) Show(ctx *gin.Context) {
	var cats []model.Cat
	var NeedToShows []NeedToShow
	p.DB.Find(&cats)
	n := len(cats)
	//遍历cats表中每一只猫记录
	for _, cat := range cats {
		var toShow NeedToShow
		toShow.CatName = cat.Name
		toShow.Color = cat.Color
		//用cat中的CategoryId去Category表中找到toShow要用的的Breed
		var category model.Category
		p.DB.Where("ID = ?", cat.CategoryId).Find(&category)
		toShow.Breed = category.Breed
		//去Appear表找locationid，再去location找place，Appear表同时也有Specific
		var appear model.Appear
		p.DB.Where("Cid = ?", cat.ID).Find(&appear)
		toShow.Specific = appear.Specific
		var location model.Location
		p.DB.Where("ID = ?", appear.LocationId).Find(&location)
		toShow.Place = location.Place
		//找到ImgUrl
		var post model.Post
		p.DB.Where("Cid = ?", cat.ID).Find(&post)
		toShow.ImgUrl = post.ImgUrl
		var user model.User
		p.DB.Where("ID = ?", post.Uid).Find(&user)
		toShow.UserName = user.Name
		NeedToShows = append(NeedToShows, toShow)
	}
	response.Success(ctx, gin.H{
		"data":        NeedToShows,
		"TotalLength": n,
	}, "")
	return
}