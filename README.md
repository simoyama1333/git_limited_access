## git_limited_access

GitHubのプライベートリポジトリをパスワードや期限設定をして公開するためのアプリケーションです。  
不特定多数に公開はしたいが、パブリックで全世界には公開したくない。cloneなどもされたくはない。
しかし不特定多数なので一人一人招待は無理という用途のために、GitHubにはそのような機能がなかったのでそのためのアプリケーションを作りました。  
  
### サンプル

サンプルは[こちら](http://35.233.244.144/)  
見ることのできるリポジトリはgit_limited_accessとMineSweeper3Dの2つです。  
MineSweeper3Dのパスワードは「Unity」です。  
フロントがReactでサーバーがGolangで動いています。DBはMongo  

### 使い方（local docker）  
  
`bash init.sh`  
`docker-compose up -d`  
`docker exec -it git_limited bash`  
`cd src/git_limited_access/golang`  
  
公開したいリポジトリを登録します。GitHubトークンを予め発行しておいてください。  
右上のSetting→Developer settings→Personal access tokensでリポジトリの権限をつけて発行です。  
以下のコマンドでクローリングしてリポジトリの構造を取得します。  
  
`export GIT_USERNAME="yourname"`  
`export GIT_TOKEN="yourtoken"`  
`./repo/repo -repo yourrepo -expire 2020-02-02 -password yourpassword`  
  
終わったらmainを実行してサーバーを立てます。完了です。  
  
`./main`
`service nginx start`

うまくいかなかった場合、ほぼReactの方に問題があるのでnodejsとnpmを入れて、
`npm install`  
`npm run buld`  
をしてみてください。
