# SIS Project: CEC Complaint Processor (Checkpoint 1)

A small synchronous, in-memory Go application. It is the core of a complaint
bot for the Central Election Committee (CEC) of KBTU: it takes a student's
complaint, validates it, decides whether it belongs to the CEC's jurisdiction
and how urgent it is, and tracks the complaint through a lifecycle.

In the SIS sense, a **task is a complaint**:

| Part | In this project |
|---|---|
| Input | `Input`: category, optional organization / school / teacher / subject, and the message |
| Processing | keyword classification of the message (`KeywordProcessor`) |
| Result | `Result`: validity (`VALID` / `REVIEW` / `INVALID`) and priority (`HIGH` / `MEDIUM` / `LOW`) |
| State | `Status`: `pending`, `processing`, `completed`, `failed` |

There is no Telegram, database or network access in this checkpoint. Everything
is stored in memory.

## Requirements

Go (see the `go` line in `go.mod`). No third-party dependencies.

## How to run

```bash
go run ./cmd/sis      # demo: creates and processes several complaints
go test ./...         # unit tests
go vet ./...
gofmt -l .            # prints nothing when the code is formatted
```

## Project structure

```
cmd/sis/main.go               demo program
internal/complaint/
  complaint.go                Complaint, Input, Result, Category, Validity, Priority
  status.go                   Status and the allowed transitions
  errors.go                   ErrInvalidInput, ErrNotFound, ErrInvalidTransition
  validate.go                 input validation
  processor.go                Processor interface and KeywordProcessor
  matcher.go                  text normalisation and keyword matching
  keywords.go                 keyword lists (Russian, English, Kazakh)
  store.go                    in-memory storage (map) and ID generation
  service.go                  Service: Create, Get, List, Process
  *_test.go                   unit tests
```

## Task lifecycle

```
pending ──> processing ──> completed
                 │
                 └───────> failed
```

| From | Allowed next status |
|---|---|
| `pending` | `processing` |
| `processing` | `completed`, `failed` |
| `completed` | none (terminal) |
| `failed` | none (terminal) |

- `Create` stores a new complaint as `pending`.
- `Process` moves it to `processing`, runs the processor, stores the result and
  moves it to `completed`. If the processor returns an error, the complaint becomes
  `failed` and has no result.
- Any other transition, for example processing a `completed` complaint again,
  is rejected with `ErrInvalidTransition` and changes nothing.
- A complaint classified as `INVALID` (outside the CEC's jurisdiction) is still
  `completed`: processing succeeded, and `INVALID` is its result.

## Validation rules

Invalid input is rejected with an error wrapping `ErrInvalidInput`, and nothing
is stored. Lengths are counted in characters (runes), not bytes, because
Cyrillic and Kazakh letters take two bytes.

| Rule | Applies to |
|---|---|
| category must be `org`, `staff`, `cafeteria`, `dormitory` or `other` | all |
| message must have at least 10 characters (after trimming spaces) | all |
| organization name must have at least 2 characters | `org` |
| school must be one of the KBTU schools in the list | `staff` |
| teacher name must have at least 3 characters | `staff` |

## How processing works

`KeywordProcessor` lower-cases the message, splits it into words and looks for
keywords from four lists (`keywords.go`).

- A keyword is a word **stem**: `нарушен` matches `нарушение`, `нарушения`. It matches only at
  the start of a word, so `кран` does not match inside `экран`.
- A keyword ending with `$` must match a **whole word**: `суп$` matches `суп`, not `супер`.
- Distinct keywords are counted once each.

**Validity**

1. `INVALID` if there are more "invalid" keywords than "valid" ones **and** at least 2 "invalid" ones.
2. Otherwise `VALID` if there is at least one "valid" keyword.
3. Otherwise `REVIEW`: nothing is conclusive, a human should look.

**Priority**: `HIGH` if any urgent keyword is present (threats, corruption,
discrimination, ...), otherwise `MEDIUM` if any medium keyword is present,
otherwise `LOW`.

## Error handling

Expected errors are returned, never `panic`. Check them with `errors.Is`:

- `ErrInvalidInput`: the input broke a validation rule;
- `ErrNotFound`: there is no complaint with that ID;
- `ErrInvalidTransition`: the operation is not allowed in the current status.

## Design decisions

- **`Processor` interface.** `Service` depends on `Processor`, not on the keyword
  classifier. Tests use a stub and a failing processor, and a later checkpoint can
  run processors in workers without changing `Service`.
- **Pointers in the store, copies outside.** The store holds `map[ID]*Complaint`
  so `Service` can update a complaint in place. `Create`, `Get`, `Process` and
  `List` return copies (including a copy of `Result`), so callers cannot change
  stored data by accident.
- **Pointer vs value receivers.** `Complaint.transition`, `KeywordProcessor` and
  the store use pointer receivers (they change state or own precompiled data).
  `Status` and `Input` methods use value receivers (small immutable values).
- **Sorted `List`.** Map iteration order is random in Go, so `List` sorts by ID.
- **IDs** are integers from 1, assigned by the store, and a rejected input does not use one.
- **Not thread-safe.** Checkpoint 1 is synchronous. Synchronisation comes with concurrency later.

## Known limitations

Keyword matching is a simple heuristic. It does not understand context or
negation ("no threats here" still matches `threat`), and word stems can match
unrelated longer words. The rules are deterministic and covered by tests.
