package mongodb

import (
	"gopkg.in/mgo.v2"
)

var (
	//c_admin_role_assoc_popedoms *mgo.Collection
	collection map[string]*mgo.Collection = make(map[string]*mgo.Collection)
)

// do this when reconnect to the database
func deinintCollections(table string) {
	//c_admin_role_assoc_popedoms = database.C(tb_admin_role_assoc_popedoms)
	collection[table] = database.C(table)
	initCollections(table, collection[table])
}

func initCollections(table string, ct *mgo.Collection) {
	initCollection(table, &ct, "")
}

func initCollection(table string, collection **mgo.Collection, indexKey ...string) {
	*collection = database.C(table)
	if len(indexKey) != 0 && indexKey[0] != "" {
		_ = (*collection).EnsureIndexKey(indexKey...)
	}
}
