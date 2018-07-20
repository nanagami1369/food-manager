package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Calendar struct {
	Day       [31]int8
	Shop      [31]string
	Food      [31]string
	Price     [31]int16
	Weekday   [31]int
	Messege   string
	Messegest Messegest
	Messege2  string
	Messege3  string
	Gool      int
	Times     Times
}
type Times struct {
	Year  int
	Month int
}
type Messegest struct {
	Err  string
	Date string
}

var Current int

func ftables(year, month int) string {
	//テーブルの名前を表示する
	return fmt.Sprint(year) + "_" + fmt.Sprint(month) + "_fdb"
}
func holiday(year, month, day int) int {
	//土曜か日曜か平日かを数字で返す(曜日の色つけに使用)
	// 日   	1
	// 土   	2
	// 平日 	0
	if weekday := fmt.Sprint(time.Date(year, time.Month(month), day, int(0), int(0), int(0), int(0), time.Local).Weekday()); weekday == "Sunday" {
		return 1
	} else if weekday == "Saturday" {
		return 2
	} else {
		return 0
	}
}
func checkday(month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if month%400 == 0 || (month%100 != 0 && month%4 == 0) {
			return 29
		} else {
			return 28
		}
	default:
		return 0
	}
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	//dbにつなぐ
	db, err := sql.Open("mysql", "pokemon:2exo4t@tcp(localhost:3306)/foodmdb")
	if err != nil {
		fmt.Println("open mysql failed")
		fmt.Println(fmt.Sprint(err))
	}
	defer db.Close()
	//今月のtable作成
	_, err = db.Exec("create table if not exists foodmdb." + ftables(int(time.Now().Year()), int(time.Now().Month())) + " (day tinyint not null default 1 primary key ,shop char(20),foodname char(20),price smallint default 0 not null);")
	if err != nil {
		fmt.Println("create table failed")
		fmt.Println(fmt.Sprint(err))
	}

	//カレンダーを作る
	Calenders := Calendar{}
	//現在日時に設定
	Calenders.Times.Year = time.Now().Year()
	Calenders.Times.Month = int(time.Now().Month())
	//データを取得する
	r.ParseForm()

	//カレンダーをめくる
	//ボタンの確認
	if button := fmt.Sprint(r.Form["operation"]); button == "[<]" {
		//戻るボタン
		Current--
	} else if button == "[>]" {
		//進むボタン
		Current++
	} else if button == "[現在の月に戻る]" {
		//戻るボタン
		Current = 0
	}
	//操作ボタン実行
	for count := 0; Current != count; {
		if Current > 0 {
			//進む
			if Calenders.Times.Month == 12 {
				Calenders.Times.Month = 1
				Calenders.Times.Year++
			} else {
				Calenders.Times.Month++
			}
			//戻る
			count++
		} else {
			if Calenders.Times.Month == 1 {
				Calenders.Times.Month = 12
				Calenders.Times.Year--

			} else {
				Calenders.Times.Month--
			}
			count--
		}
	}
	//テーブルがあるか確認
	count := 0
	err = db.QueryRow(fmt.Sprintf("select count(*) from information_schema.tables where table_schema=database() and table_name like '%d%%%d_fdb';", Calenders.Times.Year, Calenders.Times.Month)).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	//もしもなかったら
	if count == 0 {
		//未来に振り切れた場合１戻す
		if Current > 0 {
			if Calenders.Times.Month == 1 {
				Calenders.Times.Month = 12
				Calenders.Times.Year--
			} else {
				Calenders.Times.Month--
			}
			Current--
			Calenders.Messege = "それより先のデータはありません"
			Calenders.Messegest.Err = "警告"
			if Calenders.Shop[int(time.Now().Day())-1] == "" {
				Calenders.Messegest.Date = "登録"
			} else {
				Calenders.Messegest.Date = "修正"
			}

		} else if Current < 0 {
			//過去に振り切れた場合１すすめる
			if Calenders.Times.Month == 12 {
				Calenders.Times.Month = 1
				Calenders.Times.Year++
			} else {
				Calenders.Times.Month++
				Calenders.Messege = "それより前のデータはありません"
				Calenders.Messegest.Err = "警告"
			}
			Current++
		} else {
			Calenders.Messege = ""
			Calenders.Messegest.Date = "無し"
		}
	}
	//日付と土日指定だけ先に入力
	for x := 0; x < 31; x++ {
		Calenders.Day[x] = int8(x + 1)
		Calenders.Weekday[x] = holiday(Calenders.Times.Year, Calenders.Times.Month, int(x))
	}
	//データ持ってくる準備をする
	columns, err := db.Query("select * from " + ftables(Calenders.Times.Year, Calenders.Times.Month))
	if err != nil {
		fmt.Println("select table failed")
		fmt.Println(fmt.Sprint(err))
	}
	defer columns.Close()
	monthsam := 0
	daysam := 0
	//dbからデータを取得
	for columns.Next() {
		var day int8
		var shop string
		var food string
		var price int16
		if err := columns.Scan(&day, &shop, &food, &price); err != nil {
			fmt.Println(err)
		}
		if day <= 31 && day >= 1 {
			Calenders.Shop[day-1] = shop
			Calenders.Food[day-1] = food
			Calenders.Price[day-1] = price
			monthsam += int(Calenders.Price[day-1])
			daysam++
		} else {
			fmt.Println("Not day")
		}
	}

	//目標を処理する
	var month int
	var Calculation int
	//目標入力
	err = db.QueryRow(fmt.Sprint("select month from gool;")).Scan(&month)
	if err != nil {
		fmt.Println(err)
	}
	Calenders.Gool = month
	month -= monthsam
	day := checkday(int(time.Now().Month()))
	Calculation = month / (day - daysam)
	//目標変更
	if fmt.Sprint(r.Form["monthgool"]) != "[]" {
		Calenders.Messege2 = ""
		if m, _ := strconv.Atoi(r.Form.Get("monthgool")); m < 31 {
			//月がマイナスの場合
			Calenders.Messege2 = "月目標がマイナスになります"
		} else if monthsam > m {
			Calenders.Messege2 = "もうすでに使い切ってます"
		} else {
			//結果入力
			_, err := db.Exec(fmt.Sprintf("update gool set month=%d ;", m))
			if err != nil {
				fmt.Println("Not update")
			}
			//値の更新を入れる
			err = db.QueryRow(fmt.Sprint("select month from gool;")).Scan(&month)
			if err != nil {
				fmt.Println(err)
			}
			Calenders.Gool = month
			month -= monthsam
			day := checkday(int(time.Now().Month()))
			Calculation = month / (day - daysam)
			Calenders.Messege2 = fmt.Sprintf("月目標残り%d円です", month)
			Calenders.Messege3 = fmt.Sprintf("一日につき%d円使えます", Calculation)
		}
	} else {
		Calenders.Messege2 = fmt.Sprintf("月目標残り%d円です", month)
		Calenders.Messege3 = fmt.Sprintf("一日につき%d円使えます", Calculation)
	}
	//今日のデータを入力していない∧今月のカレンダーを開いている
	if Calenders.Shop[int(time.Now().Day())-1] == "" && time.Now().Year() == Calenders.Times.Year && int(time.Now().Month()) == Calenders.Times.Month && Calenders.Messege != "それより先のデータはありません" {
		switch {
		case fmt.Sprint(r.Form["shop"]) == "[]" && fmt.Sprint(r.Form["food"]) == "[]" && fmt.Sprint(r.Form["price"]) == "[]":
			{
				//データが全て空
				switch hour := time.Now().Hour(); {
				case hour < 12:
					Calenders.Messege = "昼食を食べたらデータを入力して下さい"
				case hour < 15:
					Calenders.Messege = "今日の分のデータを入力して下さい"
				case hour < 20:
					Calenders.Messege = "早急に今日の分の入力をして下さい"
				case hour < 23:
					Calenders.Messege = "お昼を食べ忘れたのですか?"
				default:
					Calenders.Messege = "入力しろ"
				}
				Calenders.Messegest.Date = "登録"
			}
		case fmt.Sprint(r.Form["shop"]) != "[]" && fmt.Sprint(r.Form["food"]) != "[]" && fmt.Sprint(r.Form["price"]) != "[]":
			{
				//データが全て埋まっていた場合
				if x, _ := strconv.Atoi(r.Form.Get("price")); x < 0 || x > 10000 {
					//値段が明らかに高すぎたり、安すぎたりしないか?
					Calenders.Messege = "不正な値です"
					Calenders.Messegest.Date = "登録"
				} else {
					//結果入力
					_, err := db.Exec(fmt.Sprintf("insert into %s (day,shop,foodname,price)values (%d,'%s','%s',%s);", ftables(time.Now().Year(), int(time.Now().Month())), time.Now().Day(), r.Form.Get("shop"), r.Form.Get("food"), r.Form.Get("price")))
					if err != nil {
						fmt.Println("Not insert")
					}
					//値の更新を入れる
					day := int(time.Now().Day())
					err = db.QueryRow("select shop,foodname,price from "+ftables(int(time.Now().Year()), int(time.Now().Month()))+" where day = "+fmt.Sprint(time.Now().Day())+";").Scan(&Calenders.Shop[day-1], &Calenders.Food[day-1], &Calenders.Price[day-1])
					if err != nil {
						fmt.Println(err)
					}
					Calenders.Messege = "データの入力が完了しました"
					Calenders.Messegest.Date = "修正"
				}
			}
		default:
			{
				//一部データが抜けている場合
				var serch string
				if fmt.Sprint(r.Form["shop"]) == "[]" {
					serch += "店の名前、"
				}
				if fmt.Sprint(r.Form["food"]) == "[]" {
					serch += "食べたもの、"
				}
				if fmt.Sprint(r.Form["price"]) == "[]" {
					serch += "値段、"
				}
				Calenders.Messege = fmt.Sprintf("%sが入力されていません", serch)
				Calenders.Messegest.Date = "登録"
			}
		}
	} else if Calenders.Shop[int(time.Now().Day())-1] != "" && time.Now().Year() == Calenders.Times.Year && int(time.Now().Month()) == Calenders.Times.Month {
		//今日のデータ入力済み∧今月カレンダーを開いている
		switch {
		case fmt.Sprint(r.Form["shop"]) == "[]" && fmt.Sprint(r.Form["food"]) == "[]" && fmt.Sprint(r.Form["price"]) == "[]":
			{
				//データが全て空
				if Calenders.Messegest.Err != "警告" {
					Calenders.Messege = "今日のデータは入力済みです"
					Calenders.Messegest.Date = "修正"
				}
			}
		case fmt.Sprint(r.Form["shop"]) != "[]" && fmt.Sprint(r.Form["food"]) != "[]" && fmt.Sprint(r.Form["price"]) != "[]":
			{
				//データが全て埋まっていた場合
				if x, _ := strconv.Atoi(r.Form.Get("price")); x < 0 || x > 10000 {
					//値段が明らかに高すぎたり、安すぎたりしないか?
					Calenders.Messege = "不正な値です"
					Calenders.Messegest.Date = "登録"
				} else {
					//結果入力
					_, err := db.Exec(fmt.Sprintf("update %s set shop = '%s',foodname = '%s',price = '%s' where day = %d;", ftables(time.Now().Year(), int(time.Now().Month())), r.Form.Get("shop"), r.Form.Get("food"), r.Form.Get("price"), time.Now().Day()))
					if err != nil {
						fmt.Println("Not update")
					}
					//値の更新を入れる
					day := int(time.Now().Day())
					err = db.QueryRow("select shop,foodname,price from "+ftables(int(time.Now().Year()), int(time.Now().Month()))+" where day = "+fmt.Sprint(time.Now().Day())+";").Scan(&Calenders.Shop[day-1], &Calenders.Food[day-1], &Calenders.Price[day-1])
					if err != nil {
						fmt.Println(err)
					}
					Calenders.Messege = "データの修正が完了しました"
					Calenders.Messegest.Date = "修正"

				}
			}
		default:
			{
				//一部データが抜けている場合
				var serch string
				if fmt.Sprint(r.Form["shop"]) == "[]" {
					serch += "店の名前、"
				}
				if fmt.Sprint(r.Form["food"]) == "[]" {
					serch += "食べたもの、"
				}
				if fmt.Sprint(r.Form["price"]) == "[]" {
					serch += "値段、"
				}
				Calenders.Messege = fmt.Sprintf("%sが入力されていません", serch)
				Calenders.Messegest.Date = "修正"

			}
		}
	} else if time.Now().Year() != Calenders.Times.Year || int(time.Now().Month()) != Calenders.Times.Month {
		//今月のカレンダーを開いていない
		//特に何もしない(メッセージは振り切れた際の表示のみ)
	}

	//テンプレートをパース
	tmpl := template.Must(template.ParseFiles("./view/view.html"))
	tmpl.Execute(w, Calenders)
}
func main() {
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(":8080", nil)
}
