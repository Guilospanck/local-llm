package domain

var Queries = map[string]map[string]string{
	"property": {
		"*":          "SELECT * FROM property",
		"*_id":       "SELECT * FROM property WHERE id in (%s)",
		"*_color":    "SELECT * FROM property WHERE color = '%s'",
		"*_price":    "SELECT * FROM property WHERE price >= %f AND price <= %f",
		"*_size_sqm": "SELECT * FROM property WHERE size_sqm >= %f AND size_sqm <= %f",
	},
	"view": {
		"id_view": `SELECT v.id from "view" v WHERE v.view ILIKE %s`,
	},
	"property_views": {
		"id_viewid": `SELECT DISTINCT property_id from property_views WHERE view_id in (%s)`,
	},
}
