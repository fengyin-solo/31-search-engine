# 倒排索引搜索引擎（search-engine）

纯 Go 标准库实现的倒排索引搜索引擎，覆盖文档管理、索引构建、分词、TF·IDF 打分、同义词扩展、拼写纠错、搜索建议、加权规则、分面搜索、热门搜索与查询日志的完整检索闭环。零第三方依赖，开箱即跑。

## 技术栈

- Go 1.22+，仅使用标准库（`net/http` + 标准库）
- 标准工程分层：`cmd` / `internal`（app/config/model/store/service/handler）/ `pkg`
- 内存存储（`store.MemoryStore`），线程安全
- 前端页面：`web/` 下纯 HTML + 原生 JS + CSS，零 CDN

## 核心算法

- **分词（Tokenizer）**：统一小写、按非字母数字切分、剔除停用词、中文按单字切分。
- **倒排索引**：文档 → 分词 → `term → postings(docID, TF)`，维护词条 `DocCount/TotalTF` 与索引 `DocCount`。
- **相关度打分**：`TF·IDF`，`idf = log(1 + N / (1 + df))`；叠加「加权规则」（命中标题/正文的词条乘以 boost 权重）。
- **同义词扩展**：查询词命中同义词库时自动扩展，提升召回。
- **拼写纠错**：词典精确匹配 + 莱文斯坦编辑距离（容错 ≤2）。
- **搜索建议**：按前缀从建议词库与词典词条合并召回，按权重排序。
- **分面搜索**：结果按文档分类聚合计数。
- **高亮**：标题中命中的查询词用 `<em>` 包裹。

## 运行

```bash
cd origin
go run ./cmd/server
# 或
go build -o server ./cmd/server && ./server
```

环境变量：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | 8080 | 监听端口 |
| `ADDR` | `:8080` | 监听地址（优先于 PORT） |
| `MAX_PAGE_SIZE` | 100 | 分页最大页大小 |
| `AUTH_TOKEN` | `dev-token` | API 鉴权令牌（`X-Auth-Token` 请求头） |
| `RATE_LIMIT` | 600 | 每 IP 每分钟限流次数 |
| `LOG_LEVEL` | info | 日志级别 debug/info/warn/error |

浏览器访问 `http://localhost:8080/` 查看看板；先调用 `POST /api/seed/demo` 初始化演示数据。

## 统一响应

```json
{"code":0,"message":"ok","data":{...}}
```

错误映射：`ValidationError→400`、`ErrNotFound→404`、`ErrConflict→409`、其他→500。

## API 一览

鉴权：除静态页面外，所有 `/api/` 接口需在请求头携带 `X-Auth-Token`。

### 文档 /api/documents
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/documents | 创建文档 |
| GET | /api/documents | 列表（status/source/category_id/keyword） |
| GET | /api/documents/{id} | 详情 |
| PUT | /api/documents/{id} | 更新 |
| DELETE | /api/documents/{id} | 删除（软删并移除索引） |
| POST | /api/documents/{id}/category | 设置分类 |

### 索引 /api/indexes
状态机：`created→ready→deleted`。
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/indexes | 创建索引 |
| GET | /api/indexes | 列表 |
| GET | /api/indexes/{id} | 详情 |
| PUT | /api/indexes/{id} | 更新 |
| DELETE | /api/indexes/{id} | 删除 |
| POST | /api/indexes/{id}/activate | 激活 |
| POST | /api/indexes/{id}/deactivate | 停用 |
| POST | /api/indexes/{id}/rebuild | 重建索引 |

### 分词器 /api/analyzers
| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/GET{id}/PUT/DELETE | /api/analyzers | 分词器 CRUD |

### 停用词 /api/stop-words
| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/DELETE{id} | /api/stop-words | 停用词管理 |

### 同义词 /api/synonyms
| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/GET{id}/PUT/DELETE | /api/synonyms | 同义词管理 |

### 词条 /api/terms
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/terms | 列表（index_id 筛选，按文档频率排序） |
| GET | /api/terms/{id} | 详情 |

### 倒排记录 /api/postings
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/postings | 列表（index_id/term/doc_id 筛选） |
| GET | /api/postings/{id} | 详情 |

### 查询记录 /api/query-logs
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/query-logs | 列表 |
| GET | /api/query-logs/{id} | 详情 |
| DELETE | /api/query-logs/{id} | 删除 |

### 热门搜索 /api/hot-searches
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/hot-searches | 列表（按热度排序） |
| GET | /api/hot-searches/{id} | 详情 |
| DELETE | /api/hot-searches/{id} | 删除 |

### 分类 /api/categories
| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/GET{id}/PUT/DELETE | /api/categories | 分类管理 |

### 搜索建议 /api/suggestions 与 /api/suggest
| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/GET{id}/PUT/DELETE | /api/suggestions | 建议词管理 |
| GET | /api/suggest?prefix=xx | 前缀自动补全 |

### 拼写纠错 /api/spell-corrections 与 /api/correct
| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/GET{id}/PUT/DELETE | /api/spell-corrections | 纠错词对管理 |
| GET | /api/correct?word=xx | 查询词纠错 |

### 加权规则 /api/boost-rules
| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/GET{id}/PUT/DELETE | /api/boost-rules | 加权规则管理 |

### 搜索
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/search?index_id=&q=&top_k= | 全文搜索（TF·IDF + 加权 + 高亮） |
| GET | /api/search/facets?index_id=&q= | 分面搜索（按分类聚合） |
| POST | /api/index-documents | 将文档索引到索引（body: index_id, doc_id） |

### 统计 /api/stats
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/stats/overview | 引擎总览 |
| GET | /api/stats/indexes | 索引文档/词条统计 |
| GET | /api/stats/queries | 查询统计 + 热门搜索 |
| GET | /api/stats/terms | TOP 词条 |
| GET | /api/stats/categories | 分类/来源统计 |

### 其他
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/seed/demo | 初始化演示数据 |
| GET | /api/export/summary | 导出全量快照汇总 |
