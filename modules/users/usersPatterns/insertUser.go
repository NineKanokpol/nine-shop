package usersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NineKanokpol/Nine-shop-test/modules/users"
	"github.com/jmoiron/sqlx"
)

// ต้อง insert ได้ทั้ง admin ใช้ design pattern แบบ factory

// โรงงานใหญ่
type IInsertUser interface {
	Customer() (IInsertUser, error)
	Admin() (IInsertUser, error)
	Result() (*users.UserPassport, error)
}

// struct ตัวแม่
type userReq struct {
	id  string
	req *users.UserRegisterRequest
	db  *sqlx.DB
}

// struct ตัวลูก
type customer struct {
	//การไม่ใช้ตัวแปรข้างหน้าจะ dot เข้าไปใน struct ได้เลยโดยไม่ต้อง dot อะไรก่อน
	*userReq
}

// struct ตัวลูก
type admin struct {
	*userReq
}

// โรงงานใหญ่ของตัว newAdmin กับ newCustomer
func InsertUser(db *sqlx.DB, req *users.UserRegisterRequest, isAdmin bool) IInsertUser {
	if isAdmin {
		return newAdmin(db, req)
	}
	return newCustomer(db, req)
}

// โรงงานย่อย
func newCustomer(db *sqlx.DB, req *users.UserRegisterRequest) IInsertUser {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

// โรงงานย่อย
func newAdmin(db *sqlx.DB, req *users.UserRegisterRequest) IInsertUser {
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

// การ insert เข้า db
func (f *userReq) Customer() (IInsertUser, error) {
	//insert ไม่เสร็จภายใน 5 วิ คือ error ละ
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	//คำสั่ง db
	// ($1,$2,$3) ตาม parameter
	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2, $3, 2)
	RETURNING "id";`

	// มี returning ต้อง ใช้ Scan
	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}

	return f, nil
}

// insert admin ขั้นที่ 1
func (f *userReq) Admin() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2, $3, 1)
	RETURNING "id";`

	if err := f.db.QueryRowContext(
		ctx,
		query,
		f.req.Email,
		f.req.Password,
		f.req.Username,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}

	return f, nil
}

func (f *userReq) Result() (*users.UserPassport, error) {
	//insert แล้ว return ออกไปเป็นข้อมูล
	// '' ชื่อ field
	// "" ข้อมูลที่จะเอามาใส่ใน field
	//ใน select return ตาม User struct
	query := `
	SELECT
		json_build_object(
			'user', "t",
			'token', NULL
		)
	FROM (
		SELECT
			"u"."id",
			"u"."email",
			"u"."username",
			"u"."role_id"
		FROM "users" "u"
		WHERE "u"."id" = $1
	) AS "t"`

	data := make([]byte, 0)
	//args ผ่านกี่ตัวก็ได้
	if err := f.db.Get(&data, query, f.id); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	//เป็น strcut ของ user passort
	user := new(users.UserPassport)
	//unmarshal แปลง bytes เป็น json struct ที่จะ pass เข้าไป
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed: %v", err)
	}
	return user, nil
}
