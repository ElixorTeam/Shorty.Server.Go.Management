package v1

type UserDetails struct {
	Roles      []string
	UserId     string
	Email      string
	Username   string
	Name       string
	FamilyName string
}

type UserUsecase interface {
}

type UserRepository interface {
}
