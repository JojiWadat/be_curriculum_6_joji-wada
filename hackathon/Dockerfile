# ビルドステージ
FROM golang:1.23.3 AS build

# 作業ディレクトリを設定
WORKDIR /app

# モジュールファイルをコピーして依存関係を取得
COPY go.mod ./
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . ./

# 実行可能ファイルをビルド
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 実行ステージ（小さいイメージを使用）
FROM alpine:3.18

# 必要なバイナリをコピー
WORKDIR /app
COPY --from=build /app/main .

# 実行可能ファイルを実行
CMD ["./main"]
