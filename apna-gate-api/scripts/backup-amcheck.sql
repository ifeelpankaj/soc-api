-- Run once as the database owner/admin, outside the webhook:
-- psql --set=ON_ERROR_STOP=1 --set=backup_role=app_user --file=scripts/backup-amcheck.sql
-- PGDATABASE must identify the same database as the API's DB_NAME.
CREATE EXTENSION IF NOT EXISTS amcheck WITH SCHEMA public;
GRANT USAGE ON SCHEMA public TO :"backup_role";
GRANT EXECUTE ON FUNCTION public.bt_index_check(regclass, boolean) TO :"backup_role";
GRANT EXECUTE ON FUNCTION public.verify_heapam(regclass, boolean, boolean, text, bigint, bigint) TO :"backup_role";
