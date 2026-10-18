# CEC Complaint Processor

A small Go program for the SIS project (Checkpoint 1).
It works with student complaints. It checks a complaint, decides if it is
important for the CEC (Central Election Committee of KBTU), and saves the result.
Everything is stored in memory. There is no database.

## What is a task here?

One task is one complaint.

| Part       | What it is                                                         |
| ---------- | ------------------------------------------------------------------ |
| Input      | category, message, and some extra fields                           |
| Processing | look for keywords in the message                                   |
| Result     | validity (VALID, REVIEW, INVALID) and priority (HIGH, MEDIUM, LOW) |

## How to run

```bash
go run ./cmd/sis    # run the demo
go test ./...       # run the tests
```

## Project files

```
cmd/sis/main.go       demo program
internal/complaint/   main code and tests
```

## Task status

Every complaint has a status.

```
pending -> processing -> completed
                |
                +------> failed
```

- `pending`: the complaint is created, but not processed yet.
- `processing`: the complaint is being processed now.
- `completed`: processing is done and the result is saved.
- `failed`: an error happened, and there is no result.

You cannot go back or skip a step. For example, you cannot process
a `completed` complaint again. The program returns an error.

## Input rules

The program rejects a complaint if:

- the category is not known;
- the message is shorter than 10 characters;
- an `org` complaint has no organization name;
- a `staff` complaint has an unknown school or no teacher name.

## How processing works

The program makes the message lowercase and looks for keywords.

- Words about threats, bribes or discrimination give **HIGH** priority.
- Words about problems or unfair grades give **MEDIUM** priority.
- Other messages get **LOW** priority.
- Words about food, cleaning or internet show that the complaint is
  **not** for the CEC. The result is **INVALID**.
- If nothing is clear, the result is **REVIEW**. A person must check it.

## Errors

The program returns errors. It does not use `panic`.

- `ErrInvalidInput`: the input is wrong.
- `ErrNotFound`: there is no complaint with this ID.
- `ErrInvalidTransition`: this action is not allowed for the current status.

## Limits

The keyword search is simple. It does not understand the meaning of a sentence.
For example, "there are no threats" still finds the word "threat".
