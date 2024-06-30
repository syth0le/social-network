package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"

	"github.com/syth0le/social-network/internal/model"
)

// файлик для отладки

func main() {
	// Connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	dialer := tarantool.NetDialer{
		Address:  "127.0.0.1:3301",
		User:     "admin",
		Password: "password",
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		fmt.Println("Connection refused:", err)
		return
	}
	defer func() {
		fmt.Println("Connection is closed")
		conn.CloseGraceful()
	}()

	tuples := [][]interface{}{
		{1, "Roxette", "name1", "user1", "pass", "", "", ""},
		{5, "Scorpions", "name2", "user2", "pass", "", "", ""},
		{6, "Ace of Base", "name3", "user3", "pass", "", "", ""},
		{7, "The Beatles", "name4", "user41", "pass", "", "", ""},
	}
	var futures []*tarantool.Future
	for _, tuple := range tuples {
		request := tarantool.NewInsertRequest("users").Tuple(tuple)
		futures = append(futures, conn.Do(request))
	}
	fmt.Println("Inserted tuples:")
	for _, future := range futures {
		result, err := future.Get()
		if err != nil {
			fmt.Println("Got an error:", err)
		} else {
			fmt.Println(result)
		}
	}

	req2 := tarantool.NewCallRequest("search_by_first_name").Args([]interface{}{"Абр"})
	res2, err := conn.Do(req2).Get()
	fmt.Println(err)
	var userList []model.TarantoolUser
	for _, r := range res2 {
		list := r.([]interface{})
		for _, rr := range list {
			rrr := rr.([]interface{})
			user := model.TarantoolUser{
				UserID:         "",
				FirstName:      rrr[1].(string),
				SecondName:     rrr[2].(string),
				Username:       rrr[3].(string),
				HashedPassword: rrr[4].(string),
				Sex:            rrr[5].(string),
				Biography:      rrr[6].(string),
				City:           rrr[7].(string),
			}
			userList = append(userList, user)
		}
	}

	fmt.Println(len(userList))

	req3 := tarantool.NewCallRequest("search_by_first_second_name").Args([]interface{}{"Абр", "Юр"})
	res3, err := conn.Do(req3).Get()
	fmt.Println(err)
	fmt.Println(res3)

	var userList2 []model.TarantoolUser
	for _, r := range res3 {
		list := r.([]interface{})
		for _, rr := range list {
			rrr := rr.([]interface{})
			user := model.TarantoolUser{
				UserID:         "",
				FirstName:      rrr[1].(string),
				SecondName:     rrr[2].(string),
				Username:       rrr[3].(string),
				HashedPassword: rrr[4].(string),
				Sex:            rrr[5].(string),
				Biography:      rrr[6].(string),
				City:           rrr[7].(string),
			}
			userList2 = append(userList2, user)
		}
	}

	fmt.Println(len(userList2))

	req4 := tarantool.NewCallRequest("search_by_first_second_name_with_offset").Args([]interface{}{"Абр", "Юр", 10, 0})
	res4, err := conn.Do(req4).Get()

	var userList3 []model.TarantoolUser
	for _, r := range res4 {
		list := r.([]interface{})
		for _, rr := range list {
			rrr := rr.([]interface{})
			user := model.TarantoolUser{
				UserID:         "",
				FirstName:      rrr[1].(string),
				SecondName:     rrr[2].(string),
				Username:       rrr[3].(string),
				HashedPassword: rrr[4].(string),
				Sex:            rrr[5].(string),
				Biography:      rrr[6].(string),
				City:           rrr[7].(string),
			}
			userList3 = append(userList3, user)
		}
	}
	fmt.Println(err)
	fmt.Println(len(userList3))

	// Select by primary key
	// data, err := conn.Do(
	// 	tarantool.NewSelectRequest("users").
	// 		Limit(10).
	// 		Iterator(tarantool.IterEq).
	// 		Key([]interface{}{uint(1)}),
	// ).Get()
	// if err != nil {
	// 	fmt.Println("Got an error:", err)
	// }
	// fmt.Println("Tuple selected by the primary key value:", data)
	//
	// // Select by primary key
	// var typ []model.TarantoolUser
	// req := tarantool.NewCallRequest("search_by_first_name").Args([]interface{}{"first_name"})
	// err = conn.Do(req).GetTyped(&typ)
	// if err != nil {
	// 	fmt.Println("Got an error:", err)
	// }
	// fmt.Println("Tuple selected by the primary key value:", typ)
	//
	// var typ2 []model.TarantoolUser
	// err = conn.Do(
	// 	tarantool.NewSelectRequest("users").
	// 		Limit(10).
	// 		Iterator(tarantool.IterEq).
	// 		Key([]interface{}{uint(1)}),
	// ).GetTyped(&typ2)
	// if err != nil {
	// 	fmt.Println("Got an error:", err)
	// }
	// fmt.Println("SELECTED 1:", typ2)
	//
	// // var typ [][]model.TarantoolUser
	// req = tarantool.NewCallRequest("search_by_first_second_name").Args([]interface{}{"fir", "sec", 10, 0})
	// data, err = conn.Do(req).Get()
	// if err != nil {
	// 	fmt.Println("Got an error:", err)
	// }
	// fmt.Println("SELECTED 2:", data)
	//
	// var tupless []model.TarantoolUser
	// err = conn.SelectTyped("users", "primary", 0, 10, tarantool.IterLe,
	// 	[]interface{}{uint(1)},
	// 	&tupless)
	//
	// fmt.Println(tupless)
}
