INSERT INTO users (id, email, password_hash)
SELECT '00000000-0000-4000-8000-000000000001', 'dev@example.com', '$2a$12$VaAqIVs3WTyN4YS1VXSLWOmvyz4ibrqjMlZemztTL.rc5hTOml0M6'
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'dev@example.com');
