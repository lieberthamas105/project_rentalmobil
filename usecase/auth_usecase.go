package usecase

import "project_rentalmobil/utils/service"

type AuthenticationUseCase interface {
	Login(username string, password string) (string, error)
}

type authenticationUseCase struct {
	userUseCase UserUseCase
	jwtService  service.JwtService
}

func (a *authenticationUseCase) Login(username string, password string) (string, error) {
	user, err := a.userUseCase.FindUserByUsernamePassword(username, password)
	if err != nil {
		return "", err
	}

	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
func NewAuthenticationUseCase(uc UserUseCase, jwtService service.JwtService) AuthenticationUseCase {
	return &authenticationUseCase{userUseCase: uc, jwtService: jwtService}
}
