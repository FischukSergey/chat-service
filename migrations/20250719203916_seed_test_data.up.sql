-- Предзаполнение тестовыми данными
INSERT INTO chats (id, client_id, created_at, updated_at)
VALUES (
    '3b1ecbf8-79ec-11ed-adc7-461e464ebed8',
    'bbc3fa26-2961-400b-beec-6fc56d509c36',
    now(),
    now()
  ) ON CONFLICT (id) DO NOTHING;
INSERT INTO problems (id, manager_id, created_at, updated_at, chat_id)
VALUES (
    '7503cc2e-79ec-11ed-b705-461e464ebed8',
    null,
    now(),
    now(),
    '3b1ecbf8-79ec-11ed-adc7-461e464ebed8'
  ) ON CONFLICT (id) DO NOTHING;
DO $$
DECLARE author_id uuid;
BEGIN FOR i IN 1..35 LOOP IF i % 4 = 0 THEN author_id = '5cdc09b6-79ee-11ed-a98d-461e464ebed8';
ELSE author_id = 'bbc3fa26-2961-400b-beec-6fc56d509c36';
END IF;
INSERT INTO messages (
    id,
    author_id,
    is_visible_for_client,
    is_visible_for_manager,
    body,
    is_blocked,
    is_service,
    created_at,
    chat_id,
    problem_id
  )
VALUES (
    gen_random_uuid(),
    author_id,
    true,
    true,
    'message #' || i,
    false,
    false,
    clock_timestamp(),
    '3b1ecbf8-79ec-11ed-adc7-461e464ebed8',
    '7503cc2e-79ec-11ed-b705-461e464ebed8'
  );
END LOOP;
END;
$$;