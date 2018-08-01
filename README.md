# food-manager  
昼メシ管理ソフト  

事前に次のmysqlで次の3つを実行しておいて下さい  
CREATE DATABASE foodmdb;  
create table gool(month int not null default 10000);  
create user 'pokemon'@'localhost'identified by '2exo4t';  
GRANT ALL PRIVILEGES ON `foodmdb`.* TO 'pokemon'@'localhost';  
