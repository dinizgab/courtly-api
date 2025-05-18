-- +goose Up
-- +goose StatementBegin
INSERT INTO companies (id, name, address, phone, email, slug) VALUES
  ('11111111-1111-1111-1111-111111111111',
   'Arena Society 10',
   'Av. Central, 1000, Campina Grande-PB',
   '+55 83 99999-0001',
   'contato@arena10.com.br',
   'arena-10'),

  ('22222222-2222-2222-2222-222222222222',
   'Beach Tennis Pro',
   'Rua Praia, 200, João Pessoa-PB',
   '+55 83 98888-0002',
   'info@beachtennispro.com.br',
   'beach-tennis-pro');

------------------------------------------------------------
-- Usuários (proprietários)
INSERT INTO users (id, email, password_hash, company_id) VALUES
  ('aaaaaaa1-0000-4000-8000-000000000001',
   'owner@arena10.com.br',
   '$2y$12$hashdemoowner1',          -- hash “dummy” de senha
   '11111111-1111-1111-1111-111111111111'),

  ('aaaaaaa2-0000-4000-8000-000000000002',
   'owner@beachtennispro.com.br',
   '$2y$12$hashdemoowner2',
   '22222222-2222-2222-2222-222222222222');

------------------------------------------------------------
-- Quadras
INSERT INTO courts (id, company_id, name, sport_type, hourly_price, is_active) VALUES
  ('33333333-3333-3333-3333-333333333331',
   '11111111-1111-1111-1111-111111111111',
   'Quadra Society 1', 'society', 120.00, TRUE),

  ('33333333-3333-3333-3333-333333333332',
   '11111111-1111-1111-1111-111111111111',
   'Quadra Society 2', 'society', 100.00, TRUE),

  ('33333333-3333-3333-3333-333333333333',
   '22222222-2222-2222-2222-222222222222',
   'Quadra Beach A', 'beach_tennis', 90.00, TRUE);

------------------------------------------------------------
-- Reservas
INSERT INTO bookings (id, court_id, start_time, end_time, status,
                      guest_name, guest_phone, guest_email, verification_code)
VALUES
  ('44444444-4444-4444-4444-444444444441',
   '33333333-3333-3333-3333-333333333331',
   '2025-05-20 18:00:00-03', '2025-05-20 19:30:00-03',
   'confirmed',
   'João Silva', '+55 83 97777-1234', 'joao@example.com', 'A1B2C3'),

  ('44444444-4444-4444-4444-444444444442',
   '33333333-3333-3333-3333-333333333332',
   '2025-05-21 20:00:00-03', '2025-05-21 21:00:00-03',
   'pending',
   'Maria Souza', '+55 83 98888-4321', 'maria@example.com', 'D4E5F6');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM bookings  WHERE id IN ('44444444-4444-4444-4444-444444444441',
                                   '44444444-4444-4444-4444-444444444442');

DELETE FROM courts    WHERE id IN ('33333333-3333-3333-3333-333333333331',
                                   '33333333-3333-3333-3333-333333333332',
                                   '33333333-3333-3333-3333-333333333333');

DELETE FROM users     WHERE id IN ('aaaaaaa1-0000-4000-8000-000000000001',
                                   'aaaaaaa2-0000-4000-8000-000000000002');

DELETE FROM companies WHERE id IN ('11111111-1111-1111-1111-111111111111',
                                   '22222222-2222-2222-2222-222222222222');
-- +goose StatementEnd
