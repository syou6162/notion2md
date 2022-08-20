## notion2md
notionの特定のページをmarkdownとして出力させるコマンドラインツールです。

### install

```
% go install github.com/syou6162/notion2md
```

### setup
環境変数`NOTION_TOKEN`を設定しましょう。tokenの取得方法は[こちら](https://www.redgregory.com/notion/2020/6/15/9zuzav95gwzwewdu1dspweqbv481s5)を参照。

### usage 
以下のコマンドでnotionのページ内容をexportし、markdownとして標準出力に出力します。

```
% notion2md --page_id https://www.notion.so/your_user_name/your_page_id
```

markdownの内容をすぐにvimで編集できるように、以下をブックマークレットに追加しておくと便利です。

```
javascript:window.prompt("Copy and paste command", "notion2md --page_id " + location.href + " | vim -");
```
