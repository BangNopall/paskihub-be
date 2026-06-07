package service

import (
	"context"
	"strconv"
	"time"

	"github.com/BangNopall/paskihub-be/domain"
	"github.com/BangNopall/paskihub-be/domain/constants"
	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/BangNopall/paskihub-be/internal/infra/env"

	"github.com/BangNopall/paskihub-be/pkg/bcrypt"
	"github.com/BangNopall/paskihub-be/pkg/gomail"
	"github.com/BangNopall/paskihub-be/pkg/helpers"
	"github.com/BangNopall/paskihub-be/pkg/log"
	"github.com/google/uuid"

	"github.com/BangNopall/paskihub-be/pkg/jwt"

	html_content "github.com/BangNopall/paskihub-be/pkg/html"
	"github.com/BangNopall/paskihub-be/pkg/redis"
	timePkg "github.com/BangNopall/paskihub-be/pkg/time"
	uuidPkg "github.com/BangNopall/paskihub-be/pkg/uuid"
)

type userService struct {
	userRepo  contracts.UserRepository
	eventRepo contracts.EventRepository
	uuid      uuidPkg.UUIDInterface
	bcrypt    bcrypt.BcryptInterface
	time      timePkg.TimeInterface
	goMail    gomail.GoMailInterface
	jwt       jwt.JwtInterface
	redis     redis.RedisInterface
	timeout   time.Duration
}

func NewUserService(
	userRepo contracts.UserRepository,
	eventRepo contracts.EventRepository,
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
		eventRepo,
		uuid,
		bcrypt,
		time,
		goMail,
		jwt,
		redis,
		timeout,
	}
}

func (s *userService) Register(ctx context.Context, role string, user dto.UserRegister) error {
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

	link := "https://paskihub.noxval.app" + "/auth/verify-email/" + user.Email + "/" + emailVerPass

	subject := "Verifikasi Email Akun PaskiHub"
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

	if role != constants.ROLE_PESERTA && role != constants.ROLE_EO {
		return domain.ErrInvalidRole
	} else if role == constants.ROLE_ADMIN {
		return domain.ErrInternalServer
	}

	newUser := entity.User{
		Id:                 uuid,
		Email:              user.Email,
		Role:               enums.Role(role),
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

	if registeredUser.IsBanned {
		return dto.UserLoginResponse{}, domain.ErrAccountBanned
	}

	if !registeredUser.EmailIsVerified {
		return dto.UserLoginResponse{}, domain.ErrCheckEmail
	}

	token, err := s.jwt.GenerateToken(registeredUser.Id, registeredUser)
	if err != nil {
		return dto.UserLoginResponse{}, err
	}

	// Update LastLoginAt
	now := s.time.Now()
	updateUser := dto.UserUpdate{
		LastLoginAt: &now,
	}
	_ = s.userRepo.UpdateUser(&updateUser, registeredUser.Id)

	var parentId *string
	if registeredUser.ParentId != nil {
		parentIdValue := registeredUser.ParentId.String()
		parentId = &parentIdValue
	}

	return dto.UserLoginResponse{
		Token:    token,
		Id:       registeredUser.Id.String(),
		Email:    registeredUser.Email,
		Role:     string(registeredUser.Role),
		ParentId: parentId,
	}, nil
}

func (s *userService) DeleteUnverifiedUsers() {
	s.userRepo.DeleteUnverifiedUser()
	log.Info(nil, "[USER SERVICE][DeleteUnverifiedUsers] deleted unverified users")
}

func (s *userService) CreateAdmin(ctx context.Context, req dto.AdminCreateRequest) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashPassword, err := s.bcrypt.Hash(req.Password)
	if err != nil {
		return err
	}

	userId, err := s.uuid.New()
	if err != nil {
		return err
	}

	newAdmin := entity.User{
		Id:              userId,
		Email:           req.Email,
		Password:        hashPassword,
		Role:            enums.Admin,
		EmailIsVerified: true,
	}

	err = s.userRepo.CreateUser(&newAdmin)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteAdmin(ctx context.Context, userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var user entity.User
	err := s.userRepo.FindUser(&user, &dto.UserParam{Id: userId})
	if err != nil {
		return err
	}

	if user.Role != enums.Admin {
		return domain.ErrForbidden
	}

	return s.userRepo.DeleteUser(ctx, userId)
}

func (s *userService) ResetAdminPassword(ctx context.Context, userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var user entity.User
	err := s.userRepo.FindUser(&user, &dto.UserParam{Id: userId})
	if err != nil {
		return err
	}

	if user.Role != enums.Admin {
		return domain.ErrForbidden
	}

	// Reuse ForgotPassword logic to send reset link
	return s.ForgotPassword(ctx, dto.UserForgotPassword{Email: user.Email})
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

func (s *userService) ForgotPassword(ctx context.Context, user dto.UserForgotPassword) error {
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

	link := "https://paskihub.noxval.app" + `/auth/forgot-password/` + forgotPasswordToken

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

func (s *userService) FetchAllUsers(ctx context.Context, role *string) ([]dto.UserResponse, error) {
	users, err := s.userRepo.FetchAllUsers(ctx, role)
	if err != nil {
		return nil, err
	}
	var res []dto.UserResponse
	for _, u := range users {
		res = append(res, *dto.UserEntityToResponse(&u))
	}
	return res, nil
}

func (s *userService) FetchUserById(ctx context.Context, userId uuid.UUID) (dto.UserResponse, error) {
	var user entity.User
	err := s.userRepo.FindUser(&user, &dto.UserParam{Id: userId})
	if err != nil {
		return dto.UserResponse{}, err
	}
	return *dto.UserEntityToResponse(&user), nil
}

func (s *userService) BanUser(ctx context.Context, userId uuid.UUID) error {
	return s.userRepo.UpdateUserStatus(ctx, userId, true)
}

func (s *userService) UnbanUser(ctx context.Context, userId uuid.UUID) error {
	return s.userRepo.UpdateUserStatus(ctx, userId, false)
}

func (s *userService) VerifyUser(ctx context.Context, userId uuid.UUID) error {
	return s.userRepo.VerifyUserEmail(ctx, userId)
}

func (s *userService) FetchUserDetailAdmin(ctx context.Context, userId uuid.UUID) (dto.AdminUserDetailResponse, error) {
	user, err := s.userRepo.FetchUserDetailAdmin(ctx, userId)
	if err != nil {
		return dto.AdminUserDetailResponse{}, err
	}

	userRes := dto.UserEntityToResponse(user)

	res := dto.AdminUserDetailResponse{
		Id:              user.Id,
		Name:            userRes.InstitutionName,
		Email:           user.Email,
		Role:            string(user.Role),
		Status:          "Active",
		EmailIsVerified: userRes.EmailIsVerified,
		IsBanned:        userRes.IsBanned,
		LastLoginAt:     userRes.LastLoginAt,
		JoinedAt:        user.CreatedAt,
		Event:           userRes.Event,
		Institution:     userRes.Institution,
	}

	if user.IsBanned {
		res.Status = "Banned"
	} else if !user.EmailIsVerified {
		res.Status = "Pending"
	}

	if len(user.Institutions) > 0 {
		inst := user.Institutions[0]
		res.SchoolName = inst.Name
		res.Phone = inst.NoWaPj
		res.Address = inst.Address
	}

	if user.Role == enums.Organizer {
		res.EOData = s.mapAdminUserEOData(user)
	} else if user.Role == enums.Peserta {
		res.PesertaData = s.mapAdminUserPesertaData(user)
	}

	return res, nil
}

func (s *userService) ArchiveOrganizer(ctx context.Context, userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var user entity.User
	err := s.userRepo.FindUser(&user, &dto.UserParam{Id: userId}, "Events")
	if err != nil {
		return err
	}

	if len(user.Events) == 0 {
		return domain.ErrNotFound
	}

	event := user.Events[0]
	update := entity.Event{
		Status: string(enums.Archived),
	}

	return s.eventRepo.UpdateEvent(ctx, &update, event.Id)
}

func (s *userService) UnarchiveOrganizer(ctx context.Context, userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var user entity.User
	err := s.userRepo.FindUser(&user, &dto.UserParam{Id: userId}, "Events")
	if err != nil {
		return err
	}

	if len(user.Events) == 0 {
		return domain.ErrNotFound
	}

	event := user.Events[0]
	update := entity.Event{
		Status: string(enums.Open),
	}

	return s.eventRepo.UpdateEvent(ctx, &update, event.Id)
}

func (s *userService) mapAdminUserEOData(user *entity.User) *dto.AdminUserEODetail {
	eoData := &dto.AdminUserEODetail{
		Panitia: []dto.AdminStaffRes{},
		Events:  []dto.AdminEventRes{},
		Judges:  []dto.AdminJudgeRes{},
	}

	for _, e := range user.Events {
		var levels []dto.AdminEventLevelRes
		for _, lvl := range e.EventLevels {
			levels = append(levels, dto.AdminEventLevelRes{
				Name:     lvl.Name,
				RegisFee: lvl.RegisFee,
				DpFee:    lvl.DpFee,
			})
		}
		eoData.Events = append(eoData.Events, dto.AdminEventRes{
			EventName:   e.Name,
			CompeDate:   e.CompeDate,
			Location:    e.Location,
			Status:      string(e.Status),
			EventLevels: levels,
		})
	}
	return eoData
}

func (s *userService) mapAdminUserPesertaData(user *entity.User) *dto.AdminUserPesertaDetail {
	pesertaData := &dto.AdminUserPesertaDetail{
		Teams:        []dto.AdminTeamRes{},
		EventHistory: []dto.AdminEventRegistrationRes{},
	}

	if len(user.Institutions) > 0 {
		for _, t := range user.Institutions[0].Teams {
			var members []dto.AdminTeamMemberRes
			for _, m := range t.TeamMembers {
				members = append(members, dto.AdminTeamMemberRes{
					Name: m.FullName,
					Role: string(m.Role),
				})
			}
			pesertaData.Teams = append(pesertaData.Teams, dto.AdminTeamRes{
				TeamName:     t.Name,
				Coach:        t.Pelatih,
				MembersCount: len(t.TeamMembers),
				Members:      members,
			})

			for _, reg := range t.Registrations {
				pesertaData.EventHistory = append(pesertaData.EventHistory, dto.AdminEventRegistrationRes{
					EventName:        reg.EventLevel.Event.Name,
					EventLevelName:   reg.EventLevel.Name,
					PaymentStatus:    string(reg.PaymentStatus),
					AssessmentStatus: string(reg.AssessmentStatus),
				})
			}
		}
	}
	return pesertaData
}
