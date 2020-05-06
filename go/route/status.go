package route

import (
	"itflow/app/handle"
	"itflow/midware"
	"itflow/network/bug"

	"github.com/hyahm/xmux"
)

var Status *xmux.GroupRoute

func init() {
	Status = xmux.NewGroupRoute().AddMidware(midware.CheckStatusPermssion).
		ApiCreateGroup("bugstatus", "bug 状态管理", "bug status").
		ApiReqHeader("X-Token", "asdfasdfasdfasdfsdf").
		ApiCodeMsg("0", "成功").ApiCodeMsg("3", "没有权限").ApiCodeMsg("10", "token 过期")

	Status.Pattern("/status/list").Post(handle.StatusList).
		ApiResStruct(&bug.ListStatus{}).ApiDescribe("获取bugstatus状态列表").
		ApiResponseTemplate(`{
			"statuslist": [
				{
					"id": 3,
					"name": "ToDoList"
				},
				{
					"id": 6,
					"name": "测试"
				},
				{
					"id": 7,
					"name": "need13"
				},
				{
					"id": 5,
					"name": "react"
				}
			],
			"code": 0
		}`)

	Status.Pattern("/status/add").Post(handle.StatusAdd).Bind(&bug.Status{}).
		End(midware.EndLog).AddMidware(midware.JsonToStruct).
		ApiDescribe("添加bug 状态").
		ApiReqStruct(&bug.Status{}).ApiRequestTemplate(`{"id": 0, "name": "普通"}`).
		ApiResStruct(&bug.ResponeStatus{}).ApiResponseTemplate(`{"id": 8,"code": 0}`)

	Status.Pattern("/status/remove").Get(handle.StatusRemove).
		End(midware.EndLog).
		ApiDescribe("删除bug 状态").ApiSupplement("当此状态有bug在使用时， 无法删除")

	Status.Pattern("/status/update").Post(handle.StatusUpdate).Bind(&bug.Status{}).AddMidware(midware.JsonToStruct).
		End(midware.EndLog).
		ApiDescribe("修改状态").
		ApiReqStruct(&bug.Status{}).ApiRequestTemplate(`{"id": 0, "name": "普通"}`).
		ApiResStruct(&bug.ResponeStatus{}).ApiResponseTemplate(`{"code": 0}`)
}
