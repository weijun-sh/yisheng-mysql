package mongodb

import (
	"gopkg.in/mgo.v2"
)

var (
	c_admin_role_assoc_popedoms *mgo.Collection
)

// do this when reconnect to the database
func deinintCollections() {
	c_admin_role_assoc_popedoms = database.C(tb_admin_role_assoc_popedoms)
}

func initCollections() {
	initCollection(tb_admin_role_assoc_popedoms, &c_admin_role_assoc_popedoms, "")
}

func initCollection(table string, collection **mgo.Collection, indexKey ...string) {
	*collection = database.C(table)
	if len(indexKey) != 0 && indexKey[0] != "" {
		_ = (*collection).EnsureIndexKey(indexKey...)
	}
}
