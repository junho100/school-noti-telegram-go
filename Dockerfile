# Build stage
FROM golang:1.22-alpine AS builder

# 필요한 시스템 패키지 설치
RUN apk add --no-cache git

# 작업 디렉터리 설정
WORKDIR /app

# Go 모듈 파일 복사 및 의존성 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 애플리케이션 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3.19

# 타임존 설정을 위한 패키지 설치
RUN apk add --no-cache tzdata

# 한국 시간대 설정
ENV TZ=Asia/Seoul

# 작업 디렉터리 설정
WORKDIR /app

# 빌드된 바이너리 복사
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# 실행 권한 설정
RUN chmod +x /app/main

# 애플리케이션 실행
CMD ["/app/main"] 