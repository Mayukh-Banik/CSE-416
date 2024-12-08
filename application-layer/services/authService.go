package services

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (as *AuthService) Login(username, password string) bool {
	// 로그인 로직 구현 (예: DB 검증)
	return username == "admin" && password == "password"
}

func (as *AuthService) Register(username, password string) bool {
	// 회원가입 로직 구현
	return true
}
