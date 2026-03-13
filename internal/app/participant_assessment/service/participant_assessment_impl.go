package service

import (
	"context"
	"errors"

	"github.com/BangNopall/paskihub-be/domain/contracts"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"github.com/google/uuid"
)

type participantAssessmentService struct {
	repo contracts.ParticipantAssessmentRepository
}

func NewParticipantAssessmentService(repo contracts.ParticipantAssessmentRepository) contracts.ParticipantAssessmentService {
	return &participantAssessmentService{
		repo: repo,
	}
}

func (s *participantAssessmentService) GetAssessmentRecap(ctx context.Context, userID string, regisID string) (*dto.AssessmentRecapResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	parsedRegisID, err := uuid.Parse(regisID)
	if err != nil {
		return nil, errors.New("invalid regis id")
	}

	isOwner, err := s.repo.GetRegistrationOwnership(ctx, parsedRegisID, parsedUserID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, errors.New("unauthorized to view this assessment")
	}

	regis, err := s.repo.GetRegistrationByID(ctx, parsedRegisID)
	if err != nil {
		return nil, err
	}

	eventID := regis.EventLevel.Event.Id
	categories, err := s.repo.GetScoreCategoriesByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	scores, err := s.repo.GetScoresByRegisID(ctx, parsedRegisID)
	if err != nil {
		return nil, err
	}

	violations, err := s.repo.GetViolationsByRegisID(ctx, parsedRegisID)
	if err != nil {
		return nil, err
	}

	var totalScore float64
	var totalViolationPoints float64

	scoreMap := make(map[string][]dto.JudgeScoreDetail)
	for _, sc := range scores {
		totalScore += sc.ScoreValue
		subCatIDStr := sc.SubCategoryID.String()
		scoreMap[subCatIDStr] = append(scoreMap[subCatIDStr], dto.JudgeScoreDetail{
			JudgeName: sc.Judge.Name,
			Value:     sc.ScoreValue,
		})
	}

	var violationDetails []dto.AssessmentViolationDetail
	for _, v := range violations {
		totalViolationPoints += v.ViolationType.Point
		violationDetails = append(violationDetails, dto.AssessmentViolationDetail{
			Name:      v.ViolationType.Name,
			Point:     v.ViolationType.Point,
			JudgeName: v.Judge.Name,
		})
	}

	var categoryDetails []dto.AssessmentCategoryDetail
	for _, c := range categories {
		catDetail := dto.AssessmentCategoryDetail{
			CategoryName: c.Name,
		}

		for _, subC := range c.SubCategories {
			subCDetail := dto.AssessmentSubCategoryDetail{
				Name:   subC.Name,
				Scores: scoreMap[subC.ID.String()],
			}
			catDetail.SubCategories = append(catDetail.SubCategories, subCDetail)
		}

		categoryDetails = append(categoryDetails, catDetail)
	}

	finalScore := totalScore - totalViolationPoints

	return &dto.AssessmentRecapResponse{
		TotalScore:           totalScore,
		TotalViolationPoints: totalViolationPoints,
		FinalScore:           finalScore,
		Violations:           violationDetails,
		Categories:           categoryDetails,
	}, nil
}
