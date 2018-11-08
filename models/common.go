/************************

	公共操作

*************************/

package models

// 返回表中的数据条数 用于数据统计
func TableNumber(table string) (count int64, err error) {
	count, err = dbc.QueryTable(table).Count()
	return
}