# Ranking System (Event Award) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement CRUD API for Event Award configurations with Many-to-Many relationship to Score Categories.

**Architecture:** Use a standard GORM Many-to-Many relationship with a join table (`event_award_score_categories`). Implementation follows the existing feature-based architecture in the `assessment` module.

**Tech Stack:** Go (Golang), Fiber v2, GORM, PostgreSQL.

---

### Task 1: Domain Entities & DTOs

**Files:**
- Modify: `domain/entity/assessment.go`
- Modify: `domain/dto/assessment_dto.go`

- [ ] **Step 1: Add EventAward entity**
Update `domain/entity/assessment.go` to include the `EventAward` struct and its relationship.

```go
type EventAward struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID         uuid.UUID       `gorm:"type:uuid;not null"`
	EventLevelID    uuid.UUID       `gorm:"type:uuid;not null"`
	Name            string          `gorm:"type:varchar(255);not null"`
	LimitRank       int             `gorm:"not null;default:1"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relation
	Event           Event           `gorm:"foreignKey:EventID;references:Id"`
	EventLevel      EventLevel      `gorm:"foreignKey:EventLevelID;references:Id"`
	ScoreCategories []ScoreCategory `gorm:"many2many:event_award_score_categories;"`
}
```

- [ ] **Step 2: Add Event Award DTOs**
Update `domain/dto/assessment_dto.go` with request and response DTOs.

```go
type CreateAwardReq struct {
	EventLevelID     uuid.UUID   `json:"event_level_id" validate:"required"`
	Name             string      `json:"name" validate:"required"`
	LimitRank        int         `json:"limit_rank" validate:"required,min=1"`
	ScoreCategoryIDs []uuid.UUID `json:"score_category_ids" validate:"required,min=1"`
}

type UpdateAwardReq struct {
	Name             string      `json:"name" validate:"required"`
	LimitRank        int         `json:"limit_rank" validate:"required,min=1"`
	ScoreCategoryIDs []uuid.UUID `json:"score_category_ids" validate:"required,min=1"`
}

type AwardRes struct {
	ID              uuid.UUID          `json:"id"`
	EventID         uuid.UUID          `json:"event_id"`
	EventLevelID    uuid.UUID          `json:"event_level_id"`
	Name            string             `json:"name"`
	LimitRank       int                `json:"limit_rank"`
	ScoreCategories []ScoreCategoryRes `json:"score_categories,omitempty"`
}
```

- [ ] **Step 3: Commit Domain Changes**
```bash
git add domain/entity/assessment.go domain/dto/assessment_dto.go
git commit -m "feat(domain): add EventAward entity and DTOs"
```

---

### Task 2: Repository Implementation

**Files:**
- Modify: `domain/contracts/assessment_contract.go`
- Modify: `internal/app/assessment/repository/assessment_repository.go`

- [ ] **Step 1: Update IAssessmentRepository contract**
Add new methods to `domain/contracts/assessment_contract.go`.

```go
type IAssessmentRepository interface {
    // ... existing methods
    CreateAward(ctx context.Context, award *entity.EventAward) error
    GetAwardsByEvent(ctx context.Context, eventId uuid.UUID, levelId *uuid.UUID) ([]entity.EventAward, error)
    FindAwardById(ctx context.Context, id uuid.UUID) (*entity.EventAward, error)
    UpdateAward(ctx context.Context, award *entity.EventAward, categoryIds []uuid.UUID) error
    DeleteAward(ctx context.Context, id uuid.UUID) error
}
```

- [ ] **Step 2: Implement Repository Methods**
Update `internal/app/assessment/repository/assessment_repository.go`.

```go
func (r *assessmentRepository) CreateAward(ctx context.Context, award *entity.EventAward) error {
	return r.db.WithContext(ctx).Create(award).Error
}

func (r *assessmentRepository) GetAwardsByEvent(ctx context.Context, eventId uuid.UUID, levelId *uuid.UUID) ([]entity.EventAward, error) {
	var awards []entity.EventAward
	query := r.db.WithContext(ctx).Preload("ScoreCategories").Where("event_id = ?", eventId)
	if levelId != nil {
		query = query.Where("event_level_id = ?", *levelId)
	}
	err := query.Find(&awards).Error
	return awards, err
}

func (r *assessmentRepository) FindAwardById(ctx context.Context, id uuid.UUID) (*entity.EventAward, error) {
	var award entity.EventAward
	err := r.db.WithContext(ctx).Preload("ScoreCategories").First(&award, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &award, nil
}

func (r *assessmentRepository) UpdateAward(ctx context.Context, award *entity.EventAward, categoryIds []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(award).Error; err != nil {
			return err
		}
		var categories []entity.ScoreCategory
		if err := tx.Where("id IN ?", categoryIds).Find(&categories).Error; err != nil {
			return err
		}
		return tx.Model(award).Association("ScoreCategories").Replace(categories)
	})
}

func (r *assessmentRepository) DeleteAward(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.EventAward{}, "id = ?", id).Error
}
```

- [ ] **Step 3: Commit Repository Changes**
```bash
git add domain/contracts/assessment_contract.go internal/app/assessment/repository/assessment_repository.go
git commit -m "feat(repo): implement EventAward repository methods"
```

---

### Task 3: Service Logic

**Files:**
- Modify: `domain/contracts/assessment_contract.go`
- Modify: `internal/app/assessment/service/assessment_service.go`

- [ ] **Step 1: Update IAssessmentService contract**
Add methods to `domain/contracts/assessment_contract.go`.

```go
type IAssessmentService interface {
    // ... existing methods
    CreateAward(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateAwardReq) (*dto.AwardRes, error)
    GetAwards(ctx context.Context, eventId, userId uuid.UUID, levelId *uuid.UUID) ([]dto.AwardRes, error)
    UpdateAward(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateAwardReq) (*dto.AwardRes, error)
    DeleteAward(ctx context.Context, eventId, userId, id uuid.UUID) error
}
```

- [ ] **Step 2: Implement Service Methods**
Update `internal/app/assessment/service/assessment_service.go`.

```go
func (s *assessmentService) CreateAward(ctx context.Context, eventId, userId uuid.UUID, req dto.CreateAwardReq) (*dto.AwardRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}

	var categories []entity.ScoreCategory
	for _, catId := range req.ScoreCategoryIDs {
		categories = append(categories, entity.ScoreCategory{ID: catId})
	}

	award := &entity.EventAward{
		EventID:         eventId,
		EventLevelID:    req.EventLevelID,
		Name:            req.Name,
		LimitRank:       req.LimitRank,
		ScoreCategories: categories,
	}

	if err := s.repo.CreateAward(ctx, award); err != nil {
		return nil, domain.ErrInternalServer
	}

	return s.mapToAwardRes(award), nil
}

func (s *assessmentService) GetAwards(ctx context.Context, eventId, userId uuid.UUID, levelId *uuid.UUID) ([]dto.AwardRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	awards, err := s.repo.GetAwardsByEvent(ctx, eventId, levelId)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	var res []dto.AwardRes
	for _, a := range awards {
		res = append(res, *s.mapToAwardRes(&a))
	}
	return res, nil
}

func (s *assessmentService) UpdateAward(ctx context.Context, eventId, userId, id uuid.UUID, req dto.UpdateAwardReq) (*dto.AwardRes, error) {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return nil, err
	}
	award, err := s.repo.FindAwardById(ctx, id)
	if err != nil || award.EventID != eventId {
		return nil, domain.ErrNotFound
	}

	award.Name = req.Name
	award.LimitRank = req.LimitRank

	if err := s.repo.UpdateAward(ctx, award, req.ScoreCategoryIDs); err != nil {
		return nil, domain.ErrInternalServer
	}

	// Re-fetch to get updated relations
	updatedAward, _ := s.repo.FindAwardById(ctx, id)
	return s.mapToAwardRes(updatedAward), nil
}

func (s *assessmentService) DeleteAward(ctx context.Context, eventId, userId, id uuid.UUID) error {
	if err := s.ensureOwnership(ctx, eventId, userId); err != nil {
		return err
	}
	award, err := s.repo.FindAwardById(ctx, id)
	if err != nil || award.EventID != eventId {
		return domain.ErrNotFound
	}
	return s.repo.DeleteAward(ctx, id)
}

func (s *assessmentService) mapToAwardRes(a *entity.EventAward) *dto.AwardRes {
	var cats []dto.ScoreCategoryRes
	for _, c := range a.ScoreCategories {
		cats = append(cats, dto.ScoreCategoryRes{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	return &dto.AwardRes{
		ID:              a.ID,
		EventID:         a.EventID,
		EventLevelID:    a.EventLevelID,
		Name:            a.Name,
		LimitRank:       a.LimitRank,
		ScoreCategories: cats,
	}
}
```

- [ ] **Step 3: Commit Service Changes**
```bash
git add internal/app/assessment/service/assessment_service.go
git commit -m "feat(service): implement EventAward service logic"
```

---

### Task 4: Controller & Routes

**Files:**
- Modify: `internal/app/assessment/controller/assessment_controller.go`

- [ ] **Step 1: Implement Controller Methods**
Add handlers to `internal/app/assessment/controller/assessment_controller.go`.

```go
func (c *assessmentController) CreateAward(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusCreated
		res     interface{}
		message string = "failed to create award"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.CreateAwardReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.CreateAward(ctx.Context(), eventId, userId, req)
	code = domain.GetCode(err)
	if err == nil { message = "success to create award"; code = http.StatusCreated }
	return nil
}

func (c *assessmentController) GetAwards(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to get awards"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }

	var levelIdPtr *uuid.UUID
	if lid := ctx.Query("level_id"); lid != "" {
		if u, err := uuid.Parse(lid); err == nil {
			levelIdPtr = &u
		}
	}

	res, err = c.svc.GetAwards(ctx.Context(), eventId, userId, levelIdPtr)
	code = domain.GetCode(err)
	if err == nil { message = "success to get awards" }
	return nil
}

func (c *assessmentController) UpdateAward(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		res     interface{}
		message string = "failed to update award"
	)
	sendResp := func() { response.Send(ctx, code, message, res, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	var req dto.UpdateAwardReq
	if errParse := ctx.BodyParser(&req); errParse != nil {
		err = errParse; code = http.StatusBadRequest; return nil
	}

	res, err = c.svc.UpdateAward(ctx.Context(), eventId, userId, id, req)
	code = domain.GetCode(err)
	if err == nil { message = "success to update award" }
	return nil
}

func (c *assessmentController) DeleteAward(ctx *fiber.Ctx) error {
	var (
		err     error
		code    int = http.StatusOK
		message string = "failed to delete award"
	)
	sendResp := func() { response.Send(ctx, code, message, nil, err) }
	defer sendResp()

	userId, err := getUserId(ctx)
	if err != nil { code = http.StatusUnauthorized; return nil }
	eventId, err := getUUIDParam(ctx, "eventId")
	if err != nil { code = http.StatusBadRequest; return nil }
	id, err := getUUIDParam(ctx, "id")
	if err != nil { code = http.StatusBadRequest; return nil }

	err = c.svc.DeleteAward(ctx.Context(), eventId, userId, id)
	code = domain.GetCode(err)
	if err == nil { message = "success to delete award" }
	return nil
}
```

- [ ] **Step 2: Register Routes**
Update `InitAssessmentController` in `internal/app/assessment/controller/assessment_controller.go`.

```go
func InitAssessmentController(...) {
    // ... existing routes
    
    // Event Awards
    group.Post("/awards", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.CreateAward)
    group.Get("/awards", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.GetAwards)
    group.Put("/awards/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.UpdateAward)
    group.Delete("/awards/:id", middleware.Authentication, middleware.RateLimiter(), middleware.AuthOrganizer, c.DeleteAward)
}
```

- [ ] **Step 3: Commit Controller Changes**
```bash
git add internal/app/assessment/controller/assessment_controller.go
git commit -m "feat(controller): add EventAward endpoints"
```

---

### Task 5: Verification

- [ ] **Step 1: Run Air and verify migrations**
Run `air` to start the server. Check logs to ensure `event_awards` and `event_award_score_categories` tables are created.

- [ ] **Step 2: Manual API Test (Postman/Curl)**
Test the CRUD flow:
1. POST `/api/v1/eo/events/{eventId}/assessment/awards`
2. GET `/api/v1/eo/events/{eventId}/assessment/awards`
3. PUT `/api/v1/eo/events/{eventId}/assessment/awards/{id}`
4. DELETE `/api/v1/eo/events/{eventId}/assessment/awards/{id}`
