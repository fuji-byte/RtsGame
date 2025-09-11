<!-- #Redis の利用
redis を入れていないなら、redis を導入する
sudo apt update
sudo apt install redis-server
sudo service redis-server start -->

<!-- wsl の redis の実行
sudo service redis-server start -->

golang 実行
air

<!-- http://localhost:8080/login?user_id=testuser123 -->

<!-- Redis クライアント確認(wsl)
redis-cli
keys \*
get session:クライアント id -->

最優先実装機能
ゲーム進行
せる接続後の動作（対象の hp が削れ、０になったら占領）

作成予定機能
Cell HP 自動増加（成長率）
sync あたりの最適化も
room 終了通知(プレイヤーに)
room キュー、ランダム参加など
observer の仕様 ゲーム開始に入れるか、ルーム内に通知するか
cell 線の双方向　おもにフロント
更新した部分だけブロードキャストする
観戦者系（観戦者がいたらゲームを終了しないなど）
ルーム時間終了後の処理
cell 座標の修正（フロントとバックエンド両方）

error 報告
なし

追加要素（オプション）
ランキング
フレンド機能
guest でもトークンを生成してユーザー認証を行う
savelog function の実装
