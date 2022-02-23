package dto

import "gin/model"

type UserDto struct{
	//UserDto... 我们希望只返回哟过户的名称和手机号，其他敏感信息不用返回
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.User) UserDto{
	return UserDto{
		Name:	user.Name,
		Telephone: user.Telephone,
	}
}