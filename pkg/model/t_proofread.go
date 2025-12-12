package model

const TableNameTProofreadResult = "t_proofread_result"

// TProofreadResult 校阅结果表
type TProofreadResult struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`      // 主键ID
	ArticleID  int64  `json:"articleId" gorm:"column:article_id;index"`          // 文章ID
	Type       string `json:"type" gorm:"column:type;type:varchar(50)"`          // 错误类型 (grammar/style等)
	Text       string `json:"text" gorm:"column:text;type:text"`                 // 错误文本
	Start      int    `json:"start" gorm:"column:start"`                         // 开始位置
	End        int    `json:"end" gorm:"column:end"`                             // 结束位置
	Suggestion string `json:"suggestion" gorm:"column:suggestion;type:text"`     // 建议（JSON数组字符串）
	Message    string `json:"message" gorm:"column:message;type:varchar(500)"`   // 错误消息
	Sentence   string `json:"sentence" gorm:"column:sentence;type:text"`         // 所在句子
	CreatedAt  int64  `json:"createdAt" gorm:"column:created_at;autoCreateTime"` // 创建时间
}

func (TProofreadResult) TableName() string {
	return TableNameTProofreadResult
}
