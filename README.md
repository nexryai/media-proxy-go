# media-proxy-go

## これは何
MisskeyのMediaProxyとして使える、Goで書かれた~~軽量なMediaProxy~~ でしたが方針がめちゃくちゃ変わりました。

### 依存関係
 - libvips

#### ビルド依存関係
 - libvips-dev
 - Go (**バージョン1.18以上**)

### 仕組み
本家のMediaProxyとは異なり、リクエストを一回ジョブキューにぶち込みダウンロードしてリサイズ圧縮してストレージに書き込んでから応答するようにしています。  
ジョブキューを使用することで同時実行数やタイムアウトを厳格に設定できるようになり、永続キャッシュの導入によりストレージを消費する代わりに高速な応答が可能となりました。  
キャッシュは72時間使用されないと消し飛ばされます。

#### ロードマップ
 - [x] ジョブキューを使う
 - [x] 永続キャッシュを導入する
 - [x] キャッシュの自動消し飛ばし
 - [ ] Nexkeyバックエンドと連携してキャッシュを事前生成するためのAPIを実装
 - [ ] キャッシュが事前生成されていないものはプロキシできないようにするモードを実装
   
### Tips
#### メモリ使用量が爆発する
Jemallocを使ってください

#### `PORT`環境変数
`PORT`を指定することでlistenするポートを指定できるようになりました。

#### docker-composeで動かす
Docker Hub に一応ビルド済みイメージを上げています。
```
services:
  app:
    image: docker.io/nexryai/mediaproxy-go:latest
    restart: always
    environment:
      - REDIS_ADDR=redis:6379
      - CACHE_DIR=/cache
    volumes:
      - proxy_cache:/cache
    ports:
      - 127.0.0.1:8080:8080
    depends_on:
      - redis

  redis:
    image: "redis:alpine"
    volumes:
      - redis_data:/data

volumes:
  proxy_cache:
    external: false

  redis_data:
    external: false
```

#### 制限
 - ~~だいぶ安定していますが、govipsを上手く扱えていないのが主な原因でまだ微妙に不安定な部分が残っています。~~ 大規模なリファクタリング中なので不安定です。
 - [Misskeyのメディアプロキシ仕様](https://github.com/misskey-dev/media-proxy/blob/master/SPECIFICATION.md)の大部分を網羅していますが、Nexkeyでの使用を想定して開発されているため完全な実装ではありません。
 - ポート80と443以外へのリクエストはドロップされます
 - （これは本家を使用しているときにも言えることですが）セキュリティ対策は基本的なものを行っていますが完全ではありません。プライベートなネットワークで実行することはあまり推奨されません。

