-- migrations/000001_create_data_table.up.sql
-- Создание таблицы
CREATE TABLE t_data (
    id SERIAL PRIMARY KEY,
    s_full VARCHAR(1000) NOT NULL,
    s_short VARCHAR(100) NOT NULL
);

CREATE INDEX idx_data_full ON t_data(s_full);
CREATE INDEX idx_data_short ON t_data(s_short);