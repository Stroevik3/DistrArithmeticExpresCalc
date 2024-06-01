# DistrArithmeticExpresCalc
Проект "Распределенный вычислитель арифметических выражений" в рамках обучение в яндекс практикум.
Он позволяет расчитать результат поступившего арифметического выражения в асинхронном режиме.
Через POST запрос передается выражение на вычисление и его id.
Через GET запрос можно получить результат обработки конкретного выражения по id или списко всех выражений.

# Запуск проекта
1) Скачать проект на свой ПК
2) Подтянуть используемые сторонние пакеты.
3) Открыть консоль и запустить оркестратор
4) Открыть вторую консоль и запустить агента
5) Отурыть третью коносль, через которую добавить выражение на вычисление и получить результат.

# Примечание
Версия go 1.22.0 +
Все примеры указаны для Windows 10+

# Описание
## Внутри проект состоит из двух модулей:
- Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. Далее будем называть его "оркестратором".
- Вычислитель, который может получить от "оркестратора" задачу, выполнить его и вернуть серверу результат. Далее будем называть его "агентом".

## Оркестратор

### POST запрос на добавление выражения и его id. 
#### Пример:
```
curl --header "Content-Type: application/json" --request POST --data {\"id\":1,\"expression\":\"2*3+1\"} http://localhost:8080/api/v1/calculate
```

## GET запрос для получения конкретного выражения. 
### Пример:
```
curl http://localhost:8080/api/v1/expressions/1
```

### GET запрос для получения списка выражений. 
#### Пример:
```
curl http://localhost:8080/api/v1/expressions
```



# Пример запросов на успешное добавлеиние выражения
```
curl --header "Content-Type: application/json" --request POST --data {\"id\":1,\"expression\":\"21*(7+5)\"} http://localhost:8080/api/v1/calculate
curl --header "Content-Type: application/json" --request POST --data {\"id\":3,\"expression\":\"-1+2\"} http://localhost:8080/api/v1/calculate
curl --header "Content-Type: application/json" --request POST --data {\"id\":4,\"expression\":\"32*0.5+16/32-(34+16)-100\"} http://localhost:8080/api/v1/calculate
curl --header "Content-Type: application/json" --request POST --data {\"id\":4,\"expression\":\"(1+(-60+4+15*0.2))*(1+0.5*3)\"} http://localhost:8080/api/v1/calculate
```

# Пример запросов на добавлеиние выражения который приведт к ошибке
## Ошибка деления на 0
```
curl --header "Content-Type: application/json" --request POST --data {\"id\":1,\"expression\":\"12/0\"} http://localhost:8080/api/v1/calculate
curl --header "Content-Type: application/json" --request POST --data {\"id\":2,\"expression\":\"12/(5-5)\"} http://localhost:8080/api/v1/calculate
```


