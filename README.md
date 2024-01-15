# media-proxy-go

## これは何
MisskeyのMediaProxyとして使える、Goで書かれた軽量なMediaProxy

### 依存関係
 - libvips

#### ビルド依存関係
 - libvips-dev
 - Go (**バージョン1.18以上**)

   
### Tips
#### メモリ使用量が爆発する
Jemallocを使ってください

#### `PORT`環境変数
`PORT`を指定することでlistenするポートを指定できるようになりました。

#### docker-composeで動かす
Docker Hub に一応ビルド済みイメージを上げています。
```
version: '3'
services:
  app:
    image: docker.io/nexryai/mediaproxy-go:latest
    ports:
      - 127.0.0.1:9090:8080
    restart: always
```

#### 制限
 - だいぶ安定していますが、govipsを上手く扱えていないのが主な原因でまだ微妙に不安定な部分が残っています。
 - [Misskeyのメディアプロキシ仕様](https://github.com/misskey-dev/media-proxy/blob/master/SPECIFICATION.md)の大部分を網羅していますが、Nexkeyでの使用を想定して開発されているため完全な実装ではありません。
 - ポート80と443以外へのリクエストはドロップされます
 - （これは本家を使用しているときにも言えることですが）セキュリティ対策は基本的なものを行っていますが完全ではありません。プライベートなネットワークで実行することはあまり推奨されません。

