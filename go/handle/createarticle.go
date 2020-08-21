package handle

import (
	"encoding/json"
	"html"
	"itflow/cache"
	"itflow/db"
	"itflow/internal/bug"
	"itflow/internal/project"
	"net/http"
	"strings"
	"time"

	"github.com/hyahm/golog"
	"github.com/hyahm/gomysql"
	"github.com/hyahm/xmux"
	//"strings"
)

type projectList struct {
	ProjectList []string `json:"projectlist"`
	Code        int      `json:"code"`
}

func GetProject(w http.ResponseWriter, r *http.Request) {

	pl := &projectList{}

	for _, v := range cache.CachePidProject {
		pl.ProjectList = append(pl.ProjectList, v.ToString())
	}
	send, _ := json.Marshal(pl)
	w.Write(send)
	return

}

func GetMyProject(w http.ResponseWriter, r *http.Request) {
	myproject := &project.MyProject{
		Name: make([]string, 0),
	}
	uid := xmux.GetData(r).Get("uid").(int64)

	w.Write(myproject.Get(uid))
	return

}

// 添加或编辑
func BugCreate(w http.ResponseWriter, r *http.Request) {
	uid := xmux.GetData(r).Get("uid").(int64)
	data := xmux.GetData(r).Data.(*bug.EditBug)
	// createdId := cache.DefaultCreateSid
	// if createdId == 0 {
	// 	golog.Error("必须给定一个状态默认值")
	// 	w.Write(errorcode.Error("必须给定一个状态默认值")
	// 	return
	// }
	var ids = make([]string, 0)
	rows, err := db.Mconn.GetRowsIn("select id from user where realname in (?)",
		(gomysql.InArgs)(data.Selectusers).ToInArgs())
	if err != nil {
		golog.Error(err)
		w.Write(data.ErrorE(err))
		return
	}
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			golog.Error(err)
			continue
		}
		ids = append(ids, id)
	}

	// bug, err := data.ToBug()
	// if err != nil {
	// 	w.Write(errorcode.ErrorE(err))
	// 	return
	// }
	// // bug.StatusId = createdId
	// bug.Uid = xmux.GetData(r).Get("uid").(int64)
	//
	// go datalog.InsertLog("bug", nickname+"create bug: "+data.Title, r.RemoteAddr, nickname, "create")

	if data.Id <= 0 {
		// 插入bug
		// err = bug.CreateBug()
		// if err != nil {
		// 	w.Write(errorcode.ErrorE(err))
		// 	return
		// }
		if cache.DefaultCreateSid == 0 {
			w.Write(data.Error("必须先设置一个默认的创建状态（系统设置->默认值->bug创建时的状态）"))
			return
		}
		createsql := `insert into bugs(uid,title,sid,content,iid,createtime,vid,spusers,lid,eid,pid) 
		values(?,?,?,?,(select ifnull(min(id),0) from importants where name=?),?,
		(select ifnull(min(id),0) from version where name=?),?,
		(select ifnull(min(id),0) from level where name=?),
		(select ifnull(min(id),0) from environment where name=?),
		(select ifnull(min(id),0) from project where name=?))`

		data.Id, err = db.Mconn.Insert(createsql, uid, data.Title,
			cache.DefaultCreateSid, html.EscapeString(data.Content), data.Important,
			time.Now().Unix(), data.Version, strings.Join(ids, ","), data.Level, data.Envname, data.Projectname)
		if err != nil {
			golog.Error(err)
			w.Write(data.ErrorE(err))
			return
		}

	} else {
		// update
		updatesql := `update bugs set title=?,content=?,
			iid=(select ifnull(min(id),0) from importants where name=?),
			updatetime=?,
			vid=(select ifnull(min(id),0) from version where name=?),
			spusers=?,
			lid=(select ifnull(min(id),0) from level where name=?),
			eid=(select ifnull(min(id),0) from environment where name=?),
			pid=(select ifnull(min(id),0) from project where name=?) 
		where uid=? and id=?`

		_, err = db.Mconn.Update(updatesql, data.Title,
			html.EscapeString(data.Content), data.Important,
			time.Now().Unix(), data.Version, strings.Join(ids, ","),
			data.Level, data.Envname, data.Projectname, uid, data.Id)
		if err != nil {
			golog.Error(err)
			w.Write(data.ErrorE(err))
			return
		}

		// err = bug.EditBug()
		// if err != nil {
		// 	w.Write(errorcode.ErrorE(err))
		// 	return
		// }
		// go datalog.InsertLog("bug", nickname+fmt.Sprintf(" update bug id: %d", data.Id), r.RemoteAddr, nickname, "update")

	}

	w.Write(data.Marshal())
	return

}
