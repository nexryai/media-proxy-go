# media-proxy-go

## これは何
MisskeyのMediaProxyとして使える、Goで書かれた軽量なMediaProxy

### 依存関係
 - libmagickwand-7

#### ビルド依存関係
 - libmagickwand-7-dev
 - Go (**バージョン1.18以上**)

   
### Tips
#### メモリ使用量が爆発する
`MALLOC_TRIM_THRESHOLD_`と`MALLOC_MMAP_THRESHOLD_`をかなり小さい値に下げないと永遠にメモリを使い続けます。原因は謎です。強い人教えてください（）

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
    environment:
      - MALLOC_TRIM_THRESHOLD_=16
      - MALLOC_MMAP_THRESHOLD_=16
```

#### `DEBUG_MODE`
`DEBUG_MODE`はデバッグログを出力するモードです。

#### 制限
 - 現在開発段階なので、本番環境での利用はあまり推奨されません（misskey.sda1.netでは実際に運用していますが）
 - [Misskeyのメディアプロキシ仕様](https://github.com/misskey-dev/media-proxy/blob/master/SPECIFICATION.md)にはなるべく追従していますが、完全な実装ではありません。以下のような制限があります。
   * フォールバック画像は実装されていません（ToDo）
   * ポート80と443以外へのリクエストはドロップされます
   * 30MB以上のファイルはプロキシできません
   * （これは本家を使用しているときにも言えることですが）セキュリティ対策は基本的なものを行っていますが完全ではありません。プライベートなネットワークで実行することはあまり推奨されません。

