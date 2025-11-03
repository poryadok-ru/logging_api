-- Скрипт для wrk с авторизацией
wrk.method = "GET"
wrk.headers["Authorization"] = "Bearer aaaaaaaa-0000-1111-2222-333333333333"
wrk.headers["Content-Type"] = "application/json"
