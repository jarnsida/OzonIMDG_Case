package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/jarsida/OzonIMDG_Case/config"
	"github.com/jarsida/OzonIMDG_Case/service"
)

//Server type определяет тип данных TCP сервера
type Server struct {
	listener         net.Listener
	quit             chan struct{}
	exited           chan struct{}
	db               memoryDB
	ttl              timeToLive
	connections      map[int]net.Conn
	connCloseTimeout time.Duration
}

//NewServer запускает сервер в горутине
func NewServer() *Server {

	// Загрузка конфига
	cfg := config.Get()

	//Резервация порта
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":"+cfg.IMDBPort))
	if err != nil {
		log.Fatal("Неудалось создать подключение ", err.Error())
	}

	//Слушать порт
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal("Неудалось создать подключение ", err.Error())
	}

	srv := &Server{
		listener:         l,
		quit:             make(chan struct{}),
		exited:           make(chan struct{}),
		db:               newDB(),
		ttl:              newTTLdb(),
		connections:      map[int]net.Conn{},
		connCloseTimeout: time.Duration(cfg.ConnCloseTimeout) * time.Second,
	}
	//Запуск экземпляра в горутине
	go srv.serve()
	return srv
}

func (srv *Server) serve() {
	var id int //Счётчик пользователей

	fmt.Println("Ожидаю подключение клиентов")
	for {
		select {

		// Отключение сервера при получении сигнала Ctrl+C
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

		//Работа сервера
		default:
			tcpListener := srv.listener.(*net.TCPListener)
			err := tcpListener.SetDeadline(time.Now().Add(2 * time.Second))
			if err != nil {
				fmt.Println("Не удалось установить listener deadline", err.Error())
			}

			//Приём вызова соединения
			conn, err := tcpListener.AcceptTCP()
			if oppErr, ok := err.(*net.OpError); ok && oppErr.Timeout() {
				continue
			}
			if err != nil {
				fmt.Println("Не удалось создать соединение", err.Error())
			}

			conn.SetReadBuffer(512 * 1024)
			conn.SetWriteBuffer(2 * 1024 * 1024)

			write(conn, "Добро пожаловть в OzonIMDB server")
			srv.connections[id] = conn

			//Счётчик пользователей в горутине
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

//write описывает функцию формирования сообщение пользователю о результате запроса
func write(conn net.Conn, s string) {
	_, err := fmt.Fprintf(conn, "%s\n-> ", s)
	if err != nil {
		log.Fatal(err)
	}
}

//handleConn сеанс соединения клиента с базой
func (srv *Server) handleConn(conn net.Conn) {
	cfg := config.Get()
	scanner := bufio.NewScanner(conn)

	//Приём запроса пользователя
	for scanner.Scan() {
		l := strings.ToLower(strings.TrimSpace(scanner.Text()))
		values := strings.Split(l, " ")

		//Ограничение работы в случае привышения объёма

		if len(srv.db.items) > cfg.MaxMemory {
			write(conn, fmt.Sprintf("DB Size excides  %d \n Supported commands shorted: delete, count,\n memstats, clean, backup, exit", cfg.MaxMemory))
			switch {
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
				mem := service.NewMonitor()
				k := mem.Get()
				write(conn, k)

				// Сервис exit отключает пользователя
			case len(values) == 1 && values[0] == "exit":
				if err := conn.Close(); err != nil {
					fmt.Println("Невозможно завершить соединение", err.Error())
				}

			case len(values) == 1 && values[0] == "clean":
				srv.db.clean()
				write(conn, "База обнулена")

			case len(values) == 1 && values[0] == "backup":
				srv.db.backUp()
				write(conn, "База сохранена во внешний файл и очищена из памяти")
				srv.db.clean()
			default:
				write(conn, fmt.Sprintf("UNKNOWN command: %s \n Supported commands: delete, count,\n memstats, clean, backup, exit", l))

			}

		} else {

			switch {

			// Сервис set записывает новую пару ключ:значение
			case len(values) == 3 && values[0] == "set":
				srv.db.set(values[1], values[2])
				write(conn, "OK")

				// Сервис set записывает новую пару ключ:значение
			case len(values) == 4 && values[0] == "set+ttl":
				srv.db.set(values[1], values[2])
				srv.ttl.setTTL(values[1], values[3])
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
				mem := service.NewMonitor()
				k := mem.Get()
				write(conn, k)

				// Сервис exit отключает пользователя
			case len(values) == 1 && values[0] == "exit":
				if err := conn.Close(); err != nil {
					fmt.Println("Невозможно завершить соединение", err.Error())
				}

			case len(values) == 1 && values[0] == "clean":
				srv.db.clean()
				write(conn, "База обнулена")

			case len(values) == 1 && values[0] == "backup":
				srv.db.backUp()
				write(conn, "База сохранена во внешний файл и очищена из памяти")
				srv.db.clean()

			default:
				write(conn, fmt.Sprintf("UNKNOWN command: %s \n Supported commands: get, set, set+ttl, delete, count,\n memstats, clean, backup, exit", l))
			}
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

//Stop Graceful Shutdown with data backup
func (srv *Server) Stop() {
	fmt.Println("Останавливаю сервер БД")
	close(srv.quit)
	<-srv.exited
	fmt.Println("Сохраняю записи во внешний файл")
	srv.db.save()
	fmt.Println("Сервер БД успешно остановлен")
}
