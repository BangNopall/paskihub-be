package service

import (
	"context"
	"strconv"

	// "hash"
	// "mime/multipart"
	// "strconv"
	"time"

	// "github.com/gofiber/fiber/v2"
	// "github.com/google/uuid"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/constants"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/internal/infra/env"

	// "github.com/BangNopall/paskihub-be/internal/infra/env"
	// "github.com/BangNopall/paskihub-be/pkg/aws"
	"github.com/BangNopall/paskihub-be/pkg/bcrypt"
	"github.com/BangNopall/paskihub-be/pkg/gomail"
	"github.com/BangNopall/paskihub-be/pkg/helpers"
	"github.com/BangNopall/paskihub-be/pkg/log"

	// html_content "github.com/BangNopall/paskihub-be/pkg/html"
	"github.com/BangNopall/paskihub-be/pkg/jwt"
	// "github.com/BangNopall/paskihub-be/pkg/log"

	html_content "github.com/BangNopall/paskihub-be/pkg/html"
	"github.com/BangNopall/paskihub-be/pkg/redis"
	timePkg "github.com/BangNopall/paskihub-be/pkg/time"
	uuidPkg "github.com/BangNopall/paskihub-be/pkg/uuid"
)

type userService struct {
	userRepo contracts.UserRepository
	uuid     uuidPkg.UUIDInterface
	bcrypt   bcrypt.BcryptInterface
	time     timePkg.TimeInterface
	goMail   gomail.GoMailInterface
	jwt      jwt.JwtInterface
	redis    redis.RedisInterface
	timeout  time.Duration
}

func NewUserService(
	userRepo contracts.UserRepository,
	uuid uuidPkg.UUIDInterface,
	bcrypt bcrypt.BcryptInterface,
	time timePkg.TimeInterface,
	goMail gomail.GoMailInterface,
	jwt jwt.JwtInterface,
	redis redis.RedisInterface,
	timeout time.Duration,
) contracts.UserService {
	return &userService{
		userRepo,
		uuid,
		bcrypt,
		time,
		goMail,
		jwt,
		redis,
		timeout,
	}
}

func (s *userService) Register(ctx context.Context, user dto.UserRegister, referer string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if user.Password != user.ConfirmPassword {
		return domain.ErrConfirmPasswordNotMatch
	}

	hashPassword, err := s.bcrypt.Hash(user.Password)
	if err != nil {
		return err
	}

	emailVerPass := helpers.GenerateRandomString(10)

	emailVerPWhash, err := s.bcrypt.Hash(emailVerPass)
	if err != nil {
		return err
	}

	currentTime := s.time.Now()
	expiredToken := s.time.Add(time.Hour * 1)

	link := "https://paskihub.com/" + "auth/verify-email/" + user.Email + "/" + emailVerPass

	subject := "Verifikasi Email Akun Hology 8.0"
	HTMLbody := html_content.GetEmailVerifHTML(link)

	sendEmail := func(email string) <-chan error {
		errCh := make(chan error, 1)

		go func() {
			defer close(errCh)
			errCh <- s.goMail.SendEmail(subject, HTMLbody, email)
		}()

		return errCh
	}

	var registeredUser entity.User
	err = s.userRepo.FindUser(&registeredUser, &dto.UserParam{Email: user.Email})

	if err == domain.ErrInternalServer {
		return domain.ErrInternalServer
	}

	expiredTime := time.Date(
		registeredUser.ExpiredToken.Year(),
		registeredUser.ExpiredToken.Month(),
		registeredUser.ExpiredToken.Day(),
		registeredUser.ExpiredToken.Hour(),
		registeredUser.ExpiredToken.Minute(),
		registeredUser.ExpiredToken.Second(),
		registeredUser.ExpiredToken.Nanosecond(),
		time.Local)

	if currentTime.Before(expiredTime) && !registeredUser.EmailIsVerified && registeredUser.Email != "" {
		return domain.ErrCheckEmail
	}

	if currentTime.After(expiredTime) && !registeredUser.EmailIsVerified && registeredUser.Email != "" {
		updateUser := dto.UserUpdate{
			Password:           hashPassword,
			EmailVerifiedToken: emailVerPWhash,
			ExpiredToken:       expiredToken,
		}

		select {
		case <-ctx.Done():
			return domain.ErrTimeout
		case err := <-sendEmail(user.Email):
			if err != nil {
				return err
			}

			err = s.userRepo.UpdateUser(&updateUser, registeredUser.Id)
			if err != nil {
				return err
			}
		}

		return nil
	}

	uuid, err := s.uuid.New()
	if err != nil {
		return err
	}

	newUser := entity.User{
		Id:                 uuid,
		Email:              user.Email,
		Role:               constants.ROLE_PESERTA,
		Password:           hashPassword,
		EmailVerifiedToken: emailVerPWhash,
		ExpiredToken:       expiredToken,
	}

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	case err := <-sendEmail(user.Email):
		if err != nil {
			return err
		}
		err = s.userRepo.CreateUser(&newUser)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *userService) Login(ctx context.Context, user dto.UserLogin) (dto.UserLoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var registeredUser entity.User
	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{Email: user.Email})
	if err == domain.ErrNotFound {
		return dto.UserLoginResponse{}, domain.ErrWrongEmailOrPassword
	}

	if err == domain.ErrInternalServer {
		return dto.UserLoginResponse{}, domain.ErrInternalServer

	}

	valid := s.bcrypt.Compare(user.Password, registeredUser.Password)
	if !valid {
		return dto.UserLoginResponse{}, domain.ErrWrongEmailOrPassword

	}

	if !registeredUser.EmailIsVerified {
		return dto.UserLoginResponse{}, domain.ErrCheckEmail
	}

	token, err := s.jwt.GenerateToken(registeredUser.Id, registeredUser)
	if err != nil {
		return dto.UserLoginResponse{}, err
	}

	return dto.UserLoginResponse{Token: token}, nil
}

func (s *userService) DeleteUnverifiedUsers() {
	s.userRepo.DeleteUnverifiedUser()
	log.Info(nil, "[USER SERVICE][DeleteUnverifiedUsers] deleted unverified users")
}

func (s *userService) VerifyEmail(ctx context.Context, email string, emailVerPass string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var registeredUser entity.User
	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{Email: email})
	if err == domain.ErrNotFound {
		return domain.ErrWrongEmailOrPassword
	}

	if err == domain.ErrInternalServer {
		return domain.ErrInternalServer
	}

	valid := s.bcrypt.Compare(emailVerPass, registeredUser.EmailVerifiedToken)
	if !valid {
		return domain.ErrWrongEmailOrPassword
	}

	currentTime := s.time.Now()
	expiredTime := time.Date(
		registeredUser.ExpiredToken.Year(),
		registeredUser.ExpiredToken.Month(),
		registeredUser.ExpiredToken.Day(),
		registeredUser.ExpiredToken.Hour(),
		registeredUser.ExpiredToken.Minute(),
		registeredUser.ExpiredToken.Second(),
		registeredUser.ExpiredToken.Nanosecond(),
		time.Local)

	if currentTime.After(expiredTime) {
		return domain.ErrCheckEmail
	}

	updateUser := dto.UserUpdate{
		EmailIsVerified: true,
	}

	err = s.userRepo.UpdateUser(&updateUser, registeredUser.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) ResetPassword(ctx context.Context, user dto.UserResetPassword, forgotPasswordToken string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if user.Password != user.ConfirmPassword {
		return domain.ErrConfirmPasswordNotMatch
	}

	var registeredUser entity.User
	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{ForgotPasswordToken: forgotPasswordToken})
	if err != nil {
		return domain.ErrInvalidToken
	}

	expiredTime := time.Date(
		registeredUser.ExpiredTokenForgot.Year(),
		registeredUser.ExpiredTokenForgot.Month(),
		registeredUser.ExpiredTokenForgot.Day(),
		registeredUser.ExpiredTokenForgot.Hour(),
		registeredUser.ExpiredTokenForgot.Minute(),
		registeredUser.ExpiredTokenForgot.Second(),
		registeredUser.ExpiredTokenForgot.Nanosecond(),
		time.Local)

	currentTime := s.time.Now()

	if currentTime.After(expiredTime) {
		return domain.ErrInvalidToken
	}

	hashPassword, err := s.bcrypt.Hash(user.Password)
	if err != nil {
		return err
	}

	updateUser := dto.UserUpdate{
		Password:           hashPassword,
		ExpiredTokenForgot: time.Now(),
	}

	err = s.userRepo.UpdateUser(&updateUser, registeredUser.Id)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *userService) ForgotPassword(ctx context.Context, user dto.UserForgotPassword, referer string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var registeredUser entity.User
	err := s.userRepo.FindUser(&registeredUser, &dto.UserParam{Email: user.Email})
	if err == domain.ErrNotFound {
		return domain.ErrUserNotFound
	}

	if err == domain.ErrInternalServer {
		return domain.ErrInternalServer
	}

	forgotPasswordToken := helpers.GenerateRandomString(64)

	currentTime := s.time.Now()
	expiredToken := s.time.Add(time.Hour * 1)

	link := "https://paskihub.com" + `/auth/reset-password/` + forgotPasswordToken

	subject := "Forgot Password"
	HTMLbody := html_content.GetEmailForgotPassword(link)

	expiredTime := time.Date(
		registeredUser.ExpiredTokenForgot.Year(),
		registeredUser.ExpiredTokenForgot.Month(),
		registeredUser.ExpiredTokenForgot.Day(),
		registeredUser.ExpiredTokenForgot.Hour(),
		registeredUser.ExpiredTokenForgot.Minute(),
		registeredUser.ExpiredTokenForgot.Second(),
		registeredUser.ExpiredTokenForgot.Nanosecond(),
		time.Local)

	if currentTime.Before(expiredTime) {
		return domain.ErrCheckEmail
	}

	updateUser := dto.UserUpdate{
		ForgotPasswordToken: forgotPasswordToken,
		ExpiredTokenForgot:  expiredToken,
	}

	err = s.userRepo.UpdateUser(&updateUser, registeredUser.Id)
	if err != nil {
		return err
	}

	err = s.goMail.SendEmail(subject, HTMLbody, user.Email)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *userService) Logout(ctx context.Context, jwtToken string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	expTime := 1
	var err error

	if env.AppEnv.JwtExpireTime != "" {
		expTime, err = strconv.Atoi(env.AppEnv.JwtExpireTime)
	}

	if err != nil {
		return err
	}

	err = s.redis.Set(ctx, jwtToken, "LOGGED OUT", time.Hour*time.Duration(expTime))

	if err != nil {
		return err
	}

	return nil
}
