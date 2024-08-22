# wb-l0
Демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе.
## Развертывание системы
Система полностью контейнеризирована для использования в Docker.
```
docker compose up --build
```
Команда собирает и запускает базу данных, брокер сообщений, сервер и скрипт для публикации сообщений.

## Настройка системы
Для настройки системы необходимо использовать JSON-конфиги, подобные тем, которые расположены в [`server/configs`](server/configs) и [`publisher/configs`](publisher/configs).
Для передачи пути к конфигу используется флаг `-c`.
```
cd server
go run cmd/main.go -c configs/local.dev.json
```

## Стресс-тестирование сервера
База данных и брокер сообщений контейнеризированы в Docker (WSL2). Рассматриваются два варианта тестирования с одинаковыми параметрами для `go-wrk`, которые оптимальны для получения наилучших безошибочных результатов на **конкретной программно-аппаратной платформе**.
### Тестирование локального сервера
```
go-wrk -c 100 http://localhost:3000/?id=116d658a-32b7-4c55-b944-832c815b4a7a

Running 10s test @ http://localhost:3000/?id=116d658a-32b7-4c55-b944-832c815b4a7a
  100 goroutine(s) running concurrently
1668463 requests in 9.952742523s, 1.47GB read
Requests/sec:           167638.52
Transfer/sec:           151.08MB
Overall Requests/sec:   166579.87
Overall Transfer/sec:   150.13MB
Fastest Request:        0s
Avg Req Time:           596µs
Slowest Request:        26.488ms
Number of Errors:       0
10%:                    0s
50%:                    0s
75%:                    0s
99%:                    0s
99.9%:                  0s
99.9999%:               0s
99.99999%:              0s
stddev:                 911µs
```

### Тестирование сервера в Docker (WSL2)
```
go-wrk -c 100 http://localhost:3000/?id=116d658a-32b7-4c55-b944-832c815b4a7a

Running 10s test @ http://localhost:3000/?id=116d658a-32b7-4c55-b944-832c815b4a7a
  100 goroutine(s) running concurrently
338466 requests in 9.981147777s, 305.03MB read
Requests/sec:           33910.53
Transfer/sec:           30.56MB
Overall Requests/sec:   33738.58
Overall Transfer/sec:   30.41MB
Fastest Request:        0s
Avg Req Time:           2.948ms
Slowest Request:        115.939ms
Number of Errors:       0
10%:                    0s
50%:                    506µs
75%:                    512µs
99%:                    516µs
99.9%:                  516µs
99.9999%:               516µs
99.99999%:              516µs
stddev:                 2.068ms
```
Для получения лучших результатов можно заменить стандартный HTTP-сервер из пакета `net/http` на сервер из [`github.com/valyala/fasthttp`](https://pkg.go.dev/github.com/valyala/fasthttp).