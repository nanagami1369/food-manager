# food-manager  
昼メシ管理ソフト  
  
事前に次のmysqlで次の5つを実行しておいて下さい  
CREATE DATABASE foodmdb;  
create table gool(month int not null default 10000);  
create user 'pokemon'@'localhost'identified by '2exo4t';  
GRANT ALL PRIVILEGES ON `foodmdb`.* TO 'pokemon'@'localhost';  
insert into gool (month) values(10000);  

# 使い方  
上記4つを行った上で  

./main  
http://localhost:8080/ につなぐ  

店、昼飯、値段を入力、Enter  
で今日のデータを入力
日付の左右でカレンダーの月を変える
#おまけ  
月目標を設定できる  
目標に対する残り残金と一日に付きいくら使えるのか計算  
フォームに金額を入れて送信すると目標額が更新される
