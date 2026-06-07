DROP INDEX IF EXISTS idx_data_full;
DROP INDEX IF EXISTS idx_data_short;

CREATE UNIQUE INDEX IF NOT EXISTS idx_data_full ON t_data(s_full);
CREATE UNIQUE INDEX IF NOT EXISTS idx_data_short ON t_data(s_short); 