CREATE TABLE action_instances (
  id SERIAL PRIMARY KEY,
  trigger_instance_id INTEGER,
  output_data jsonb,
  status INTEGER,
  created_at TIMESTAMP,
  upated_at TIMESTAMP
);