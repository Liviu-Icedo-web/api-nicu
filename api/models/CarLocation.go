package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

//COmment
type CarLocation struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	CarID     uint32    `gorm:"int" json:"user_id"`
	Car       Car       `json:"car"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Country   string    `json:"country"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *CarLocation) Prepare() {
	u.ID = 0
	u.Street = html.EscapeString(strings.TrimSpace(u.Street))
	u.City = html.EscapeString(strings.TrimSpace(u.City))
	u.State = html.EscapeString(strings.TrimSpace(u.State))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (p *CarLocation) Validate() error {

	if p.Street == "" {
		return errors.New("Required Street")
	}
	if p.City == "" {
		return errors.New("Required City")
	}

	if p.State == "" {
		return errors.New("Required State")
	}

	return nil
}

func (p *CarLocation) SaveCarLocation(db *gorm.DB) (*CarLocation, error) {
	var err error
	err = db.Debug().Model(&Car{}).Create(&p).Error
	if err != nil {
		return &CarLocation{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Car{}).Where("id = ?", p.CarID).Take(&p.Car).Error
		if err != nil {
			return &CarLocation{}, err
		}
	}
	return p, nil
}

func (p *CarLocation) FindAllCarsLocation(db *gorm.DB) (*[]CarLocation, error) {
	var err error
	posts := []CarLocation{}
	err = db.Debug().Model(&CarLocation{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]CarLocation{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&CarLocation{}).Where("id = ?", posts[i].CarID).Take(&posts[i].Car).Error
			if err != nil {
				return &[]CarLocation{}, err
			}
		}
	}
	return &posts, nil
}

func (p *CarLocation) FindCarLocationByID(db *gorm.DB, pid uint64) (*CarLocation, error) {
	var err error
	err = db.Debug().Model(&Car{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &CarLocation{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Car{}).Where("id = ?", p.CarID).Take(&p.Car).Error
		if err != nil {
			return &CarLocation{}, err
		}
	}
	return p, nil
}

func (p *CarLocation) UpdateACarLocation(db *gorm.DB, pid uint64) (*CarLocation, error) {

	var err error
	db = db.Debug().Model(&CarLocation{}).Where("id = ?", pid).Take(&CarLocation{}).UpdateColumns(
		map[string]interface{}{
			"car_id":     p.CarID,
			"Street":     p.Street,
			"city":       p.City,
			"sate":       p.State,
			"country":    p.Country,
			"updated_at": time.Now(),
		},
	)
	err = db.Debug().Model(&CarLocation{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &CarLocation{}, err
	}
	/*if p.ID != 0 { LIVIU AREGLA ESTO, cuando haces un UPDATE no guarda bien el owner.Aun que  no lo necesitas
		err = db.Debug().Model(&User{}).Where("id = ?", p.User_id).Take(&p.Owner).Error
		if err != nil {
			return &Car{}, err
		}
	}*/
	return p, nil
}

func (p *CarLocation) DeleteACar(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&CarLocation{}).Where("id = ? and car_id = ?", pid, uid).Take(&CarLocation{}).Delete(&CarLocation{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Location not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
