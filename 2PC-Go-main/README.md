# 2PC-Go 

This simple POC show how to implement a 2PC protocol in Go.

It is used to determine if a transaction can be committed or not.

## How to use it

### Start the coordinator

```bash
go run coordinator.go common.go model.go
```

### Start the participant

```bash
go run participant.go common.go model.go 0
```

```bash
go run participant.go common.go model.go 1
```
