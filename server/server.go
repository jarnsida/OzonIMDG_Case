package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/jarsida/OzonIMDG_Case/service"
)

type Server struct {
	listener         net.Listener
	quit             chan struct{}
	exited           chan struct{}
	db               memoryDB
	connections      map[int]net.Conn
	connCloseTimeout time.Duration
}

func NewServer() *Server {
	//cfg := config.Get()

	laddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		log.Fatal("Неудалось создать подключение ", err.Error())
	}

	l, err := net.ListenTCP("tcp", laddr) //!заменить порт на порт из конфига
	if err != nil {
		log.Fatal("Неудалось создать подключение ", err.Error())
	}

	srv := &Server{
		listener:         l,
		quit:             make(chan struct{}),
		exited:           make(chan struct{}),
		db:               newDB(),
		connections:      map[int]net.Conn{},
		connCloseTimeout: 4 * time.Second,
	}
	go srv.serve()
	return srv
}

func (srv *Server) serve() {
	var id int //Счётчик пользователей

	fmt.Println("Ожидаю подключение клиентов")
	for {
		select {
		case <-srv.quit:
			fmt.Println("Завершение работы сервера БД")
			err := srv.listener.Close()
			if err != nil {
				fmt.Println("Не удалось отключить клиента", err.Error())
			}
			if len(srv.connections) > 0 {
				srv.warnConnections(srv.connCloseTimeout)
				<-time.After(srv.connCloseTimeout)
				srv.closeConnections()
			}
			close(srv.exited)
			return
		default:
			tcpListener := srv.listener.(*net.TCPListener)
			err := tcpListener.SetDeadline(time.Now().Add(2 * time.Second))
			if err != nil {
				fmt.Println("Не удалось установить listener deadline", err.Error())
			}

			conn, err := tcpListener.Accept()
			if oppErr, ok := err.(*net.OpError); ok && oppErr.Timeout() {
				continue
			}
			if err != nil {
				fmt.Println("Не удалось создать соединение", err.Error())
			}

			write(conn, "Добро пожаловть в OzonIMDB server")
			srv.connections[id] = conn
			go func(connID int) {
				fmt.Println("Клиент с id", connID, "подключён")
				srv.handleConn(conn)
				delete(srv.connections, connID)
				fmt.Println("Клиент с id", connID, "отключён")
			}(id)
			id++
		}
	}
}

func write(conn net.Conn, s string) {
	_, err := fmt.Fprintf(conn, "%s\n-> ", s)
	if err != nil {
		log.Fatal(err)
	}
}

func (srv *Server) handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	mem := service.NewMonitor()
	for scanner.Scan() {
		l := strings.ToLower(strings.TrimSpace(scanner.Text()))
		values := strings.Split(l, " ")

		switch {

		// Сервис set записывает новую пару ключ:значение
		case len(values) == 3 && values[0] == "set":
			srv.db.set(values[1], values[2])
			write(conn, "OK")

		// Сервис get выводит значение записи с указанным ключом
		case len(values) == 2 && values[0] == "get":
			k := values[1]
			val, found := srv.db.get(k)
			if !found {
				write(conn, fmt.Sprintf("key %s not found", k))
			} else {
				write(conn, val)
			}

		// Сервис delete удлаяет запись с указанным ключом
		case len(values) == 2 && values[0] == "delete":
			srv.db.delete(values[1])
			write(conn, "OK")

		// Сервис count выводит количество записей в базе
		case len(values) == 1 && values[0] == "count":
			k := srv.db.count()
			write(conn, strconv.Itoa(k))

		// Сервис memstats выводит информацию о состоянии памяти
		case len(values) == 1 && values[0] == "memstats":
			k := mem.Get()
			write(conn, k)
		// Сервис exit отключает пользователя
		case len(values) == 1 && values[0] == "exit":
			if err := conn.Close(); err != nil {
				fmt.Println("Невозможно завершить соединение", err.Error())
			}

		default:
			write(conn, fmt.Sprintf("UNKNOWN command: %s \n Supported commands: get, set, delete, count, exit", l))
		}
	}
}

func (srv *Server) warnConnections(timeout time.Duration) {
	for _, conn := range srv.connections {
		write(conn, fmt.Sprintf("Остановка сервера произойдёт через: %s", timeout.String()))
	}
}

func (srv *Server) closeConnections() {
	fmt.Println("Закрываю все соединения")
	for id, conn := range srv.connections {
		err := conn.Close()
		if err != nil {
			fmt.Println("Не могу завершить соединение с User id:", id)
		}
	}
}

func (srv *Server) Stop() {
	fmt.Println("Останавливаю сервер БД")
	close(srv.quit)
	<-srv.exited
	fmt.Println("Сохраняю записи во внешний файл")
	srv.db.save()
	fmt.Println("Сервер БД успешно остановлен")
}
