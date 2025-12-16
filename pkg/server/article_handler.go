package server

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"proofread-show/pkg/db"
	"proofread-show/pkg/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ArticleListResponse 文章列表响应
type ArticleListResponse struct {
	Articles   []ArticleItem `json:"articles"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
	TotalPages int           `json:"totalPages"`
}

// ArticleItem 文章列表项
type ArticleItem struct {
	ArticleID int64 `json:"articleId"`
}

// ArticleDetailResponse 文章详情响应
type ArticleDetailResponse struct {
	ArticleID          int64                    `json:"articleId"`
	Content            string                   `json:"content"`
	Results            []model.TProofreadResult `json:"results"`
	HighlightedContent string                   `json:"highlightedContent"`
}

// ArticleWithHighlight 带高亮的文章
type ArticleWithHighlight struct {
	ArticleID          int64  `json:"articleId"`
	Title              string `json:"title"`
	HighlightedContent string `json:"highlightedContent"`
}

// GetAllArticles 获取所有文章（平铺展示）
func GetAllArticles(c *gin.Context) {
	// 查询有校阅数据的文章ID（去重）
	var articleIDs []int64
	if err := db.GetDB().
		Table(model.TableNameTProofreadResult).
		Select("DISTINCT article_id").
		Where("article_id IS NOT NULL").
		Pluck("article_id", &articleIDs).Error; err != nil {
		zap.S().Errorf("查询文章ID列表失败: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	if len(articleIDs) == 0 {
		c.HTML(http.StatusOK, "articles.tpl", gin.H{
			"articles": []ArticleWithHighlight{},
		})
		return
	}

	// 查询所有文章及其校阅结果
	var articlesWithHighlight []ArticleWithHighlight

	for _, articleID := range articleIDs {
		// 查询文章内容
		var article model.TArticle
		if err := db.GetDB().
			Where("article_id = ?", articleID).
			First(&article).Error; err != nil {
			zap.S().Warnf("查询文章 %d 失败: %s", articleID, err.Error())
			continue
		}

		// 查询校阅结果
		var results []model.TProofreadResult
		if err := db.GetDB().
			Where("article_id = ?", articleID).
			Order("start ASC").
			Find(&results).Error; err != nil {
			zap.S().Warnf("查询文章 %d 的校阅结果失败: %s", articleID, err.Error())
			continue
		}

		// 去除原文中的 HTML 标签
		plainContent := stripHTMLTags(article.Content)

		// 生成高亮后的内容
		highlightedContent := highlightContent(plainContent, results)

		articlesWithHighlight = append(articlesWithHighlight, ArticleWithHighlight{
			ArticleID:          articleID,
			Title:              article.Title,
			HighlightedContent: highlightedContent,
		})
	}

	c.HTML(http.StatusOK, "articles.tpl", gin.H{
		"articles": articlesWithHighlight,
	})
}

// highlightContent 将原文内容按照校阅结果进行高亮标记
func highlightContent(content string, results []model.TProofreadResult) string {
	if len(results) == 0 {
		return html.EscapeString(content)
	}

	// 按位置从后往前排序，避免插入时位置偏移
	sortedResults := make([]model.TProofreadResult, len(results))
	copy(sortedResults, results)
	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].Start > sortedResults[j].Start
	})

	// 将内容转换为 rune 切片以正确处理中文
	contentRunes := []rune(content)

	// 从后往前插入高亮标记
	for _, result := range sortedResults {
		start := result.Start
		end := result.End

		// 确保位置有效
		if start < 0 || end > len(contentRunes) || start >= end {
			continue
		}

		// 构建 tooltip 内容
		tooltipHTML := buildTooltip(result)

		// 插入结束标记
		endTag := "</span>"
		contentRunes = insertRunes(contentRunes, end, []rune(endTag))

		// 插入开始标记和 tooltip
		startTag := fmt.Sprintf(`<span class="highlight"><span class="tooltip">%s</span>`, tooltipHTML)
		contentRunes = insertRunes(contentRunes, start, []rune(startTag))
	}

	// 先转义整个内容（包括我们插入的标签）
	escaped := html.EscapeString(string(contentRunes))

	// 恢复我们插入的 HTML 标签（需要按照嵌套顺序恢复）
	// 先恢复最外层的 highlight span
	escaped = strings.ReplaceAll(escaped, "&lt;span class=&#34;highlight&#34;&gt;", "<span class=\"highlight\">")
	// 恢复 tooltip span
	escaped = strings.ReplaceAll(escaped, "&lt;span class=&#34;tooltip&#34;&gt;", "<span class=\"tooltip\">")
	// 恢复 tooltip-content span
	escaped = strings.ReplaceAll(escaped, "&lt;span class=&#34;tooltip-content&#34;&gt;", "<span class=\"tooltip-content\">")
	escaped = strings.ReplaceAll(escaped, "&lt;span class=&#34;tooltip-content tooltip-message&#34;&gt;", "<span class=\"tooltip-content tooltip-message\">")
	// 恢复 tooltip-message span
	escaped = strings.ReplaceAll(escaped, "&lt;span class=&#34;tooltip-message&#34;&gt;", "<span class=\"tooltip-message\">")
	// 恢复 tooltip-suggestion span
	escaped = strings.ReplaceAll(escaped, "&lt;span class=&#34;tooltip-suggestion&#34;&gt;", "<span class=\"tooltip-suggestion\">")
	// 恢复所有结束标签
	escaped = strings.ReplaceAll(escaped, "&lt;/span&gt;", "</span>")

	return escaped
}

// buildTooltip 构建 tooltip HTML 内容
func buildTooltip(result model.TProofreadResult) string {
	var parts []string

	// 处理建议
	if result.Suggestion != "" {
		var suggestions []string
		if err := json.Unmarshal([]byte(result.Suggestion), &suggestions); err != nil {
			// 如果不是 JSON 数组，直接使用原值
			suggestions = []string{result.Suggestion}
		}
		if len(suggestions) > 0 {
			suggestionText := html.EscapeString(suggestions[0])
			if len(suggestions) > 1 {
				for i := 1; i < len(suggestions); i++ {
					suggestionText += ", " + html.EscapeString(suggestions[i])
				}
			}
			parts = append(parts, fmt.Sprintf(`<span class="tooltip-content">建议: <span class="tooltip-suggestion">%s</span></span>`, suggestionText))
		}
	}

	// 处理消息
	if result.Message != "" {
		messageText := html.EscapeString(result.Message)
		parts = append(parts, fmt.Sprintf(`<span class="tooltip-content">原因: <span class="tooltip-message">%s</span></span>`, messageText))
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "")
}

// stripHTMLTags 去除 HTML 标签，只保留文本内容
func stripHTMLTags(htmlContent string) string {
	// 使用正则表达式去除所有 HTML 标签
	re := regexp.MustCompile(`<[^>]*>`)
	plainText := re.ReplaceAllString(htmlContent, "")

	// 解码 HTML 实体（如 &nbsp; &lt; 等）
	plainText = html.UnescapeString(plainText)

	return plainText
}

// insertRunes 在指定位置插入 rune 切片
func insertRunes(slice []rune, index int, insert []rune) []rune {
	if index < 0 || index > len(slice) {
		return slice
	}
	result := make([]rune, 0, len(slice)+len(insert))
	result = append(result, slice[:index]...)
	result = append(result, insert...)
	result = append(result, slice[index:]...)
	return result
}
