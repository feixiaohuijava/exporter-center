package webhookmanager

import (
	"gorm.io/gorm"
	"time"
)

type WebhookmanagerPhonerecord struct {
	ID          int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	Username    string    `gorm:"column:username;type:varchar(100)" json:"username"`
	Phone       string    `gorm:"column:phone;type:varchar(100)" json:"phone"`
	AlertName   string    `gorm:"column:alertName;type:varchar(100)" json:"alert_name"`
	Result      bool      `gorm:"column:result;type:tinyint(1);not null" json:"result"`
	ResultMsg   string    `gorm:"column:resultMsg;type:longtext" json:"result_msg"`
	CreatedTime time.Time `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
	Project     string    `gorm:"column:project;type:varchar(100)" json:"project"`
	VoiceName   string    `gorm:"column:voiceName;type:varchar(100)" json:"voice_name"`
}

func (phonerecord WebhookmanagerPhonerecord) TableName() string {
	return "webhookmanager_phonerecord"
}

func InsertPhoneRecord(db *gorm.DB, phoneRecord WebhookmanagerPhonerecord) {
	db.Create(&phoneRecord)
}
