package usecase

import (
	"context"
	"hexagony/internal/users/domain"

	"github.com/google/uuid"
)

type userUseCase struct {
	userRepository domain.UserRepository
}

func NewUserUseCase(ur domain.UserRepository) domain.UserUseCase {
	return &userUseCase{userRepository: ur}
}

func (u *userUseCase) FindAll(ctx context.Context) ([]*domain.UserList, error) {
	user, err := u.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) FindByID(ctx context.Context, uuid uuid.UUID) (*domain.UserList, error) {
	user, err := u.userRepository.FindByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Add(ctx context.Context, user *domain.User) error {
	if err := u.userRepository.Add(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) Update(ctx context.Context, uuid uuid.UUID, user *domain.User) error {
	if err := u.userRepository.Update(ctx, uuid, user); err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) Delete(ctx context.Context, uuid uuid.UUID) error {
	if err := u.userRepository.Delete(ctx, uuid); err != nil {
		return err
	}
	return nil
}
