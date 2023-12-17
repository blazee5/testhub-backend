package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/user"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
)

type Service struct {
	log       *zap.SugaredLogger
	repo      user.Repository
	redisRepo user.RedisRepository
	awsRepo   user.AWSRepository
	tracer    trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo user.Repository, redisRepo user.RedisRepository, awsRepo user.AWSRepository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, redisRepo: redisRepo, awsRepo: awsRepo, tracer: tracer}
}

func (s *Service) GetByID(ctx context.Context, userID int) (models.UserInfo, error) {
	ctx, span := s.tracer.Start(ctx, "userService.GetByID")
	defer span.End()

	cachedUser, err := s.redisRepo.GetByIDCtx(ctx, strconv.Itoa(userID))

	if err != nil {
		s.log.Infof("cannot get user by id in redis: %v", err)
	}

	if cachedUser != nil {
		return *cachedUser, nil
	}

	user, err := s.repo.GetByID(ctx, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.UserInfo{}, err
	}

	if err := s.redisRepo.SetUserCtx(ctx, strconv.Itoa(user.User.ID), 600, &user); err != nil {
		s.log.Infof("error while save user to cache: %v", err)
	}

	return user, nil
}

func (s *Service) ChangeAvatar(ctx context.Context, userID int, fileHeader *multipart.FileHeader) error {
	ctx, span := s.tracer.Start(ctx, "userService.ChangeAvatar")
	defer span.End()

	avatar, err := s.repo.GetAvatarByID(ctx, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	file, err := fileHeader.Open()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	bytes, err := io.ReadAll(file)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	contentType := http.DetectContentType(bytes)

	uuid, err := uuid.NewUUID()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if avatar == "" {
		avatar = uuid.String() + filepath.Ext(fileHeader.Filename)
	}

	err = s.awsRepo.SaveFile(ctx, avatar, contentType, bytes)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.ChangeAvatar(ctx, userID, "avatar/"+avatar)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if err = s.redisRepo.DeleteUserCtx(ctx, strconv.Itoa(userID)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) Update(ctx context.Context, userID int, input domain.UpdateUser) error {
	ctx, span := s.tracer.Start(ctx, "userService.Update")
	defer span.End()

	err := s.repo.Update(ctx, userID, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if err = s.redisRepo.DeleteUserCtx(ctx, strconv.Itoa(userID)); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, userID int) error {
	ctx, span := s.tracer.Start(ctx, "userService.Delete")
	defer span.End()

	if err := s.repo.Delete(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}
