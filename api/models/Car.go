package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Car struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	User_id   uint32    `gorm:"int" json:"user_id"`
	Owner     User      `json:"user"`
	Brand     string    `gorm:"size:255;not null;unique" json:"brand"`
	Year      int       `gorm:"size:255;not null;unique" json:"year"`
	Hp        int       `gorm:"int" json:"hp"`
	Doors     int       `gorm:"int" json:"doors"`
	Seats     int       `gorm:"int" json:"seats"`
	Insurance string    `gorm:"size:255;not null;unique" json:"insurance"`
	Images    string    `gorm:"json" json:"images"`
	Town      string    `gorm: json:"town"`
	PriceDay  float64   `gorm: json:"price_day"`
	PriceHour float64   `gorm: json:"price_hour"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *Car) Prepare() {
	u.ID = 0
	u.Brand = html.EscapeString(strings.TrimSpace(u.Brand))
	u.Insurance = html.EscapeString(strings.TrimSpace(u.Insurance))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (p *Car) Validate() error {

	if p.Brand == "" {
		return errors.New("Required Brand")
	}
	if p.Year < 1 {
		return errors.New("Required Year")
	}
	if p.Seats < 1 {
		return errors.New("Required Seats")
	}
	return nil
}

func (p *Car) SaveCar(db *gorm.DB) (*Car, error) {
	var err error
	err = db.Debug().Model(&Car{}).Create(&p).Error
	if err != nil {
		return &Car{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.User_id).Take(&p.Owner).Error
		if err != nil {
			return &Car{}, err
		}
	}
	return p, nil
}

func (p *Car) FindAllCars(db *gorm.DB) (*[]Car, error) {
	var err error
	posts := []Car{}
	err = db.Debug().Model(&Car{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Car{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].User_id).Take(&posts[i].Owner).Error
			if err != nil {
				return &[]Car{}, err
			}
		}
	}
	return &posts, nil
}

func (p *Car) FindCarByID(db *gorm.DB, pid uint64) (*Car, error) {
	var err error
	err = db.Debug().Model(&Car{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Car{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.User_id).Take(&p.Owner).Error
		if err != nil {
			return &Car{}, err
		}
	}
	return p, nil
}

func (p *Car) UpdateACar(db *gorm.DB, pid uint64) (*Car, error) {

	var err error
	db = db.Debug().Model(&Car{}).Where("id = ?", pid).Take(&Car{}).UpdateColumns(
		map[string]interface{}{
			"user_id":    p.User_id,
			"brand":      p.Brand,
			"year":       p.Year,
			"doors":      p.Doors,
			"hp":         p.Hp,
			"seats":      p.Seats,
			"images":     p.Images,
			"insurance":  p.Insurance,
			"updated_at": time.Now(),
		},
	)
	err = db.Debug().Model(&Car{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Car{}, err
	}
	fmt.Println("USerr:", p.User_id)
	fmt.Println("owner:", p.Owner)
	/*if p.ID != 0 { LIVIU AREGLA ESTO, cuando haces un UPDATE no guarda bien el owner.Aun que  no lo necesitas
		err = db.Debug().Model(&User{}).Where("id = ?", p.User_id).Take(&p.Owner).Error
		if err != nil {
			return &Car{}, err
		}
	}*/
	return p, nil
}

func (p *Car) DeleteACar(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Car{}).Where("id = ? and user_id = ?", pid, uid).Take(&Car{}).Delete(&Car{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Car not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
