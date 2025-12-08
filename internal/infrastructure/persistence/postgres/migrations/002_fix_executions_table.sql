-- Drop the problematic tables first
DROP TABLE IF EXISTS execution_node_data CASCADE;
DROP TABLE IF EXISTS executions CASCADE;

-- Recreate executions table without partitioning for now
CREATE TABLE executions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workflow_id UUID NOT NULL REFERENCES workflows(id),
    workflow_version INT NOT NULL,
    status VARCHAR(50) NOT NULL, -- waiting, running, success, error, cancelled
    mode VARCHAR(50) NOT NULL, -- manual, trigger, webhook, schedule
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    execution_time_ms INT,
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    error_message TEXT,
    error_node VARCHAR(255),
    retry_of UUID,
    retry_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add self-referential foreign key after table creation
ALTER TABLE executions ADD CONSTRAINT fk_retry_of FOREIGN KEY (retry_of) REFERENCES executions(id);

-- Node execution data
CREATE TABLE execution_node_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    execution_id UUID NOT NULL REFERENCES executions(id) ON DELETE CASCADE,
    node_id VARCHAR(255) NOT NULL,
    node_type VARCHAR(100) NOT NULL,
    node_name VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    input_data JSONB DEFAULT '{}',
    output_data JSONB DEFAULT '{}',
    error_message TEXT,
    execution_time_ms INT,
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    retry_count INT DEFAULT 0
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_executions_workflow_status ON executions(workflow_id, status);
CREATE INDEX IF NOT EXISTS idx_executions_created_at ON executions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_node_execution_data_execution ON execution_node_data(execution_id);
