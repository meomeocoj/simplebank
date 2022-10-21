ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";
-- ALTER TABLE IF EXISTS "accounts" DROP INDEX IF EXISTS "accounts_owner_currency_idx";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
DROP TABLE users;