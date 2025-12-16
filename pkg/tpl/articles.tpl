<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel=”shortcut icon” href=”favicon.ico” type=”image/x-icon” />
    <title>文章校阅展示</title>
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
            line-height: 1.8;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1 {
            color: #333;
            margin-bottom: 30px;
            font-size: 28px;
            text-align: center;
        }
        .article-item {
            background: white;
            padding: 30px;
            margin-bottom: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .article-header {
            margin-bottom: 20px;
            padding-bottom: 15px;
            border-bottom: 2px solid #eee;
        }
        .article-header h2 {
            color: #333;
            font-size: 20px;
            margin-bottom: 5px;
        }
        .article-id {
            color: #666;
            font-size: 14px;
        }
        .article-content {
            font-size: 16px;
            color: #333;
            white-space: pre-wrap;
            word-wrap: break-word;
        }
        .highlight {
            background: linear-gradient(120deg, #ffd54f 0%, #ffb74d 100%);
            padding: 2px 4px;
            border-radius: 3px;
            cursor: pointer;
            position: relative;
            display: inline-block;
            transition: all 0.2s;
        }
        .highlight:hover {
            background: linear-gradient(120deg, #ffc107 0%, #ff9800 100%);
            box-shadow: 0 2px 8px rgba(255, 193, 7, 0.4);
        }
        .tooltip {
            position: absolute;
            bottom: 100%;
            left: 50%;
            transform: translateX(-50%);
            margin-bottom: 8px;
            padding: 12px 16px;
            background: #333;
            color: white;
            border-radius: 6px;
            font-size: 14px;
            white-space: nowrap;
            opacity: 0;
            pointer-events: none;
            transition: opacity 0.3s;
            z-index: 1000;
            box-shadow: 0 4px 12px rgba(0,0,0,0.3);
        }
        .tooltip::after {
            content: '';
            position: absolute;
            top: 100%;
            left: 50%;
            transform: translateX(-50%);
            border: 6px solid transparent;
            border-top-color: #333;
        }
        .highlight:hover .tooltip {
            opacity: 1;
        }
        .tooltip-content {
            display: block;
            margin-bottom: 6px;
        }
        .tooltip-suggestion {
            color: #4CAF50;
            font-weight: bold;
        }
        .tooltip-message {
            color: #ffd54f;
            font-size: 13px;
        }
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #999;
            font-size: 16px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>文章校阅展示</h1>
        {{if .articles}}
            {{range .articles}}
            <div class="article-item">
                <div class="article-header">
                    {{if .Title}}
                    <h2>{{.Title}}</h2>
                    {{else}}
                    <h2>文章 ID: {{.ArticleID}}</h2>
                    {{end}}
                    <div class="article-id">Article ID: {{.ArticleID}}</div>
                </div>
                <div class="article-content">{{.HighlightedContent | safeHTML}}</div>
            </div>
            {{end}}
        {{else}}
            <div class="empty-state">暂无文章</div>
        {{end}}
    </div>
</body>
</html>

