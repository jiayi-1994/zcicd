package repository

import (
	"context"

	"github.com/zcicd/zcicd-server/internal/auth/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := r.db.WithContext(ctx).Model(&model.User{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *UserRepository) GetUserRoles(ctx context.Context, userID string) ([]model.UserRole, error) {
	var roles []model.UserRole
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *UserRepository) AddUserRole(ctx context.Context, role *model.UserRole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *UserRepository) RemoveUserRole(ctx context.Context, userID, role, scopeType, scopeID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role = ? AND scope_type = ? AND scope_id = ?", userID, role, scopeType, scopeID).
		Delete(&model.UserRole{}).Error
}
