package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `db:"id" json:"id"` //* query ผ่าน struct เอาข้อมูลเข้าไปใน struct ทั้งก้อน
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	RoleId   int    `db:"role_id" json:"role_id"`
}

type UserRegisterRequest struct {
	//เพิ่ม form เนื่องจากใช้ x-www-form ใน postman ซึ่งปลอดภัยกว่า raw json
	Email    string `db:"email" json:"emai" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
	Username string `db:"username" json:"username" form:"username"`
}

// ไม่เก็บรหัส password ของ user โดยเข้ารหัสไม่ให้ admin เห็น
// Vvvvvเป็น method ไปใช้กับ struct ตัวไหนเอามาไว้ข้างหน้า
func (obj *UserRegisterRequest) BcryptHashing() error {
	// เลข 10 ยิงเยอะยิงแฮ็คยาก แต่เวลาถอดรหัสจะใช้ทรัพยากรเยอะ
	hashedPasswords, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hashed password failed: %v", err)
	}
	//112345 -> djfhfuyrhjsddh
	obj.Password = string(hashedPasswords)
	return nil
}

// ทำ sign in
type UserCredentials struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserCredentialChecker struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
	RoleId   int    `db:"role_id"`
}

// ฟังก์ชั่นเช็ค pattern email
func (obj *UserRegisterRequest) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
}

type UserPassport struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}

type UserToken struct {
	Id           string `db:"id" json:"id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
}

////////------------------------------////////////////

type UserClaims struct {
	//*ห้ามเอาข้อมุล sensitive เข้ามา
	Id     string `db:"id" json:"id"`
	RoleId int    `db:"role" json:"role"`
}

type UserRefreshCredential struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type Oauth struct {
	Id     string `db:"id" json:"id"`
	UserId string `db:"user_id" json:"user_id"`
}

type UserRemoveCredential struct {
	OauthId string `json:"oauth_id" form:"oauth_id"`
}
