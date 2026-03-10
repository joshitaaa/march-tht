# Christmas Gift Redemption Console App

A Golang console app for gift redemption. It supports:

- Lookup team by `staff_pass_id` from a CSV mapping file
- Check redemption eligibility at team level
- Create redemption record only once per team

Redemption data is stored in SQLite with this schema:
- `team_name` (unique)
- `redeemed_at` (epoch milliseconds)

## Assumptions / Design Notes

- Each `staff_pass_id` maps to exactly one `team_name`.
- Mapping is loaded once at startup.
- Team-level idempotency is enforced by SQLite unique constraint on `team_name`.
- Designed as a small interactive program for testing purposes instead of full blown HTTP REST service.
- Redemption data is stored in SQLite instead of CSV for eg. since it would be easier to impose the uniqueness of `team_name` to say that is has already been redeemed.

## Steps To Run Locally

1. Start the app:

```bash
go run .
```

2. Choose from menu:

- `1` Lookup team
- `2` Redeem gift
- `3` Exit

3. Example redeem flow:

```bash
Choose [1-3]: 2
Enter staff_pass_id: S1234567A

================================
REDEEM RESULT
--------------------------------
Status: redeemed
Team: Platform Team
RedeemedAt: 1736500000000
================================
Press Enter to continue...
```

## Run Tests

```bash
go test ./... -v -cover
```

## Console Output Samples

### Redeem success

```text
================================
REDEEM RESULT
--------------------------------
Status: redeemed
Team: Platform Team
RedeemedAt: 1736500000000
================================
```

Already redeemed output:

```text
================================
REDEEM RESULT
--------------------------------
Status: already_redeemed
Team: Platform Team
================================
```

Not found output:

```text
================================
REDEEM RESULT
--------------------------------
Status: staff_not_found
================================
```

### Lookup found

```text
================================
LOOKUP RESULT
--------------------------------
Status: found
Team: Platform Team
================================
```
