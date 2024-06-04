package models

type User struct {
	UserName string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Query struct {
	// d     []string `vd:"@:len($)>0 && $[0]=='D'; msg:sprintf('invalid d: %v',$)"`
	Name  string `uri:"name" json:"name" vd:"len($)>0 && $!='admin';msg:sprintf('无效参数 name:%v',$)"`
	Phone string `uri:"phone" json:"phone" vd:"phone($)"`
	Email string `uri:"email" json:"email" vd:"email($)"`
}
