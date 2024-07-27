package mysql

import (
	"GameApp/entity"
	"database/sql"
	"fmt"
	"log"
)

// this function will return true if user is unique and vice versa
func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	//user := entity.User{}
	//var created_at []uint8

	row := d.db.QueryRow(`select * from users where phone_number = ?`, phoneNumber)
	//err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &created_at)
	user, err := scanUser(row)

	//fmt.Printf("user phone number %s --- req phone number %s", user.PhoneNumber, phoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		log.Fatalln("err : %w", err)
		return false, err
		//return false, fmt.Errorf("Can't scan query result: %w", err)
	} else if user.PhoneNumber == phoneNumber {
		return false, fmt.Errorf("This user already exist !!!")
	} else {
		return true, nil
	}
}

func (d *MySQLDB) Register(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`insert into users(name, phone_number, password) values(?, ?, ?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("Can't execute cpmmand: %w", err)
	}

	// error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *MySQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error) {
	//user := entity.User{}
	//var created_at []uint8

	row := d.db.QueryRow(`select * from users where phone_number = ?`, phoneNumber)
	user, err := scanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, fmt.Errorf("Can't scan query result: %w", err)
	}
	fmt.Println(user.Name)
	return user, true, nil
}

func (d *MySQLDB) GetUserByID(userID uint) (entity.User, error) {
	//user := entity.User{}
	//var created_at []uint8

	row := d.db.QueryRow(`select * from users where id = ?`, userID)
	//err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &created_at)
	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("Record not found.")
		}
		return entity.User{}, fmt.Errorf("Can't scan query result: %w", err)
	}
	fmt.Println(user.Name)
	return user, nil
}
func scanUser(row *sql.Row) (entity.User, error) {
	var created_at []uint8
	var user entity.User

	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &created_at)

	return user, err
}
