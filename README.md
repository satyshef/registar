# Registar

Регистратор Телеграм аккаунтов

Для сборки Docker образа используется скрипт build. 
Образ основан на Alpine. 
Точка входа в контейнере скрипт entrypoint.sh
    CONFIG - Название конфигурационного файла
    NUMBER - Номер телефона. Если не указан то пробует использовать SERVICE
    SERVICE - Название СМС сервиса

Директория приложения /app. Директория с профилями по умолчанию /app/profiles
Переменные окружения:


 p - Путь к директории с профилями ("./profiles")
 d - Путь к директории конфигураций ("./data/config")
 r - Путь к дериктории с сервисами ("./data/service")
 n - Номер телефона. Если указан СМС сервис не используется
 c - Название конфигурационного файла ("bot.toml")
 a - Количество аккаунтов (1)
 s - Имя СМС сервиса. Поддерживаются:
      365sms
      sms3t
      5sim
      sms-acktiwator

    