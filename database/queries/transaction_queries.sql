-- +Query untuk Mengambil Informasi Transfer dan Detail Akun Pengirim dan Penerima
SELECT
    t.id AS transfer_id,
    t.amount,
    t.transfer_type,
    t.fee,
    t.status AS transfer_status,
    t.created_at AS transfer_created_at,
    t.completed_at AS transfer_completed_at,
    sa.id AS sender_account_id,
    sa.account_number AS sender_account_number,
    sa.balance AS sender_balance,
    ra.id AS receiver_account_id,
    ra.account_number AS receiver_account_number,
    ra.balance AS receiver_balance
FROM
    transfers t
JOIN
    accounts sa ON t.sender_account_id = sa.id
JOIN
    accounts ra ON t.receiver_account_id = ra.id;



-- +Query untuk Mengambil Semua Transfer yang Melibatkan Akun Tertentu
SELECT
    t.id AS transfer_id,
    t.amount,
    t.transfer_type,
    t.fee,
    t.status AS transfer_status,
    t.created_at AS transfer_created_at,
    t.completed_at AS transfer_completed_at,
    sa.account_number AS sender_account_number,
    ra.account_number AS receiver_account_number
FROM
    transfers t
JOIN
    accounts sa ON t.sender_account_id = sa.id
JOIN
    accounts ra ON t.receiver_account_id = ra.id
WHERE
    sa.id =91206  OR ra.id = 34621;  -- Ganti $1 dengan ID akun yang dicari;  


-- +Query untuk Mengambil Log Transaksi dan Detail Transfer
SELECT
    l.id AS log_id,
    l.log_message,
    l.created_at AS log_created_at,
    t.id AS transfer_id,
    t.amount,
    t.transfer_type,
    t.fee,
    t.status AS transfer_status,
    t.created_at AS transfer_created_at
FROM
    transaction_logs l
JOIN
    transfers t ON l.transfer_id = t.id;

-- Query untuk Mengambil Akun dengan Transfer yang Tertunda
SELECT
    a.id AS account_id,
    a.account_number,
    a.balance,
    t.id AS transfer_id,
    t.amount,
    t.transfer_type,
    t.fee,
    t.status AS transfer_status,
    t.created_at AS transfer_created_at
FROM
    accounts a
JOIN
    transfers t ON a.id = t.sender_account_id OR a.id = t.receiver_account_id
WHERE
    t.status = 'pending';

--  Query untuk Mengambil Semua Akun dan Status Aktivitas Pengguna
SELECT
    a.id AS account_id,
    a.account_number,
    a.balance,
    a.account_type,
    u.id AS user_id,
    u.username,
    u.email,
    u.is_active
FROM
    accounts a
JOIN
    users u ON a.user_id = u.id;

-- untuk mengecek pengguna yang sudah terdaftar di tabel users namun belum memiliki akun di tabel accounts, Anda dapat menggunakan query SQL berikut:
SELECT
    u.id AS user_id,
    u.username,
    u.email,
    u.phonenumber
FROM
    users u
LEFT JOIN
    accounts a ON u.id = a.user_id
WHERE
    a.user_id IS NULL;

-- MENGABUNGNA TABEL USERS DAN ACCOUNTS
SELECT
    u.id AS user_id,
    u.username,
    u.email,
    u.phonenumber,
    u.password,
    u.is_active,
    u.created_at AS user_created_at,
    u.created_by,
    u.modified_at,
    u.modified_by,
    u.activated_by,
    u.activated_at,
    a.id AS account_id,
    a.account_number,
    a.balance,
    a.account_type,
    a.created_at AS account_created_at
FROM
    public.users u
LEFT JOIN
    public.accounts a ON u.id = a.user_id
ORDER BY
    u.id ASC;
