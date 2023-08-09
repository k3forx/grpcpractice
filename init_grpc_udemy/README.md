# Go言語で学ぶ実践gRPC入門

## Protocol Buffersとは

- Googleによって2008年にオープンソース化されたスキーマ言語

## Protocol Buffersの基礎

### Protocol Buffersの特徴

- **gRPCのデータフォーマット**として使用されている
- **プログラミング言語からは独立**しており、様々な言語に変換可能
- **バイナリ形式にシリアライズ**するので、サイズが小さく高速な通信が可能
- **型安全**にデータのやり取りが可能
- **JSONに変換**することも可能

### Protocol Buffersを使用した開発の進め方

1. スキーマの定義
1. 開発言語のオブジェクトを自動生成
1. バイナリ形式にシリアライズ

### `message` とは

- 複数フィールドを持つことができる型定義
  - それぞれのフィールドはスカラ型もしくはコンポジット型
- 各言語のコードとしてコンパイルした場合、構造体やクラスとして変換される
- フィールド番号 (タグ番号) はシリアライズされた後のデータでフィールドを区別するために使用される

### スカラー型

- [Scalar Value Types](https://protobuf.dev/programming-guides/proto3/#scalar)

### Tag

- Protocol Buffersではフィールドはフィールド名ではなく、タグ番号によって識別される
- 重複は許されず、一意である必要がある
- タグの最小値は1、最大値は2^29-1 (536,870,911)
- 1-15万までは1byteで表すことができるので、よく使うフィールドには1-15番を割り当てる
- タグは連番にする必要はないので、あまり使わないフィールドはあえて16番以降を割り当てることも可能
- タグ番号を予約するなど、安全にProtocol Buffersを使用する方法も用意されている

### 列挙型

- `0` をunknownとすることが多い
- [Enumerations](https://protobuf.dev/programming-guides/proto3/#enum)

```proto
enum Occupation {
  OCCUPATION_UNKNOWN = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}
```

### その他の型

- [Specifying Filed Labels](https://protobuf.dev/programming-guides/proto3/#field-labels)
- [Oneof](https://protobuf.dev/programming-guides/proto3/#oneof)

- `repeated`: 0以上の繰り返しの値を扱う。順番は保持される。
- `map`: マップ。`map<string, int32> hoge = 4;` のような記述をする。
- `oneof`: いずれかの値が入る時に使う。

### デフォルト値

- 定義したmessageでデータをやり取りする際に、定義したフィールドがセットされていない場合、そのフィールドのデフォルト値が設定される
- デフォルト値は型によって決定される

## GoでのProtocol Buffersの操作

### protoファイルのコンパイル

- `-IPATH, --proto_path=PATH`: protoファイルのimport文のバスを指定する。複数の場合は `:` を使う。
- `--go_out=OUT_DIR`: Go言語への変換。変換したファイルの保存先を指定する。

```bash
$ protoc --version
libprotoc 23.4

protoc -I. --go_out=. proto/*.proto
```

protoファイルに `option go_package = "./path"` を記載し、コンパイルしたファイルのパッケージ名を指定する

## gRPCの基礎

### gRPCの概要

- Googleが開発したRPCのためのプロトコル

### gRPCの特徴

- データフォーマットにProtocol Buffersを使用
- IDLからサーバー側、クライアント側に必要なソースコードを生成
- 通信にはHTTP/2を使用
- 特定の言語やプラットフォームに依存しない

### gRPCが適したケース

- マイクロサービス間の通信
- モバイルユーザーが利用するサービス
  - 通信量が削減できるため、通信容量制限にかかりにくい
- 速度が求められる場合

### HTTP/1.1の課題

- リクエストの多重化
  - 1リクエストに対して1レスポンスという制約がある
- プロトコルオーバーヘッド
  - クッキーやトークンなどを毎回リクエストヘッダに付与してリクエストするため、オーバーヘッドが大きくなる

### HTTP/2の特徴

- ストリームという概念を導入
  - 1つのTCP接続を用いて、複数のリクエスト/レスポンスのやり取りが可能
  - TCP接続を減らすことができるので、サーバーの負荷軽減
- ヘッダーの圧縮
  - ヘッダーをHPACKという圧縮方式で圧縮し、さらにキャッシュを行うことで、差分のみを送受信することで効率化
- サーバープッシュ
  - クライアントからのリクエストなしにサーバーからデータを送信できる
  - 事前に必要と思われるリソースを送信しておくことで、ラウンドトリップの回数を削減しリソース読み込みまでの時間を短縮

### Serviceとは

- RPC (メソッド) の実装単位
  - サービス内に定義するメソッドがエンドポイントになる
  - 1サービス内に複数のメソッドを定義できる
- サービス名、メソッド名、引数、戻り値を定義する必要がある
- コンパイルしてgoファイルに変換すると、インターフェイスとなる
  - アプリケーション側でこのインターフェイスを実装する

### gRPCの通信方式

- Unary RPC
- Server Streaming RPC
- Client Streaming RPC
- Bidirectional Streaming RPC

#### Unary RPC

- 1リクエスト1レスポンスの方式
- 通常の関数のコールのように扱える
- 用途
  - API

#### Server Streaming RPC

- 1リクエスト、複数レスポンスの方式
- クライアントはサーバーから送信完了の信号が送信されるまでストリームのメッセージを読み続ける
- 用途
  - サーバーからのプッシュ通知など

#### Client Streaming RPC

- 複数リクエスト、1レスポンスの方式
- サーバーはクライアントからリクエスト完了の信号が送信されるまでストリームからメッセージを読み続け、レスポンスを返さない
- 用途
  - 大きなファイルのアップロードなど

#### Bidirectional Streaming RPC

- 複数リクエスト、複数レスポンスの方式
- クライアントとサーバーのストリームが独立しており、リクエストとレスポンスはどのような順序でも良い
- 用途
  - チャットやオンライン対戦ゲームなど

## 通信方式別gRPCの詳細

### Unary RPC

- コンパイル

```bash
$ protoc -I. --go_out=. --go-grpc_out=. proto/*.proto
```

- サーバー起動

```bash
$ go run server/main.go
```

- クライアント起動

```bash
$ go run client/main.go
filenames: [name.txt sports.txt]
```

## gRPCの応用
