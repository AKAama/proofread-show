<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>文章校阅列表</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background-color: #f5f5f5;
            padding: 20px;
            line-height: 1.6;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            margin-bottom: 30px;
            font-size: 28px;
        }
        .article-list {
            list-style: none;
        }
        .article-item {
            padding: 15px;
            margin-bottom: 10px;
            background: #f9f9f9;
            border-left: 4px solid #4CAF50;
            border-radius: 4px;
            transition: all 0.3s;
        }
        .article-item:hover {
            background: #f0f0f0;
            transform: translateX(5px);
        }
        .article-item a {
            text-decoration: none;
            color: #333;
            font-size: 16px;
            display: block;
        }
        .article-item a:hover {
            color: #4CAF50;
        }
        .pagination {
            margin-top: 30px;
            text-align: center;
        }
        .pagination a, .pagination span {
            display: inline-block;
            padding: 8px 12px;
            margin: 0 4px;
            text-decoration: none;
            border: 1px solid #ddd;
            border-radius: 4px;
            color: #333;
        }
        .pagination a:hover {
            background: #4CAF50;
            color: white;
            border-color: #4CAF50;
        }
        .pagination .current {
            background: #4CAF50;
            color: white;
            border-color: #4CAF50;
        }
        .pagination .disabled {
            color: #ccc;
            cursor: not-allowed;
            pointer-events: none;
        }
        .pagination .ellipsis {
            padding: 8px 4px;
            color: #999;
            pointer-events: none;
        }
        .info {
            margin-bottom: 20px;
            color: #666;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>文章校阅列表</h1>
        <div class="info">
            共 {{.total}} 篇文章，第 {{.page}} / {{.totalPages}} 页
        </div>
        <ul class="article-list">
            {{range .articles}}
            <li class="article-item">
                <a href="/articles/{{.ArticleID}}">文章 ID: {{.ArticleID}}</a>
            </li>
            {{else}}
            <li style="text-align: center; padding: 40px; color: #999;">
                暂无文章
            </li>
            {{end}}
        </ul>
        
        {{if gt .totalPages 1}}
        <div class="pagination">
            {{if gt .page 1}}
            <a href="/articles?page={{sub .page 1}}&pageSize={{.pageSize}}">上一页</a>
            {{else}}
            <span class="disabled">上一页</span>
            {{end}}
            
            {{$pageSize := .pageSize}}
            {{$currentPage := .page}}
            {{range .pages}}
                {{if eq . "..."}}
                <span class="ellipsis">...</span>
                {{else if eq . $currentPage}}
                <span class="current">{{.}}</span>
                {{else}}
                <a href="/articles?page={{.}}&pageSize={{$pageSize}}">{{.}}</a>
                {{end}}
            {{end}}
            
            {{if lt .page .totalPages}}
            <a href="/articles?page={{add .page 1}}&pageSize={{.pageSize}}">下一页</a>
            {{else}}
            <span class="disabled">下一页</span>
            {{end}}
        </div>
        {{end}}
    </div>
</body>
</html>

