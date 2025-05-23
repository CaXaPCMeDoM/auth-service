### "Продовый" запуск:
Если вы находитесь в директории проекта:
```
docker-compose up --build -d 
```

### Локальный запуск
Это не предусмотрено конфигурацией, поэтому придется ввести некоторые изменение:
1. Во-первых, добавим в compose связь базы с внешним миром:
в compose нужно прокинуть порты для базы(для примера возьмем 5432):
```
port:
    - 5432:5432
```

2. В .env изменим PG_URL, заменив auth-db-postgres на localhost
```
PG_URL=postgres://caxap:1234@localhost:5432/db?sslmode=disable
```
Теперь можем запускать стандартно Go приложение из main (./cmd/app/main.go)


### Ручки
Так как swagger я не добавил(не было в тз), то вот ручки моего приложения:
1. "Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса":
```
/v1/auth/tokens?user_id=
```
2. "Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов":
```
/v1/auth/refresh
```
с примерным Body:
```
{
    "refresh": "11e71771-951f-4844-962d-b9f66533fb47",
    "access": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6Ijo6MSIsInN1YiI6ImJmNWQ3YmFlLWJiZjUtNDUxNC1hNTBlLWY1YTJmMzM5N2RmMSIsImV4cCI6MTc0NjM0NDY4NywiaWF0IjoxNzQ2MzQ0MDg3LCJqdGkiOiJlNDM0NDI2ZS01YTE3LTRmZmQtOWI2My04MWQ2NDgwZjZjZjgifQ.lGmpg-lS1ATCg7aTcPn1mdPfz8dB0PbPhcz_9DKgXD3-j5Eo9UsEoL6qL9BDXOYFSyIMXEd26qj_aw9thhyorA"
}
```

### ВАЖНО!
Так как добавление юзера не было описано в задаче, то я добавил тестовые данные при миграции. Файл ```./migrations/00000003_create_test_user.up.sql```
Там указаны мои креды, на которых я тестировал своё приложение.

- Также специально создал почту, с которой отправляются варнинги о смене ip. Её данные указаны в .env
  ![image](https://github.com/user-attachments/assets/dc7bcc52-db88-4550-979e-193ebc27c71a)

### Комментарии по реализации:
1. GPT и другие LLM не использовались. Обращался только к документации. Некоторые ресурсы оставил ссылкой в качестве комментариев
2. Считаем, что в базе уже лежат пользователи с соответствующими емайлами.
3. Предусмотрен случай разницы времени на инфре и в приложении
4. Для попыток повторного использования данные о сессии удаляются после выдачи
5. Для упрощения отправка возможна только с одной почтой
6. .env не добавлял в gitignore для примера данных.
7. Добавил обертку для транзакций, чтобы не мешать сервисный слой с инфраструктурным. Просто сохраняется транзакционный контекст и переопределяются методы обращения к бд для "жонглирования" между двумя возможными состояниями.
