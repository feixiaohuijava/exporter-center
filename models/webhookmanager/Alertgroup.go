package webhookmanager

import "time"

// WebhookmanagerAlertgroup [...]
type WebhookmanagerAlertgroup struct {
	ID                     int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	AlertSource            string    `gorm:"unique_index:webhookmanager_alertgroup_alertSource_groupName_9ef626b6_uniq;column:alertSource;type:varchar(100)" json:"alert_source"`
	GroupName              string    `gorm:"unique_index:webhookmanager_alertgroup_alertSource_groupName_9ef626b6_uniq;column:groupName;type:varchar(100)" json:"group_name"`
	GroupFile              string    `gorm:"column:groupFile;type:varchar(255)" json:"group_file"`
	GroupInterval          int       `gorm:"column:groupInterval;type:int(11)" json:"group_interval"`
	Env                    string    `gorm:"column:env;type:varchar(10)" json:"env"`
	CreatedTime            time.Time `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
	UpdateTime             time.Time `gorm:"column:updateTime;type:datetime(6)" json:"update_time"`
	PrometheusRuleOperator string    `gorm:"column:prometheusRuleOperator;type:varchar(100)" json:"prometheus_rule_operator"`
}
