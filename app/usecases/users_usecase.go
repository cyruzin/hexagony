package usecase

import (
	"context"
	"hexagony/app/domain"

	"github.com/google/uuid"
)

type userUseCase struct {
	userRepository domain.UsersRepository
}

func NewUserUseCase(ur domain.UsersRepository) domain.UsersUseCase {
	return &userUseCase{userRepository: ur}
}

func (u *userUseCase) FindAll(ctx context.Context) ([]*domain.UsersList, error) {
	user, err := u.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) FindByID(ctx context.Context, uuid uuid.UUID) (*domain.UsersList, error) {
	user, err := u.userRepository.FindByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Add(ctx context.Context, user *domain.Users) error {
	if err := u.userRepository.Add(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) Update(ctx context.Context, uuid uuid.UUID, user *domain.Users) error {
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
