Тестовый проект для MEDODS. 

Выплонен по [CLEAN архитектуре](https://github.com/alfssobsd/notes/blob/master/golang/arch/golang_arch_description.md) 

Для того, чтобы все работало достаточно запустить docker compose c вылидными файлами окружения:

Файл в корне .database-env c окружением БД:

POSTGRES_PASSWORD

POSTGRES_USER

Файл в директории executor с окружением сервера: 
|Переменная| Значение|
|---|---|
|SERVER_HTTP_PORT   | Порт сервера |
|SERVER_HTTP_HOST   | Хост сервера(локалхост) |
|DATABASE_HOST      | Хост БД, в случае докеризации имя в композе (postgres) |
|DATABASE_USER      | Юзер БД, как и в конфиге БД |
|DATABASE_PASSWORD  | Пароль БД, как и в конфиге |
|DATABASE_NAME      | Имя БД, по умолчанию postgres |
|MAILER_FROM        | Ящик-отправитель для конфигурации SMTP |
|MAILER_PASSWORD    | Пароль от SMTP |
|MAILER_SMTP_HOST   | Хост SMTP |
|MAILER_SMTP_PORT   | Порт SMTP |