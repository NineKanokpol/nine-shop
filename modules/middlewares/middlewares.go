package middlewares

//ขั้นตอนที่ 2 ในการทำ role base access control
type Role struct {
	Id    int    `db:"id"`
	Title string `db:"title"`
}
