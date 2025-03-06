INSERT INTO commissions
    (
     name, description, budget, currency, location, deadline, status, winner, entries
     ) VALUES (
               'test', 'test', 100, 'USD', 'test', '2020-01-01', 'open', '', ''
               );
SELECT LAST_INSERT_ID();


