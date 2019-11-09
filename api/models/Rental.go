package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Rental struct {
	ID             uint32    `gorm:"primary_key;auto_increment" json:"id"`
	CarID          uint32    `gorm:"int" json:"Car_id"`
	Car            Car       `json:"Car"`
	UserID         uint32    `gorm:"int" json:"user_id"`
	User           User      `json:"user"`
	PickupLocation uint32    `gorm:"int" json:"pickup_location"`
	StartDate      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *Rental) Prepare() {
	u.ID = 0
	//u.StartDate = html.EscapeString(strings.TrimSpace(u.StartDate))
	//u.EndDate = html.EscapeString(strings.TrimSpace(u.EndDate))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (p *Rental) Validate() error {

	if p.StartDate.IsZero() {
		return errors.New("Required Start Date")
	}
	if p.EndDate.IsZero() {
		return errors.New("Required End Date")
	}

	return nil
}

func (p *Rental) SaveRental(db *gorm.DB) (*Rental, error) {
	var err error
	err = db.Debug().Model(&Rental{}).Create(&p).Error
	if err != nil {
		return &Rental{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Rental{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Rental{}, err
		}
	}
	return p, nil
}

func (p *Rental) FindAllRentals(db *gorm.DB) (*[]Rental, error) {
	var err error
	posts := []Rental{}
	err = db.Debug().Model(&Rental{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Rental{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].UserID).Take(&posts[i].User).Error
			if err != nil {
				return &[]Rental{}, err
			}
		}
	}
	return &posts, nil
}

func (p *Rental) FindRentalByID(db *gorm.DB, pid uint64) (*Rental, error) {
	var err error
	err = db.Debug().Model(&Rental{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Rental{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Rental{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &Rental{}, err
		}
	}
	return p, nil
}

func (p *Rental) UpdateARental(db *gorm.DB, pid uint64) (*Rental, error) {

	var err error
	db = db.Debug().Model(&Rental{}).Where("id = ?", pid).Take(&Rental{}).UpdateColumns(
		map[string]interface{}{
			"car_id":          p.CarID,
			"user_id":         p.UserID,
			"pick_uplocation": p.PickupLocation,
			"start_date":      p.StartDate,
			"end_date":        p.EndDate,
			"updated_at":      time.Now(),
		},
	)
	err = db.Debug().Model(&Rental{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Rental{}, err
	}
	/*if p.ID != 0 { LIVIU AREGLA ESTO, cuando haces un UPDATE no guarda bien el owner.Aun que  no lo necesitas
		err = db.Debug().Model(&User{}).Where("id = ?", p.User_id).Take(&p.Owner).Error
		if err != nil {
			return &Rental{}, err
		}
	}*/
	return p, nil
}

func (p *Rental) DeleteARental(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Rental{}).Where("id = ? and user_id = ?", pid, uid).Take(&Rental{}).Delete(&Rental{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Rental not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
