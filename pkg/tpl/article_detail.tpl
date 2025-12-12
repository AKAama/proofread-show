<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>文章校阅详情 - {{.article.ArticleID}}</title>
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
            max-width: 1000px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .header {
            margin-bottom: 30px;
            padding-bottom: 20px;
            border-bottom: 2px solid #eee;
        }
        .header h1 {
            color: #333;
            font-size: 24px;
            margin-bottom: 10px;
        }
        .header a {
            color: #4CAF50;
            text-decoration: none;
            font-size: 14px;
        }
        .header a:hover {
            text-decoration: underline;
        }
        .content {
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
        .back-link {
            display: inline-block;
            margin-top: 30px;
            padding: 10px 20px;
            background: #4CAF50;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            transition: background 0.3s;
        }
        .back-link:hover {
            background: #45a049;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>文章 ID: {{.article.ArticleID}}</h1>
            <a href="/articles">← 返回文章列表</a>
        </div>
        <div class="content">{{.highlightedContent | safeHTML}}</div>
        <a href="/articles" class="back-link">返回列表</a>
    </div>
</body>
</html>

