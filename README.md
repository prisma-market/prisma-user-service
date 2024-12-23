# prisma-user-service
This is prisma market's user service


## 실행방법

```bash
# 서비스 빌드 및 시작
docker-compose up -d --build

# 로그 확인
docker-compose logs -f
```



시작:
```bash
docker-compose up -d      # 백그라운드로 실행
docker-compose up -d --build  # 이미지 새로 빌드하고 실행
```

종료:
```bash
docker-compose down      # 컨테이너, 네트워크 모두 종료/제거
docker-compose down -v   # 볼륨까지 모두 제거
```

그 외 유용한 명령어들:
```bash
docker-compose ps        # 실행 중인 컨테이너 상태 확인
docker-compose logs -f   # 로그 실시간 확인
docker-compose restart   # 모든 서비스 재시작
docker-compose stop     # 컨테이너만 중지 (삭제는 안 함)
docker-compose start    # 중지된 컨테이너 시작
```

특정 서비스만 조작하고 싶을 때:
```bash
docker-compose restart user-service  # user-service만 재시작
docker-compose logs mongodb         # mongodb 로그만 보기
```