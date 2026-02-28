package repository

import (
	"github.com/zcicd/zcicd-server/internal/deploy/model"
	"gorm.io/gorm"
)

type EnvRepository struct {
	db *gorm.DB
}

func NewEnvRepository(db *gorm.DB) *EnvRepository {
	return &EnvRepository{db: db}
}

// Variables

func (r *EnvRepository) ListVariables(envID string) ([]model.EnvVariable, error) {
	var vars []model.EnvVariable
	err := r.db.Where("environment_id = ?", envID).Order("var_key").Find(&vars).Error
	return vars, err
}

func (r *EnvRepository) GetVariable(id string) (*model.EnvVariable, error) {
	var v model.EnvVariable
	err := r.db.Where("id = ?", id).First(&v).Error
	return &v, err
}

func (r *EnvRepository) CreateVariable(v *model.EnvVariable) error {
	return r.db.Create(v).Error
}

func (r *EnvRepository) UpdateVariable(v *model.EnvVariable) error {
	return r.db.Save(v).Error
}

func (r *EnvRepository) DeleteVariable(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.EnvVariable{}).Error
}

func (r *EnvRepository) BatchUpsertVariables(envID string, vars []model.EnvVariable) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, v := range vars {
			v.EnvironmentID = envID
			result := tx.Where("environment_id = ? AND var_key = ?", envID, v.VarKey).
				Assign(model.EnvVariable{VarValue: v.VarValue, IsSecret: v.IsSecret, Description: v.Description}).
				FirstOrCreate(&v)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

// Resource Quotas

func (r *EnvRepository) GetQuota(envID string) (*model.EnvResourceQuota, error) {
	var q model.EnvResourceQuota
	err := r.db.Where("environment_id = ?", envID).First(&q).Error
	return &q, err
}

func (r *EnvRepository) UpsertQuota(q *model.EnvResourceQuota) error {
	var existing model.EnvResourceQuota
	err := r.db.Where("environment_id = ?", q.EnvironmentID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(q).Error
	}
	if err != nil {
		return err
	}
	q.ID = existing.ID
	return r.db.Save(q).Error
}
