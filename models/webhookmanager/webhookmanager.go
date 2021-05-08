package webhookmanager

import (
	"time"
)

// WebhookmanagerAlertorder [...]
type WebhookmanagerAlertorder struct {
	ID              int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	Kind            string    `gorm:"column:kind;type:varchar(50)" json:"kind"`
	Alername        string    `gorm:"column:alername;type:varchar(100)" json:"alername"`
	Status          string    `gorm:"column:status;type:varchar(50)" json:"status"`
	Principal       string    `gorm:"column:principal;type:varchar(100)" json:"principal"`
	Comment         string    `gorm:"column:comment;type:longtext" json:"comment"`
	CreatedTime     time.Time `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
	UpdateTime      time.Time `gorm:"column:updateTime;type:datetime(6)" json:"update_time"`
	PrometheusURL   string    `gorm:"column:prometheusUrl;type:varchar(100)" json:"prometheus_url"`
	AlertDuty       string    `gorm:"column:alertDuty;type:varchar(50)" json:"alert_duty"`
	AlertLabel      string    `gorm:"column:alertLabel;type:longtext" json:"alert_label"`
	AlertReason     string    `gorm:"column:alertReason;type:longtext" json:"alert_reason"`
	AlertResolution string    `gorm:"column:alertResolution;type:json;not null" json:"alert_resolution"`
	FollowingTime   time.Time `gorm:"column:followingTime;type:datetime(6)" json:"following_time"`
	EndTime         time.Time `gorm:"column:endTime;type:datetime(6)" json:"end_time"`
}

// WebhookmanagerAlertruleAutoHealPackage [...]
type WebhookmanagerAlertruleAutoHealPackage struct {
	ID                            int                           `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	AlertruleID                   int                           `gorm:"unique_index:webhookmanager_alertrule_alertrule_id_autohealpac_0884264f_uniq;column:alertrule_id;type:int(11);not null" json:"alertrule_id"`
	WebhookmanagerAlertrule       WebhookmanagerAlertrule       `gorm:"association_foreignkey:alertrule_id;foreignkey:id" json:"webhookmanager_alertrule_list"`
	AutohealpackageID             int                           `gorm:"unique_index:webhookmanager_alertrule_alertrule_id_autohealpac_0884264f_uniq;index:webhookmanager_alert_autohealpackage_id_416a3ed7_fk_webhookma;column:autohealpackage_id;type:int(11);not null" json:"autohealpackage_id"`
	WebhookmanagerAutohealpackage WebhookmanagerAutohealpackage `gorm:"association_foreignkey:autohealpackage_id;foreignkey:id" json:"webhookmanager_autohealpackage_list"`
}

// WebhookmanagerAlertruleRobotPrimary [...]
type WebhookmanagerAlertruleRobotPrimary struct {
	ID                      int                     `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	AlertruleID             int                     `gorm:"unique_index:webhookmanager_alertrule_alertrule_id_robot_id_567349bb_uniq;column:alertrule_id;type:int(11);not null" json:"alertrule_id"`
	WebhookmanagerAlertrule WebhookmanagerAlertrule `gorm:"association_foreignkey:alertrule_id;foreignkey:id" json:"webhookmanager_alertrule_list"`
	RobotID                 int                     `gorm:"unique_index:webhookmanager_alertrule_alertrule_id_robot_id_567349bb_uniq;index:webhookmanager_alert_robot_id_1aa26d23_fk_webhookma;column:robot_id;type:int(11);not null" json:"robot_id"`
	WebhookmanagerRobot     WebhookmanagerRobot     `gorm:"association_foreignkey:robot_id;foreignkey:id" json:"webhookmanager_robot_list"`
}

// WebhookmanagerAlertruleRobotSecond [...]
type WebhookmanagerAlertruleRobotSecond struct {
	ID                      int                     `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	AlertruleID             int                     `gorm:"unique_index:webhookmanager_alertrule_alertrule_id_robot_id_b2e8441d_uniq;column:alertrule_id;type:int(11);not null" json:"alertrule_id"`
	WebhookmanagerAlertrule WebhookmanagerAlertrule `gorm:"association_foreignkey:alertrule_id;foreignkey:id" json:"webhookmanager_alertrule_list"`
	RobotID                 int                     `gorm:"unique_index:webhookmanager_alertrule_alertrule_id_robot_id_b2e8441d_uniq;index:webhookmanager_alert_robot_id_676a503c_fk_webhookma;column:robot_id;type:int(11);not null" json:"robot_id"`
	WebhookmanagerRobot     WebhookmanagerRobot     `gorm:"association_foreignkey:robot_id;foreignkey:id" json:"webhookmanager_robot_list"`
}

// WebhookmanagerAlialertrule [...]
type WebhookmanagerAlialertrule struct {
	ID                  int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	RuleName            string    `gorm:"column:RuleName;type:varchar(100)" json:"rule_name"`
	AlertState          string    `gorm:"column:AlertState;type:varchar(100)" json:"alert_state"`
	ContactGroups       string    `gorm:"column:ContactGroups;type:varchar(100)" json:"contact_groups"`
	Dimensions          string    `gorm:"column:Dimensions;type:json" json:"dimensions"`
	EffectiveInterval   string    `gorm:"column:EffectiveInterval;type:varchar(100)" json:"effective_interval"`
	EnableState         bool      `gorm:"column:EnableState;type:tinyint(1);not null" json:"enable_state"`
	Escalations         string    `gorm:"column:Escalations;type:json" json:"escalations"`
	GroupID             string    `gorm:"column:GroupId;type:varchar(100)" json:"group_id"`
	GroupName           string    `gorm:"column:GroupName;type:varchar(100)" json:"group_name"`
	MailSubject         string    `gorm:"column:MailSubject;type:varchar(100)" json:"mail_subject"`
	MetricName          string    `gorm:"column:MetricName;type:varchar(100)" json:"metric_name"`
	Namespace           string    `gorm:"column:Namespace;type:varchar(100)" json:"namespace"`
	NoEffectiveInterval string    `gorm:"column:NoEffectiveInterval;type:varchar(100)" json:"no_effective_interval"`
	Period              int       `gorm:"column:Period;type:int(11)" json:"period"`
	Resources           string    `gorm:"column:Resources;type:json" json:"resources"`
	RuleID              string    `gorm:"column:RuleId;type:varchar(100)" json:"rule_id"`
	SilenceTime         int       `gorm:"column:SilenceTime;type:int(11)" json:"silence_time"`
	SourceType          string    `gorm:"column:SourceType;type:varchar(100)" json:"source_type"`
	Webhook             string    `gorm:"column:Webhook;type:varchar(100)" json:"webhook"`
	CreatedTime         time.Time `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
	UpdateTime          time.Time `gorm:"column:updateTime;type:datetime(6)" json:"update_time"`
	CallFlag            bool      `gorm:"column:callFlag;type:tinyint(1);not null" json:"call_flag"`
}

// WebhookmanagerAlialertruleRobotPrimary [...]
type WebhookmanagerAlialertruleRobotPrimary struct {
	ID                         int                        `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	AlialertruleID             int                        `gorm:"unique_index:webhookmanager_alialertr_alialertrule_id_robot_id_aebe0032_uniq;column:alialertrule_id;type:int(11);not null" json:"alialertrule_id"`
	WebhookmanagerAlialertrule WebhookmanagerAlialertrule `gorm:"association_foreignkey:alialertrule_id;foreignkey:id" json:"webhookmanager_alialertrule_list"`
	RobotID                    int                        `gorm:"unique_index:webhookmanager_alialertr_alialertrule_id_robot_id_aebe0032_uniq;index:webhookmanager_alial_robot_id_0fa87f4c_fk_webhookma;column:robot_id;type:int(11);not null" json:"robot_id"`
	WebhookmanagerRobot        WebhookmanagerRobot        `gorm:"association_foreignkey:robot_id;foreignkey:id" json:"webhookmanager_robot_list"`
}

// WebhookmanagerAlialertruleRobotSecond [...]
type WebhookmanagerAlialertruleRobotSecond struct {
	ID                         int                        `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	AlialertruleID             int                        `gorm:"unique_index:webhookmanager_alialertr_alialertrule_id_robot_id_2a8f1276_uniq;column:alialertrule_id;type:int(11);not null" json:"alialertrule_id"`
	WebhookmanagerAlialertrule WebhookmanagerAlialertrule `gorm:"association_foreignkey:alialertrule_id;foreignkey:id" json:"webhookmanager_alialertrule_list"`
	RobotID                    int                        `gorm:"unique_index:webhookmanager_alialertr_alialertrule_id_robot_id_2a8f1276_uniq;index:webhookmanager_alial_robot_id_197b8c71_fk_webhookma;column:robot_id;type:int(11);not null" json:"robot_id"`
	WebhookmanagerRobot        WebhookmanagerRobot        `gorm:"association_foreignkey:robot_id;foreignkey:id" json:"webhookmanager_robot_list"`
}

// WebhookmanagerAutohealpackage [...]
type WebhookmanagerAutohealpackage struct {
	ID                int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	HealName          string    `gorm:"column:healName;type:varchar(50)" json:"heal_name"`
	HealOperation     string    `gorm:"column:healOperation;type:longtext" json:"heal_operation"`
	CreatedTime       time.Time `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
	UpdateTime        time.Time `gorm:"column:updateTime;type:datetime(6)" json:"update_time"`
	Lock              bool      `gorm:"column:lock;type:tinyint(1);not null" json:"lock"`
	DurationEndTime   string    `gorm:"column:durationEndTime;type:varchar(20)" json:"duration_end_time"`
	DurationStartTime string    `gorm:"column:durationStartTime;type:varchar(20)" json:"duration_start_time"`
}

// WebhookmanagerGrouplabelkey [...]
type WebhookmanagerGrouplabelkey struct {
	ID          int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	Key         string    `gorm:"column:key;type:varchar(50)" json:"key"`
	KeyCallback string    `gorm:"column:keyCallback;type:varchar(100)" json:"key_callback"`
	CreatedTime time.Time `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
	UpdateTime  time.Time `gorm:"column:updateTime;type:datetime(6)" json:"update_time"`
}

// WebhookmanagerRabbitmqpublish [...]
type WebhookmanagerRabbitmqpublish struct {
	ID           int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	Namespace    string    `gorm:"column:namespace;type:varchar(100)" json:"namespace"`
	SrcQueueName string    `gorm:"column:src_queue_name;type:varchar(200)" json:"src_queue_name"`
	DstQueueName string    `gorm:"column:dst_queue_name;type:varchar(200)" json:"dst_queue_name"`
	PublishFlag  string    `gorm:"column:publishFlag;type:varchar(10)" json:"publish_flag"`
	CreatedTime  time.Time `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
}

// WebhookmanagerRobot [...]
type WebhookmanagerRobot struct {
	ID               int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	RobotType        string    `gorm:"column:robot_type;type:varchar(30)" json:"robot_type"`
	RobotDescription string    `gorm:"column:robot_description;type:longtext" json:"robot_description"`
	RobotName        string    `gorm:"unique;column:robot_name;type:varchar(100)" json:"robot_name"`
	RobotURL         string    `gorm:"column:robot_url;type:longtext" json:"robot_url"`
	RobotURLLabel    string    `gorm:"column:robot_url_label;type:json" json:"robot_url_label"`
	CreatedTime      time.Time `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
	UpdateTime       time.Time `gorm:"column:updateTime;type:datetime(6)" json:"update_time"`
	Channel          string    `gorm:"column:channel;type:varchar(100)" json:"channel"`
	Project          string    `gorm:"column:project;type:varchar(100)" json:"project"`
}

// WebhookmanagerSilencehistory [...]
type WebhookmanagerSilencehistory struct {
	ID          int       `gorm:"primary_key;column:id;type:int(11);not null" json:"-"`
	AlertRuleID int       `gorm:"column:alertRuleId;type:int(11);not null" json:"alert_rule_id"`
	Comment     string    `gorm:"column:comment;type:longtext" json:"comment"`
	CreatedBy   string    `gorm:"column:createdBy;type:varchar(100)" json:"created_by"`
	Hours       int       `gorm:"column:hours;type:int(11)" json:"hours"`
	Minutes     int       `gorm:"column:minutes;type:int(11)" json:"minutes"`
	Matchers    string    `gorm:"column:matchers;type:json" json:"matchers"`
	Flag        string    `gorm:"column:flag;type:longtext" json:"flag"`
	CreatedTime time.Time `gorm:"column:createdTime;type:datetime(6)" json:"created_time"`
}
