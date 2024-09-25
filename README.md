# Golang-Short-Link

## How to design HTTP Router and Handler?


## How to add Middleware when dealing with HTTP processing?


## How to use Interface in Go to implement scalability design?


## How to generate short links using Redis auto-increase sequence?

**API:** 
1. POST /api/shorten
2. Get /api/info?shortlink=shortlink
3. GET/:shortlink - return 302 code

1. POST /api/shorten:
   - Params:
  
    | Name                 | Type                          |Description|
    |----------------------|-------------------------------|-------------|
    | url                  | string                        |Required. URL to shorten.|
    | expiration_in_minutes| int                           |Required. expiration of short link in minutes. e.g. value 0 represents permanent|

   - Response:
     ```
       {
          "shortlink": "p"
       }
     ```

 2. Get /api/info?shortlink=shortlink
    - Params:
        | Name                 | Type                          |Description|
        |----------------------|-------------------------------|-------------|
        | shortlink            | string                        |Required. id of shortened. e.g. P|
    - Response:
     ```
       {
          "url": "https://www.example.com",
          "created_at": "2024-09-14 20:59:50.46866058 +0800 CST",
          "expiration_in_minutes": 60
       }
     ```
 3. GET/:shortlink - return 302 code
 
       为什么用302，不用301？  答：302 是临时的，使用301会永久保存在用户的浏览器缓存中，这样无法跟踪用户的行为。

## How to use redis to convert link to shortlink?

1. **URL Hash Mapping** (URLHashKey)
- **目的**: 将 URL 的哈希值映射到短链接 ID。
- **工作流程**:
  
            1. 用户提供一个 URL（例如 https://example.com/some/long/url。
  
            2. 计算这个 URL 的 SHA1 哈希，假设得到 h = abc123。
  
            3. 使用 fmt.Sprintf(URLHashKey, h) 生成键 "urlhash:abc123:url"。
  
            4. 如果这个哈希值之前没有存储过，则进入下一步生成短链接 ID；如果存在，则直接返回存储的短链接 ID。


2. **Shortlink Mapping** (ShortlinkKey)
- **目的**: 将短链接 ID 映射到原始 URL。
- **工作流程**:
  
               1. 在生成短链接 ID 时，假设计算出的短链接 ID（经过 base62 编码）为 xyz456。
  
               2. 使用 fmt.Sprintf(ShortlinkKey, eid) 生成键 "shortlink:xyz456:url"。

               3. 将原始 URL（https://example.com/some/long/url）存储在这个键下。

3. **Shortlink Detail Mapping** (ShortlinkDetailKey)
- **目的**: 存储短链接的详细信息。
- **工作流程**:
  
               1. 使用 fmt.Sprintf(ShortlinkDetailKey, eid) 生成键，例如 "shortlinkxyz456:detail"。

               2.将短链接的详细信息（如原始 URL、创建时间、过期时间）序列化为 JSON，并存储在这个键下。



## Mapping process
1. 用户请求:
输入 URL → 计算哈希 → 检查 urlhash:abc123:url 键。
如果未找到，生成短链接 ID（如 xyz456）。
2. 存储数据:
将短链接 ID xyz456 存储到 shortlink:xyz456:url 键。
将哈希值与短链接 ID 关联，存储到 urlhash:abc123:url 键。
将短链接的详细信息存储到 shortlinkxyz456:detail 键。

