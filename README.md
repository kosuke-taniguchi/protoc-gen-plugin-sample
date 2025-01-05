## 概要
protocコマンドからgotemplateを利用してコード生成をするサンプルスクリプト
db層のCRUD処理をprotoのmessageから自動生成する

## 使い方
**自作プラグインのビルド**
```
make install
```

**コード生成**

`./gen/go/mysql`にdb層の実装が、`./gen/go/hoge`内にはprotoのmessageが生成される

```
make protoc-plugin
```
