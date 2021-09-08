package mongodb

import (
	"time"

	"github.com/weijun-sh/yisheng-mysql/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Add_admin_role_assoc_popedoms(ms *Struct_admin_role_assoc_popedoms, overwrite bool) (err error) {
	if overwrite {
		_, err = c_admin_role_assoc_popedoms.UpsertId(ms.Id, ms)
		return err
	} else {
		err = c_admin_role_assoc_popedoms.Insert(ms)
	}
	if err == nil {
		log.Info("[mongodb] Add admin_role_assoc_popedoms success", "role", ms)
	} else {
		log.Warn("[mongodb] Add admin_role_assoc_popedoms failed", "role", ms, "err", err)
	}
	return err
}

