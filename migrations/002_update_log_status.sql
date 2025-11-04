-- Миграция: изменение статусов логов на стандартные уровни логирования
-- Дата: 2025-11-04

CREATE TYPE log_status AS ENUM ('Debug', 'Info', 'Warning', 'Error', 'Critical');
COMMENT ON TYPE log_status IS 'Уровень лога: Debug / Info / Warning / Error / Critical';

ALTER TABLE logs ADD COLUMN status_new log_status;

UPDATE logs SET status_new = CASE 
    WHEN status::text = 'success' THEN 'Info'::log_status
    WHEN status::text = 'warning' THEN 'Warning'::log_status
    WHEN status::text = 'error' THEN 'Error'::log_status
    ELSE 'Info'::log_status  
END;

ALTER TABLE logs DROP COLUMN status;

ALTER TABLE logs RENAME COLUMN status_new TO status;

ALTER TABLE logs ALTER COLUMN status SET NOT NULL;
ALTER TABLE logs ALTER COLUMN status SET DEFAULT 'Info'::log_status;

COMMENT ON COLUMN logs.status IS 'Уровень лога: Debug (отладка), Info (информация), Warning (предупреждение), Error (ошибка), Critical (критическая ошибка)';

