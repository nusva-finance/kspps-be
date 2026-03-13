package repositories

import (
	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

type MemberRepository struct {
	db *gorm.DB
}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{
		db: config.GetDB(),
	}
}

func (r *MemberRepository) Create(member *models.Member) error {
	return r.db.Create(member).Error
}

func (r *MemberRepository) FindByID(id uint) (*models.Member, error) {
	var member models.Member
	err := r.db.First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepository) FindByMemberNo(memberNo string) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("member_no = ?", memberNo).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepository) FindByNIK(nik string) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("ktp_no = ?", nik).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepository) FindByKtpNo(ktpNo string) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("ktp_no = ?", ktpNo).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepository) List(offset, limit int) ([]models.Member, int64, error) {
	var members []models.Member
	var total int64

	err := r.db.Model(&models.Member{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

func (r *MemberRepository) Update(member *models.Member) error {
	return r.db.Save(member).Error
}

func (r *MemberRepository) Delete(id uint) error {
	return r.db.Delete(&models.Member{}, id).Error
}

func (r *MemberRepository) Search(keyword string, offset, limit int) ([]models.Member, int64, error) {
	var members []models.Member
	var total int64

	query := r.db.Model(&models.Member{}).Where("is_active = ?", true)

	if keyword != "" {
		query = query.Where(
			"member_no ILIKE ? OR full_name ILIKE ? OR ktp_no ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
		)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Find(&members).Error
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}



func (r *MemberRepository) FindByKtp(ktpNo string) (*models.Member, error) {
	var member models.Member
	// Mencari berdasarkan kolom ktp_no sesuai di models/member.go
	err := r.db.Where("ktp_no = ?", ktpNo).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

