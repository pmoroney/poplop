package db

import (
	"database/sql"
	"log"
	"pmoroney/poplop"
)

/*
CREATE TABLE `Scheme` (
	`Name` varchar(255) NOT NULL DEFAULT '',
	`Counter` int(11) NOT NULL DEFAULT '0',
	`Username` varchar(255) NOT NULL DEFAULT '',
	`URL` varchar(255) NOT NULL DEFAULT '',
	`Notes` varchar(255) NOT NULL DEFAULT '',
	`Forbidden` varchar(255) CHARACTER SET ascii NOT NULL DEFAULT '',
	`MaxLength` int(11) NOT NULL DEFAULT '0',
	`Legacy` tinyint(1) NOT NULL DEFAULT '0',
	`CreateDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`UpdateDate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`Name`),
	KEY `URL` (`URL`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
*/

func GetScheme(name string) (poplop.Scheme, error) {
	var s poplop.Scheme
	err := db.QueryRow("SELECT Name, Counter, Username, URL, Notes, Forbidden, MaxLength, Legacy from Scheme where Name = ?", name).Scan(&s.Name, &s.Counter, &s.Username, &s.URL, &s.Notes, &s.Forbidden, &s.MaxLength, &s.Legacy)
	return s, err
}

func InsertScheme(s poplop.Scheme) error {
	_, err := db.Exec("INSERT INTO Scheme (Name, Counter, Username, URL, Notes, Forbidden, MaxLength, Legacy) values (?, ?, ?, ?, ?, ?, ?, ?)", s.Name, s.Counter, s.Username, s.URL, s.Notes, s.Forbidden, s.MaxLength, s.Legacy)
	return err
}

func UpdateScheme(s poplop.Scheme, oldName string) error {
	if oldName == "" {
		oldName = s.Name
	}
	_, err := db.Exec("UPDATE Scheme SET Name = ?, Counter = ?, Username = ?, URL = ?, Notes = ?, Forbidden = ?, MaxLength = ?, Legacy = ? WHERE Name = ?", s.Name, s.Counter, s.Username, s.URL, s.Notes, s.Forbidden, s.MaxLength, s.Legacy, oldName)
	return err
}

var db *sql.DB

func Connect() {
	var err error
	db, err = sql.Open("mysql", "username:password@unix(/var/run/mysqld/mysqld.sock)/poplop")
	if err != nil {
		log.Fatal(err)
	}
}
