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
