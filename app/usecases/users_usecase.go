package usecase

import (
	"context"
	"hexagony/app/domain"

	"github.com/google/uuid"
)

type usersUseCase struct {
	usersRepository domain.UsersRepository
}

func NewUserUseCase(ur domain.UsersRepository) domain.UsersUseCase {
	return &usersUseCase{usersRepository: ur}
}

func (u *usersUseCase) FindAll(ctx context.Context) ([]*domain.UsersList, error) {
	user, err := u.usersRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usersUseCase) FindByID(ctx context.Context, uuid uuid.UUID) (*domain.UsersList, error) {
	user, err := u.usersRepository.FindByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *usersUseCase) Add(ctx context.Context, user *domain.Users) error {
	if err := u.usersRepository.Add(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *usersUseCase) Update(ctx context.Context, uuid uuid.UUID, user *domain.Users) error {
	if err := u.usersRepository.Update(ctx, uuid, user); err != nil {
		return err
	}
	return nil
}

func (u *usersUseCase) Delete(ctx context.Context, uuid uuid.UUID) error {
	if err := u.usersRepository.Delete(ctx, uuid); err != nil {
		return err
	}
	return nil
}
