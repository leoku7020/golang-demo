## 後端開發環境

### 技術關鍵字
 - MySQL 8.0
 - PhpMyAdmin
 - Docker 24.0.7
 - Docker-Compose 2.23.3
 - Golang 1.19.0
 - protobuf v3.21.12
 - protoc-gen-go v1.28.1
 - protoc-gen-go-grpc v1.2.0

### 服務port
 - Aurora MySQL 8.0 (MySQL 服務)
   - `127.0.0.1:33060`
 - PhpMyAdmin - MySQL (DB tool)
   - `127.0.0.1:8080`
 - Golang RestFul API
   - `localhost:8000`
### 啟動環境
   - 先安裝 docker & docker-compose
   - 確認docker 有安裝好 `$ docker --version` 有出現版本即安裝成功
   - 確認docker-compose 有安裝好 `$ docker-compose --version` 有出現版本即安裝成功
   - `$ docker-compose up -d`
### DB GUI Tools
   - PhpMyAdmin - MySQL
### Make 指令
   - `$ make init` 複製 env
   - `$ make generate` 初始化grpc proto
### 安裝 protobuf v3.21.12 for Mac 如果安裝最新版 會有 proto 編譯過不了的問題
- brew install protobuf@21
- echo 'export PATH="/usr/local/opt/protobuf@21/bin:$PATH"' >> ~/.zshrc
- export LDFLAGS="-L/usr/local/opt/protobuf@21/lib"
- export CPPFLAGS="-I/usr/local/opt/protobuf@21/include"
- export PKG_CONFIG_PATH="/usr/local/opt/protobuf@21/lib/pkgconfig"
- protoc --version ; 顯示 libprotoc 3.21.12 為 ok

### 安裝 protoc-gen-go v1.28.1
- go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

### 安裝 protoc-gen-go-grpc v1.2.0
- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

### 本機環境建置 流程
- 先用 `make init` 複製出 .env 基本上不用改
- `make generate` 跑一下 proto
- docker 啟動 mysql & phpmyadmin 服務 `docker-compose up -d`
- 訪問 `localhost:8080` 看到 phpmyadmin mysql GUI 介面 新增 database name example 空表的即可
- 啟動服務 兩者擇一
  1. 用Goland跑專案
    - 打开 "Run/Debug Configurations" 窗口：在顶部菜单中选择 "Run" > "Edit Configurations..."。
    - 点击 "+" 按钮创建一个新的 Go 运行配置 選擇 go build
       - Run kind: package
       - Package path: demo/cmd/api/example
  1. Cmd 直接跑起來
    - `go run cmd/api/example/main.go`
- 會自動跑db migration 會新增兩張表 members (for Login) items (for ListItem) 且 自動新增假資料 

### 測試api
  - Login Success 帳密 admin/admin
    - ![login success.png](imgs%2Flogin%20success.png)
  - Login Req Failed
    - ![login_failed2.png](imgs%2Flogin_failed2.png)
  - Login Failed
    - ![login failed.png](imgs%2Flogin%20failed.png)
  - ListItems Success
    - ![listItem.png](imgs%2FlistItem.png)
  - ListItems Failed
    - ![listitem_failed.png](imgs%2Flistitem_failed.png)
### DB 資料
  - Members
    - ![db_members.png](imgs%2Fdb_members.png)
  - Items
    - ![db_items.png](imgs%2Fdb_items.png)