docker-compose up -d

Write-Host "Waiting for database to start..."
Start-Sleep -Seconds 3

Write-Host "Running migrations..."
Get-Content database\migrations\000001_init_users.up.sql | docker exec -i finance-postgres psql -U postgres -d finance
Get-Content database\migrations\000002_create_accounts.up.sql | docker exec -i finance-postgres psql -U postgres -d finance
Get-Content database\migrations\000003_create_categories.up.sql | docker exec -i finance-postgres psql -U postgres -d finance
Get-Content database\migrations\000004_create_transactions.up.sql | docker exec -i finance-postgres psql -U postgres -d finance

Write-Host "Database is ready!"
