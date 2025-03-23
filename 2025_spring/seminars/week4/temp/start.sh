# Проверьте, корректно ли у вас установлен postgres в зависимости от вашей ОС.

# Настройка БД
initdb -D ./data

# Запуск
pg_ctl -D ./data -l logfile start

# Запуск синхронный
postgres -D ./data
