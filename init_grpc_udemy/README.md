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

## 通信方式別gRPCの詳細

## gRPCの応用
