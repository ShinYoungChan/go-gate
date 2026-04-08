# go-gate

## 🛠 Backend (Go Gin Framework) 설정

### 📦 라이브러리 설치
터미널에서 아래 명령어를 실행하여 필요한 패키지를 설치합니다.

```bash
# Gin Web Framework 설치
go get -u github.com/gin-gonic/gin (https://github.com/gin-gonic/gin)

# JWT 인증 라이브러리 설치
go get github.com/golang-jwt/jwt/v5 (https://github.com/golang-jwt/jwt/v5)

# GORM 및 SQLite 드라이버 설치
go get -u gorm.io/gorm
go get -u github.com/glebarez/sqlite (https://github.com/glebarez/sqlite)

# GORM 라이브러리 드라이버 교체
go get -u github.com/glebarez/sqlite

#flutter dio 패키지 설치
flutter pub add dio

#flutter chrome 실행
flutter run -d chrome

```


### 테스트 데이터 삽입
users 데이터는 회원가입 데이터 전송으로 생성

```sql
-- 1. 입장 장소 (위치 설정)
INSERT INTO locations (place_name, category, lat, lon, address, image_url) VALUES
('강남 피트니스 센터', '헬스장', 37.4979, 127.0276, '서울시 강남구 테헤란로 123', 'https://picsum.photos/id/10/200/200'),
('서초 바디빌드', '헬스장', 37.4849, 127.0348, '서울시 서초구 강남대로 456', 'https://picsum.photos/id/11/200/200'),
('홍대 머슬팩토리', '헬스장', 37.5565, 126.9239, '서울시 마포구 양화로 78', 'https://picsum.photos/id/12/200/200'),
('잠실 요가 앤 필라테스', '요가', 37.5133, 127.1001, '서울시 송파구 올림픽로 567', 'https://picsum.photos/id/13/200/200');

-- 2. 멤버십 상품 종류 (아이템 목록)
INSERT INTO membership_items (title, type, duration_days, amount)
VALUES 
('10회 이용권', 'count', 90, 150000),
('1개월 자유이용권', 'period', 30, 200000);

-- 3. 사용자의 멤버십 보유 현황 (방금 가입한 유저 ID가 1번이라고 가정)
-- 10회권 중 10회가 남았고, 오늘부터 90일간 유효한 상태
INSERT INTO user_memberships (user_id, item_id, stt_dt, end_dt, is_count_type, count, is_valid, amount)
VALUES (1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '90 days', true, 10, true, 150000);
```