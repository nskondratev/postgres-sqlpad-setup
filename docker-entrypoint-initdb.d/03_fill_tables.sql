COPY users_evaluation_of_satisfaction FROM '/pg_init_data/users_evaluation_of_satisfaction.csv' DELIMITER ',' QUOTE '"' CSV HEADER WHERE request_id IS NOT NULL;

COPY support_tickets FROM '/pg_init_data/support_tickets.csv' DELIMITER ',' QUOTE '"' CSV HEADER;

COPY new_items_by_support_users FROM '/pg_init_data/new_items_by_support_users.csv' DELIMITER ',' QUOTE '"' CSV HEADER;
