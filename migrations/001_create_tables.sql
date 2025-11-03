-- Включаем расширение для UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================
-- ENUM типы
-- ============================================
CREATE TYPE bot_type AS ENUM ('AI', 'Backend', 'Frontend', 'Robot');
CREATE TYPE bot_lang AS ENUM ('Python', 'Go', 'N8N', 'PIX', 'JS', 'C', 'Other');
CREATE TYPE eff_status AS ENUM ('success', 'warning', 'error');

COMMENT ON TYPE bot_type IS 'Тип системы или робота: AI, Backend, Frontend, Robot';
COMMENT ON TYPE bot_lang IS 'Основной язык или технология, на которой реализован бот/система';
COMMENT ON TYPE eff_status IS 'Статус выполнения робота или системы: success / warning / error';

-- ============================================
-- Таблица владельцев (пользователей системы)
-- ============================================
CREATE TABLE IF NOT EXISTS owners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE owners IS 'Владельцы ботов - пользователи системы';
COMMENT ON COLUMN owners.id IS 'Уникальный идентификатор владельца (UUID)';
COMMENT ON COLUMN owners.full_name IS 'Полное имя владельца';
COMMENT ON COLUMN owners.is_active IS 'Активен ли владелец';
COMMENT ON COLUMN owners.created_at IS 'Дата и время создания записи';

-- ============================================
-- Таблица ботов/роботов
-- ============================================
CREATE TABLE IF NOT EXISTS bots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    bot_type bot_type NOT NULL,
    language bot_lang NOT NULL,
    description TEXT,
    tags TEXT[],
    owner_id UUID REFERENCES owners(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE bots IS 'Боты и роботы, которые пишут логи';
COMMENT ON COLUMN bots.id IS 'Уникальный идентификатор бота (UUID)';
COMMENT ON COLUMN bots.code IS 'Уникальный код бота для идентификации';
COMMENT ON COLUMN bots.name IS 'Название бота';
COMMENT ON COLUMN bots.bot_type IS 'Тип бота: AI, Backend, Frontend, Robot';
COMMENT ON COLUMN bots.language IS 'Язык или технология реализации';
COMMENT ON COLUMN bots.description IS 'Описание бота';
COMMENT ON COLUMN bots.tags IS 'Теги для категоризации';
COMMENT ON COLUMN bots.owner_id IS 'Владелец бота (может быть NULL если владелец удалён)';
COMMENT ON COLUMN bots.is_active IS 'Активен ли бот';
COMMENT ON COLUMN bots.created_at IS 'Дата и время создания';
COMMENT ON COLUMN bots.updated_at IS 'Дата и время последнего обновления';

-- ============================================
-- Таблица токенов (для аутентификации API)
-- ============================================
-- Каждый бот может иметь несколько токенов
-- При удалении бота удаляются его токены (CASCADE)
CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Токен = UUID, используется как ключ API
    bot_id UUID REFERENCES bots(id) ON DELETE CASCADE, -- Бот, которому принадлежит токен (NULL для админских токенов)
    name VARCHAR(100) NOT NULL,                    -- Название токена ("Production", "Test Instance")
    is_active BOOLEAN DEFAULT true,                -- Активен ли токен (можно деактивировать)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), -- Дата и время создания токена
    is_admin BOOLEAN NOT NULL DEFAULT false        -- Флаг, указывающий, является ли токен административным
);
COMMENT ON TABLE tokens IS 'Таблица для хранения API токенов, используемых для аутентификации ботов.';
COMMENT ON COLUMN tokens.id IS 'Уникальный идентификатор токена (UUID), который также является самим токеном.';
COMMENT ON COLUMN tokens.bot_id IS 'Ссылка на бота, которому принадлежит данный токен. NULL для админских токенов. При удалении бота, все его токены также удаляются.';
COMMENT ON COLUMN tokens.name IS 'Человекочитаемое название токена, например, "Токен для продакшн-сервера".';
COMMENT ON COLUMN tokens.is_active IS 'Флаг активности токена. Неактивные токены не могут быть использованы.';
COMMENT ON COLUMN tokens.created_at IS 'Дата и время создания токена.';
COMMENT ON COLUMN tokens.is_admin IS 'Флаг, указывающий, обладает ли токен административными привилегиями.';

-- ============================================
-- Таблица логов
-- ============================================
CREATE TABLE IF NOT EXISTS logs (
    id BIGSERIAL PRIMARY KEY,
    bot_id UUID REFERENCES bots(id) ON DELETE SET NULL,
    status eff_status DEFAULT 'success',
    msg TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE logs IS 'Логи от ботов. Сохраняются даже после удаления бота';
COMMENT ON COLUMN logs.id IS 'Уникальный идентификатор лога (автоинкремент)';
COMMENT ON COLUMN logs.bot_id IS 'Бот, от которого пришёл лог (может быть NULL если бот удалён)';
COMMENT ON COLUMN logs.status IS 'Статус лога: success, warning, error';
COMMENT ON COLUMN logs.msg IS 'Текст сообщения лога';
COMMENT ON COLUMN logs.created_at IS 'Дата и время создания лога';

-- ============================================
-- Таблица эффективных запусков
-- ============================================
CREATE TABLE IF NOT EXISTS eff_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bot_id UUID NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    period_from TIMESTAMP WITH TIME ZONE,
    period_to TIMESTAMP WITH TIME ZONE,
    status eff_status NOT NULL DEFAULT 'success',
    host TEXT,
    extra JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT eff_runs_check CHECK (period_to > period_from)
);

COMMENT ON TABLE eff_runs IS 'Информация о запусках ботов за определённый период времени';
COMMENT ON COLUMN eff_runs.id IS 'Уникальный идентификатор запуска (UUID)';
COMMENT ON COLUMN eff_runs.bot_id IS 'Бот, который выполнялся';
COMMENT ON COLUMN eff_runs.period_from IS 'Начало периода выполнения';
COMMENT ON COLUMN eff_runs.period_to IS 'Конец периода выполнения';
COMMENT ON COLUMN eff_runs.status IS 'Результат выполнения: success, warning, error';
COMMENT ON COLUMN eff_runs.host IS 'Хост или сервер, на котором выполнялся бот';
COMMENT ON COLUMN eff_runs.extra IS 'Дополнительные метаданные в формате JSON';
COMMENT ON COLUMN eff_runs.created_at IS 'Дата и время создания записи';

-- ============================================
-- Индексы для производительности
-- ============================================

-- Индексы для таблицы bots
CREATE INDEX idx_bots_code ON bots(code);
CREATE INDEX idx_bots_owner ON bots(owner_id);
CREATE INDEX idx_bots_type ON bots(bot_type);

-- Индексы для таблицы tokens
CREATE INDEX idx_tokens_bot ON tokens(bot_id);
CREATE INDEX idx_tokens_active ON tokens(is_active);

-- Индексы для таблицы logs
CREATE INDEX idx_logs_bot ON logs(bot_id);
CREATE INDEX idx_logs_status ON logs(status);
CREATE INDEX idx_logs_created ON logs(created_at DESC);

-- Индексы для таблицы eff_runs
CREATE INDEX idx_eff_runs_bot_period ON eff_runs(bot_id, period_to DESC);
CREATE INDEX idx_eff_runs_status ON eff_runs(status);
CREATE INDEX idx_eff_runs_extra_gin ON eff_runs USING gin(extra jsonb_path_ops);
