package server

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"sort"
	"strconv"
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
	ArticleID int64                    `json:"articleId"`
	Content   string                   `json:"content"`
	Results   []model.TProofreadResult `json:"results"`
}

// GetArticleList 获取文章列表（分页）
func GetArticleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var articles []ArticleItem
	var total int64

	// 查询有校阅数据的文章ID（去重）
	query := db.GetDB().
		Table(model.TableNameTProofreadResult).
		Select("DISTINCT article_id").
		Where("article_id IS NOT NULL")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		zap.S().Errorf("查询文章总数失败: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 分页查询
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Scan(&articles).Error; err != nil {
		zap.S().Errorf("查询文章列表失败: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response := ArticleListResponse{
		Articles:   articles,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	// 生成分页页码数组（带省略号）
	pages := generatePaginationPages(response.Page, response.TotalPages)

	c.HTML(http.StatusOK, "article_list.tpl", gin.H{
		"articles":   response.Articles,
		"total":      response.Total,
		"page":       response.Page,
		"pageSize":   response.PageSize,
		"totalPages": response.TotalPages,
		"pages":      pages,
	})
}

// GetArticleDetail 获取文章详情（包含校阅结果）
func GetArticleDetail(c *gin.Context) {
	articleIDStr := c.Param("id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 查询文章内容
	var article model.TArticle
	if err := db.GetDB().
		Where("article_id = ?", articleID).
		First(&article).Error; err != nil {
		zap.S().Errorf("查询文章失败: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 查询校阅结果
	var results []model.TProofreadResult
	if err := db.GetDB().
		Where("article_id = ?", articleID).
		Order("start ASC").
		Find(&results).Error; err != nil {
		zap.S().Errorf("查询校阅结果失败: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 去除原文中的 HTML 标签
	plainContent := stripHTMLTags(article.Content)

	// 生成高亮后的内容
	highlightedContent := highlightContent(plainContent, results)

	response := ArticleDetailResponse{
		ArticleID: article.ArticleID,
		Content:   article.Content,
		Results:   results,
	}

	c.HTML(http.StatusOK, "article_detail.tpl", gin.H{
		"article":            response,
		"highlightedContent": highlightedContent,
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
		parts = append(parts, fmt.Sprintf(`<span class="tooltip-content tooltip-message">%s</span>`, messageText))
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "")
}

// generatePaginationPages 生成分页页码数组，中间用省略号
func generatePaginationPages(currentPage, totalPages int) []interface{} {
	const maxVisiblePages = 7 // 最多显示7个页码（包括省略号）

	if totalPages <= maxVisiblePages {
		// 如果总页数不多，直接返回所有页码
		pages := make([]interface{}, totalPages)
		for i := 1; i <= totalPages; i++ {
			pages[i-1] = i
		}
		return pages
	}

	var pages []interface{}

	// 总是显示第一页
	pages = append(pages, 1)

	// 计算开始和结束页码
	startPage := currentPage - 2
	endPage := currentPage + 2

	if startPage < 2 {
		startPage = 2
		endPage = startPage + 4
		if endPage > totalPages-1 {
			endPage = totalPages - 1
			startPage = endPage - 4
			if startPage < 2 {
				startPage = 2
			}
		}
	} else if endPage > totalPages-1 {
		endPage = totalPages - 1
		startPage = endPage - 4
		if startPage < 2 {
			startPage = 2
		}
	}

	// 如果开始页码不是2，添加省略号
	if startPage > 2 {
		pages = append(pages, "...")
	}

	// 添加中间页码
	for i := startPage; i <= endPage; i++ {
		pages = append(pages, i)
	}

	// 如果结束页码不是倒数第二页，添加省略号
	if endPage < totalPages-1 {
		pages = append(pages, "...")
	}

	pages = append(pages, totalPages)

	return pages
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
