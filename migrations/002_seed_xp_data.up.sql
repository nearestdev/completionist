INSERT INTO xp_actions (action_key, xp_amount) VALUES 
('list_completion', 50),
('review', 20),
('comment', 5)
ON CONFLICT (action_key) DO UPDATE SET xp_amount = EXCLUDED.xp_amount;


INSERT INTO level_definitions (level, xp_required)
SELECT generate_series, (generate_series - 1) * 100
FROM generate_series(1, 100)
ON CONFLICT (level) DO NOTHING;