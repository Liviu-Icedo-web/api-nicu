package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Rental struct {
	ID              uint32        `gorm:"primary_key;auto_increment" json:"id"`
	CarID           uint32        `gorm:"int" json:"Car_id"`
	Car             Car           `json:"car"`
	Owner           User          `json:"user"`
	UserID          uint32        `gorm:"int" json:"user_id"`
	CarLocation     []CarLocation `json:"car_location"`
	PickupLocation  uint32        `gorm:"int" json:"pickup_location"`
	DropoffLocation uint32        `gorm:"int" json:"dropoff_location"`
	StartDate       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"start_date"`
	EndDate         time.Time     `json:"end_date"`
	CreatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
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
	return p, nil
}

func (p *Rental) FindUserRentals(db *gorm.DB) (*[]Rental, error) {
	var err error
	posts := []Rental{}
	err = db.Debug().Model(&Rental{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Rental{}, err
	}
	return &posts, nil
}

func (p *Rental) FindRentalUserID(db *gorm.DB, pid uint64) (*[]Rental, error) {
	var err error
	r := []Rental{}
	err = db.Debug().Model(&Rental{}).Where("user_id = ?", pid).Find(&r).Error
	if err != nil {
		return &[]Rental{}, err
	}
	if len(r) > 0 {
		for i, _ := range r {
			err := db.Debug().Model(&User{}).Where("id = ?", r[i].UserID).Take(&r[i].Owner).Error
			if err != nil {
				return &[]Rental{}, err
			}
			err = db.Debug().Model(&CarLocation{}).Where("car_id = ?", r[i].CarID).Find(&r[i].CarLocation).Error
			if err != nil {
				return &[]Rental{}, err
			}
			err = db.Debug().Model(&Car{}).Where("id = ?", r[i].CarID).Find(&r[i].Car).Error
			if err != nil {
				return &[]Rental{}, err
			}

		}
	}

	return &r, nil
}
func (p *Rental) FindRentalCarID(db *gorm.DB, pid uint64) (*[]Rental, error) {
	var err error
	r := []Rental{}
	err = db.Debug().Model(&Rental{}).Where("car_id = ?", pid).Find(&r).Error
	if err != nil {
		return &[]Rental{}, err
	}
	if len(r) > 0 {
		for i, _ := range r {
			err = db.Debug().Model(&Car{}).Where("id = ?", r[i].CarID).Find(&r[i].Car).Error
			if err != nil {
				return &[]Rental{}, err
			}

		}
	}
	fmt.Println("*** FindRenatlCar:", r)

	return &r, nil
}

func (p *Rental) UpdateARental(db *gorm.DB, pid uint64) (*Rental, error) {

	var err error
	db = db.Debug().Model(&Rental{}).Where("id = ?", pid).Take(&Rental{}).UpdateColumns(
		map[string]interface{}{
			"car_id":           p.CarID,
			"user_id":          p.UserID,
			"pickup_location":  p.PickupLocation,
			"dropoff_location": p.DropoffLocation,
			"start_date":       p.StartDate,
			"end_date":         p.EndDate,
			"updated_at":       time.Now(),
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
