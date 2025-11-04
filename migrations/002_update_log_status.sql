-- Миграция: изменение статусов логов на стандартные уровни логирования
-- Дата: 2025-11-04

-- Шаг 1: Создаём новый ENUM тип для статусов логов (стандартные уровни логирования)
CREATE TYPE log_status AS ENUM ('Debug', 'Info', 'Warning', 'Error', 'Critical');
COMMENT ON TYPE log_status IS 'Уровень лога: Debug / Info / Warning / Error / Critical';

-- Шаг 2: Добавляем временную колонку с новым типом
ALTER TABLE logs ADD COLUMN status_new log_status;

-- Шаг 3: Копируем данные с преобразованием (старые значения -> новые)
UPDATE logs SET status_new = CASE 
    WHEN status::text = 'success' THEN 'Info'::log_status
    WHEN status::text = 'warning' THEN 'Warning'::log_status
    WHEN status::text = 'error' THEN 'Error'::log_status
    ELSE 'Info'::log_status  -- На случай если есть другие значения
END;

-- Шаг 4: Удаляем старую колонку
ALTER TABLE logs DROP COLUMN status;

-- Шаг 5: Переименовываем новую колонку
ALTER TABLE logs RENAME COLUMN status_new TO status;

-- Шаг 6: Устанавливаем NOT NULL и значение по умолчанию
ALTER TABLE logs ALTER COLUMN status SET NOT NULL;
ALTER TABLE logs ALTER COLUMN status SET DEFAULT 'Info'::log_status;

-- Шаг 7: Обновляем комментарий
COMMENT ON COLUMN logs.status IS 'Уровень лога: Debug (отладка), Info (информация), Warning (предупреждение), Error (ошибка), Critical (критическая ошибка)';

