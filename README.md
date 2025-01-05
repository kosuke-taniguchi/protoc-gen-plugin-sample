## 概要
protocコマンドからgotemplateを利用してコード生成をするサンプルスクリプト

## 使い方
**自作プラグインのビルド**
```
make install
```

**コード生成**

`./gen/go/mysql`にgotemplateを利用したコードが、`./gen/go/hoge`内にはprotoのmessageが生成される

```
make protoc-plugin
```
