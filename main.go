package main

import (
	"database/sql"
	"fmt"
	_ "mysql"
	"net/http"
)

//定义全局变量
var (
	Db       *sql.DB
	err      error
	username string
	password string
)

//打开数据库，我用的mysql，创建了database：login，然后里面有一张table: logincheck,然后有两个属性username和password
func init() {
	Db, err = sql.Open("mysql", "root:123123@/login")
	if err != nil {
		fmt.Println(err)
	}
}

//这是个检查注册信息的函数，思路是：
//传入注册的用户名username。用sql语句把所有的username导入一个切片
//遍历切片看传入的username是否和切片中元素重复，如果重复则传出false，否则true
func check_register(username string) bool {
	//写sql语句
	sqlStr := "select username from logincheck "
	//执行
	row, err := Db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
	}
	//创建切片
	n := make([]string, 30)

	for row.Next() {
		//声明
		var id string
		err := row.Scan(&id)
		if err != nil {
			fmt.Print(err)
		}
		n = append(n, id)

	}

	for _, i := range n {
		if i == username {
			return false
		}
	}

	return true

}

//这是个检查登录信息的函数。如果输入的用户名和密码与数据库中的一组用户名密码重复，则传出false，否则true
func check_login(username string, password string) bool {
	//写sql语句
	sqlStr := "select username,password from logincheck "
	//执行
	row, err := Db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
	}
	//创建切片

	for row.Next() {
		//声明
		var a, b string

		err := row.Scan(&a, &b)
		if err != nil {
			fmt.Print(err)
		}

		if a == username && b == password {
			return true
		}

	}

	return false

}

//如果注册用户名不重复，则将此用户名&密码导入数据库
func insert(username string, password string) {
	sqlStr := "insert into logincheck(username,password) values(?,?)"
	_, err := Db.Exec(sqlStr, username, password)
	if err != nil {
		fmt.Println("导入异常")
	} else {
		fmt.Println("导入成功")
	}
}

//注册端口
func register(w http.ResponseWriter, r *http.Request) {

	username = r.FormValue("username")
	password = r.FormValue("password")
	a := check_register(username)
	if a {
		//如果注册用户名不重复，则将此用户名&密码导入数据库
		go insert(username, password) //此insert函数在上面
		fmt.Fprintln(w, "恭喜您注册成功")
	} else {
		fmt.Fprintln(w, "抱歉，用户名重复，请重新选择用户名")
	}

}

//登录端口
func login(w http.ResponseWriter, r *http.Request) {
	username = r.FormValue("username")
	password = r.FormValue("password")
	a := check_login(username, password)
	if a {
		fmt.Fprintln(w, "恭喜您登录成功")
	} else {
		fmt.Fprintln(w, "用户名或密码错误，请重新输入")
	}

}

func main() {
	servemux := http.NewServeMux()
	serve := &http.Server{
		Addr:    ":2233",
		Handler: servemux,
	}

	servemux.HandleFunc("/login", login)
	servemux.HandleFunc("/register", register)
	serve.ListenAndServe()

}
