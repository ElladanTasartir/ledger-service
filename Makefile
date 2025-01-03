.PHONY: run
run:
	@air -c ./build/.air.toml

.PHONY: migration-up
migration-up:
	@migrate -source file://db/migrations -database postgresql://ledger:pg_ledger@localhost:5432/ledger?sslmode=disable --verbose up

.PHONY: migration-down
migration-down:
	@migrate -source file://db/migrations -database postgresql://ledger:pg_ledger@localhost:5432/ledger?sslmode=disable --verbose down 1