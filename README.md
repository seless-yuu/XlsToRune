# XlsToRune
ゲームマスターデータとして記述されたExcelファイルをゲーム内パラメータとして解釈する為の<br>
jsonファイル(rune形式とします)へと変換するツールを扱うリポジトリです

## インストール
### バイナリをダウンロードする
Releases ページから各バイナリをダウンロードしてください<br>
### リポジトリからビルドする
リポジトリをクローンしてビルドします。golang環境が必要です<br>
#### Windows環境
build_[各プラットフォーム名].bat のシェルスクリプトを起動してください<br>
### Mac(Apple Silicon)環境
build_arm64.sh のシェルスクリプトを起動してください<br>
#### Linux,wsl環境
build_[各プラットフォーム名].sh のシェルスクリプトを起動してください<br>

buildフォルダへ各バイナリが出力されます<br>

## 使用方法
### コマンドラインオプション
XlsToRune [-i 入力ファイル] [-o 出力ファイル] [-enum enum出力フォルダ] [-class 型クラス出力フォルダ] [-enum_namespace enumのnamespace名]<br>

## タイプテーブル
### 使用できる型
enum<br>
int<br>
int2<br>
int3<br>
float<br>
float2<br>
float3<br>
float4<br>
string<br>

・"#"が先頭につく型の列はコメント扱いとして出力されません<br>
・各型は文字列として格納されます。その解釈は扱う側に任せられます<br>

## ChangeLog

### v1.02
- int2,3 float2,3,4 の型を追加
- Config.ScriptableObjectDirectoryが未設定の場合は
  runeファイルと同じフォルダに出力される扱いで出力クラスを修正
