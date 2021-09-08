package mongodb

import (
	"fmt"
	//"github.com/weijun-sh/yisheng-mysql/log"
)

//func Add_admin_role_assoc_popedoms(ms *Struct_admin_role_assoc_popedoms, overwrite bool) (err error) {
//	if overwrite {
//		_, err = c_admin_role_assoc_popedoms.UpsertId(ms.Id, ms)
//		return err
//	} else {
//		err = c_admin_role_assoc_popedoms.Insert(ms)
//	}
//	if err == nil {
//		log.Info("[mongodb] Add admin_role_assoc_popedoms success", "role", ms)
//	} else {
//		log.Warn("[mongodb] Add admin_role_assoc_popedoms failed", "role", ms, "err", err)
//	}
//	return err
//}

func Insert(table string, docs []interface{}) (err error) {
	deinintCollections(table)
	c := collection[table]
	err = c.Insert(docs...)
	if err != nil {
		fmt.Printf("Insert, docs: %v, err: %v\n", docs, err)
	}
	return nil
}

