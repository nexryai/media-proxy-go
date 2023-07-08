# media-proxy-go

## これは何
MisskeyのMediaProxyとして使える、Goで書かれた軽量なMediaProxy

### 依存関係
 - libwebp

#### ビルド依存関係
 - libwebp-dev (RHEL系ならlibwebp-devel)
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
 - 現在開発段階なので、本番環境での利用は推奨されません
 - [Misskeyのメディアプロキシ仕様](https://github.com/misskey-dev/media-proxy/blob/master/SPECIFICATION.md)にはなるべく追従していますが、完全な実装ではありません。以下のような制限があります。
   * staticが指定されていない場合、gif画像とAnimated WebPは無変換でプロキシされます
   * プロキシ対象がapngかつstaticが指定されていない場合、本家と同様そのまま無変換で応答されますが、staticの場合は静止webpに変換されます。
   * フォールバック画像は実装されていません（ToDo）
   * ポート80と443以外へのリクエストはドロップされます
   * 変換クエリが存在しない場合の挙動が本家とは異なります。変換クエリが存在しないかつプロキシ対象のファイルが画像ファイルである場合、webpへのリサイズなしでの変換が行われます。画像ファイルではないかつブラウザセーフなファイル形式の場合そのままプロキシされます。
   * 30MB以上のファイルはプロキシできません
   * （これは本家を使用しているときにも言えることですが）セキュリティ対策は基本的なものを行っていますが完全ではありません。プライベートなネットワークで実行することはあまり推奨されません。

