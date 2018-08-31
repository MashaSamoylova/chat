package main

import (
    "fmt"
    "net"
)

type User struct {
    username string
    password string
    channel chan string
    connection net.Conn
    friend *User
}

var users [200]User
var i int

func getUsernames() string {
    var k int
    var usernames string
    for k < i {
        usernames += users[k].username + " "
        k++
    }
    usernames += "\n"
    return usernames
}

func getUser(username string) *User {
    var k int
    k = 0
    for k < i {
        if users[k].username == username {
            return &users[k]
        }
        k++
    }
    return nil
}

func listenMessages(user *User) {
    for {
        msg := <-user.channel
        _, err := user.connection.Write([]byte(msg))
        if err != nil {
            return
        }
    }
}

func openRoom(user *User) {
    user.connection.Write([]byte("Чтобы выйти из чат комнаты напишите 'q'\n"))
    for {
        var buf [512]byte

        n, err := user.connection.Read(buf[0:])
        if err != nil {
            return
        }
        msg := string(buf[:n])
        if msg == "q\n" {
            user.connection.Write([]byte("Выходим из чат комнаты с " + user.friend.username + "\n"))
            user.friend = nil
            return
        }
        msg = user.username + ": " + msg
        user.friend.channel <- msg
    }
}

func handleUser(user *User) {
    var buf [512]byte
    _, err := user.connection.Write([]byte("Привет, как тебя зовут?\n"))

    if err != nil {
        return
    }
    n, err := user.connection.Read(buf[0:])
    if err != nil {
        return
    }
    user.username = string(buf[:n-1])
    user.connection.Write([]byte("Привет, " + user.username + "!\n"))

    go listenMessages(user)

    for {
        _, err = user.connection.Write([]byte("С кем открыть чат-комнату?\n"))
        if err != nil {
            return
        }

        _, err = user.connection.Write([]byte(getUsernames()))
        if err != nil {
            return
        }

        n, err = user.connection.Read(buf[0:])
        if err != nil {
            return
        }

        recipier_name := string(buf[:n-1])
        recipier := getUser(recipier_name)

        if recipier == nil {
            user.connection.Write([]byte("Здесь нет " + recipier_name + "\n"))
            continue
        }
        user.friend = recipier
        openRoom(user)
    }
}

func main() {
    // слушать порт
    ln, err := net.Listen("tcp", ":1212")
    if err != nil {
        fmt.Println(err)
        return
    }

    for {
        new_connection, err := ln.Accept()

        fmt.Println("new connection")
        if err != nil {
            fmt.Println(err)
            continue
        }

        users[i] = User{"","", make(chan string), new_connection, nil}
        go handleUser(&users[i])
        i++
    }
}
