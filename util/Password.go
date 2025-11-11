package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// GenerateFromPassword 会自动生成盐值
	// 第二个参数是 cost（计算成本），值越大越安全但越慢，推荐 10-14
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
