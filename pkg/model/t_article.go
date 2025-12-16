package model

const TableNameTArticle = "t_article_proofread"

type TArticle struct {
	ArticleID int64  `json:"articleId" gorm:"column:article_id;primaryKey"` // 文章id
	Title     string `json:"title" gorm:"column:title"`                     // 文章标题
	Content   string `json:"content" gorm:"column:content"`                 //文章内容
}

func (*TArticle) TableName() string {
	return TableNameTArticle
}
