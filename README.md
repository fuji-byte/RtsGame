#Redis の利用
redis を入れていないなら、redis を導入する
sudo apt update
sudo apt install redis-server
sudo service redis-server start

wsl の redis の実行
sudo service redis-server start

golang 実行
air

http://localhost:8080/login?user_id=testuser123

Redis クライアント確認(wsl)
redis-cli
keys \*
get session:クライアント id

作成予定機能
guest でもトークンを生成してユーザー認証を行う
フレンド機能
sync 当たりの最適化も
gameroom をキーで管理するか
host が抜けたときの処理書いてある？
