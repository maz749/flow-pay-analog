@echo off
echo Applying missing migrations to FlowPay database...
echo.

echo Copying migration files to container...
docker cp migrations/002_add_cancellation_instructions.sql flowpay-postgres:/tmp/002.sql
docker cp migrations/003_add_teams_and_roles.sql flowpay-postgres:/tmp/003.sql

echo.
echo Applying migration 002 (cancellation instructions)...
docker exec flowpay-postgres psql -U flowpay -d flowpay_db -f /tmp/002.sql

echo.
echo Applying migration 003 (teams and roles)...
docker exec flowpay-postgres psql -U flowpay -d flowpay_db -f /tmp/003.sql

echo.
echo Cleaning up temporary files...
docker exec flowpay-postgres rm /tmp/002.sql /tmp/003.sql

echo.
echo Restarting application...
docker restart flowpay-app

echo.
echo Done! Migrations applied successfully.
echo The application should now be ready at http://localhost:8080
pause
