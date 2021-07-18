Simple IMDG service for Ozon Fintech internship

## :scroll: Задача    
### Написать реализацию простого in-memory key-value хранилища (аля redis).
Требования:
- [X] Возможность добавить, искать и удалять произвольный набор байт по ключу
- [X] консистентность при параллельных запросах к хранилищу
- [ ] соблюдение ограничения на размер базы (по объему, задаваемый из конфига)
- [ ] тестовый пример: клиент, пишущий и читающий из хранилища
____
Далее по возможности/желанию можно развить работу:
- [ ] сделать бенчмарк
- [ ] добавить поддержку TTL для ключей
- [X] добавить поддержку персистентности, чтобы хранилище переживало рестарт
- [ ] добавить поддержку типов
- [X] добавить поддержку репликации
    * etc    
    
## :clipboard: 1 этап решения

- [X] TCP(telnet) server launch
- [X] Users connect/act/disconect
- [X] Data consistency (Mutex)
- [X] Graceful Shutdown with data backup 
- [X] set/get/delete/count operations defined for [string]string
- [ ] соблюдение ограничения на размер базы (по объему, задаваемый из конфига)
- [ ] тестовый пример: клиент, пишущий и читающий из хранилища

## 2 этап решения


## Запуск работы программы
```
source run.sh PORT
````