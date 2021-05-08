package webhookmanager

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

// WebhookmanagerAlertrule [...]
type WebhookmanagerAlertrule struct {
	ID                          int                         `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	RuleName                    string                      `gorm:"unique_index:webhookmanager_alertrule_ruleName_alertGroup_id_8bd06451_uniq;column:ruleName;type:varchar(100)" json:"rule_name"`
	RuleQuery                   string                      `gorm:"column:ruleQuery;type:longtext" json:"rule_query"`
	RuleDuration                string                      `gorm:"column:ruleDuration;type:varchar(100)" json:"rule_duration"`
	RuleLabels                  string                      `gorm:"column:ruleLabels;type:json" json:"rule_labels"`
	RuleAnnotations             string                      `gorm:"column:ruleAnnotations;type:json" json:"rule_annotations"`
	RuleAlerts                  string                      `gorm:"column:ruleAlerts;type:json" json:"rule_alerts"`
	RuleHealth                  string                      `gorm:"column:ruleHealth;type:varchar(30)" json:"rule_health"`
	RuleLevel                   string                      `gorm:"column:ruleLevel;type:varchar(30)" json:"rule_level"`
	AlertGroupID                int                         `gorm:"unique_index:webhookmanager_alertrule_ruleName_alertGroup_id_8bd06451_uniq;index:webhookmanager_alert_alertGroup_id_de323c0a_fk_webhookma;column:alertGroup_id;type:int(11)" json:"alert_group_id"`
	WebhookmanagerAlertgroup    WebhookmanagerAlertgroup    `gorm:"association_foreignkey:alertGroup_id;foreignkey:id" json:"webhookmanager_alertgroup_list"`
	CreatedTime                 time.Time                   `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
	UpdateTime                  time.Time                   `gorm:"column:updateTime;type:datetime(6)" json:"update_time"`
	CallOpsFlag                 bool                        `gorm:"column:CallOpsFlag;type:tinyint(1);not null" json:"call_ops_flag"`
	Type                        string                      `gorm:"column:type;type:varchar(30)" json:"type"`
	VoiceName                   string                      `gorm:"column:voiceName;type:varchar(10)" json:"voiceName"`
	NoticeFlag                  bool                        `gorm:"column:noticeFlag;type:tinyint(1);not null" json:"notice_flag"`
	Subsystems                  string                      `gorm:"column:subsystems;type:json" json:"subsystems"`
	MutipleNoticeFlag           bool                        `gorm:"column:mutipleNoticeFlag;type:tinyint(1);not null" json:"mutiple_notice_flag"`
	GroupLabelKeyID             int                         `gorm:"index:webhookmanager_alert_groupLabelKey_id_fe374438_fk_webhookma;column:groupLabelKey_id;type:int(11)" json:"group_label_key_id"`
	WebhookmanagerGrouplabelkey WebhookmanagerGrouplabelkey `gorm:"association_foreignkey:groupLabelKey_id;foreignkey:id" json:"webhookmanager_grouplabelkey_list"`
	NotiteKeyValue              string                      `gorm:"column:notite_key_value;type:json" json:"notite_key_value"`
	DistributedKey              string                      `gorm:"column:distributedKey;type:varchar(100)" json:"distributed_key"`
	CallDevFlag                 bool                        `gorm:"column:callDevFlag;type:tinyint(1);not null" json:"call_dev_flag"`
	CallTestFlag                bool                        `gorm:"column:callTestFlag;type:tinyint(1);not null" json:"call_test_flag"`
	CallProductFlag             bool                        `gorm:"column:callProductFlag;type:tinyint(1);not null" json:"call_product_flag"`
}

func (webhookmanagerAlertrule WebhookmanagerAlertrule) TableName() string {
	return "webhookmanager_alertrule"
}

// 根据id返回监控规则这个对象
func FindAlertById(db *gorm.DB, alertRuleId int, logger *logrus.Logger) WebhookmanagerAlertrule {
	var alertRule WebhookmanagerAlertrule
	db.Where(&WebhookmanagerAlertrule{ID: alertRuleId}).First(&alertRule)
	if (alertRule == WebhookmanagerAlertrule{}) {
		logger.Error("根据id查不到监控规则")
		panic("根据id查不到监控规则")
	} else {
		return alertRule
	}
}
