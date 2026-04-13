# OddsHub

OddsHub is a service for collecting, storing, and serving sports betting odds from multiple providers.

## Features
- Odds aggregation
- Historical storage
- REST API

## Tech stack
- Go
- PostgreSQL
- Redis

## Folders structure
oddshub/
├── cmd/
│   └── api/            # входная точка (main.go)
├── internal/
│   ├── collector/      # тянет данные из API
│   ├── storage/        # работа с БД
│   ├── service/        # бизнес-логика
│   └── transport/      # HTTP (handlers)
├── pkg/                # общие утилиты
├── configs/
├── migrations/
├── go.mod
└── README.md
