# protoc-gen-proto

<p align="center">
  <a href="README.md">English</a> •
  <a href="README.ja.md">日本語 (Japanese)</a>
</p>

これは protoc プラグインのテンプレートです。
デフォルトのプラグイン実装は、`.proto` ファイルから `.proto` ファイルを生成します。
入力された `.proto` ファイルと同じ `.proto` ファイルを生成する処理を把握することで、
新しく protoc プラグインを作成する際に役立ちます。

## Tasks

```shell
# generate code
make generated
```

```shell
# run test
make test
```

```shell
# generate testdata
make testdata
```

```shell
# update test expected data
go run ./cmd/protoc-gen-proto < testdata/stdin.bin.pb > testdata/stdout.bin.pb
```

```shell
# formatting .proto files
make format-proto
```
