package routegroup

import (
	"itflow/handle"
	"itflow/midware"
	"itflow/model"

	"github.com/hyahm/xmux"
)

var UserGroup *xmux.GroupRoute

func init() {
	UserGroup = xmux.NewGroupRoute().AddMidware(midware.CheckUserGroupPermssion).ApiCreateGroup("usergroup", "用户组相关", "user group")
	UserGroup.ApiCodeField("code").ApiCodeMsg("0", "成功")
	UserGroup.ApiCodeField("code").ApiCodeMsg("20", "token过期")
	UserGroup.ApiCodeField("code").ApiCodeMsg("", "其他错误,请查看返回的msg")
	UserGroup.ApiReqHeader("X-Token", "xxxxxxxxxxxxxxxxxxxxxxxxxx")

	UserGroup.Pattern("/group/get").Post(handle.GroupGet).ApiResponseTemplate(`"grouplist":[],"code":0`)

	UserGroup.Pattern("/group/add").Post(handle.GroupAdd).Bind(&model.Get_groups{}).AddMidware(midware.JsonToStruct).
		End(midware.EndLog)

	UserGroup.Pattern("/group/del").Get(handle.GroupDel).End(midware.EndLog)

	UserGroup.Pattern("/group/update").Post(handle.GroupUpdate).Bind(&model.Get_groups{}).AddMidware(midware.JsonToStruct).
		End(midware.EndLog)
}