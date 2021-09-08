package mongodb

import (
//"gopkg.in/mgo.v2/bson"
)

const (
	tb_admin_role_assoc_popedoms string = "admin_role_assoc_popedoms"
)

//admin_role_assoc_popedoms, [id roleId popedomCode
type Struct_admin_role_assoc_popedoms struct {
	Id          string `bson:"id"`
	RoleId      string `bson:"roleId"`
	PopedomCode string `bson:"popedomCode"`
}
