# 프로젝트: 개인 웹/모바일 애플리케이션 (Go + Flutter)

## 시스템 아키텍처 및 디렉토리 구조
- Go 백엔드는 **도메인 중심(기능별) 디렉토리 구조**를 사용합니다.
- 관련된 모든 레이어(Layer)는 하나의 도메인/기능 폴더 내에 모아둡니다.
- 구조 예시:
  - `go-gate/internal/user/user_handler.go` (라우팅, gin.Context 처리, HTTP 요청/응답 담당)
  - `internal/user/user_service.go` (순수 비즈니스 로직 담당, 절대 gin.Context 사용 금지)
  - `internal/user/user_repository.go` (데이터베이스 쿼리 및 데이터 접근 담당)
  - `internal/user/user_model.go` (도메인 전용 구조체 정의)

## 빌드, 테스트 및 개발 명령어
- **백엔드 (Go)**:
  - 서버 실행: `cd go-gate && go run main.go`
  - 모듈 정리: `cd go-gate && go mod tidy`
- **프론트엔드 (Flutter)**:
  - 앱 실행 (디버그): `cd qr_app && flutter run -d chrome`
  - 패키지 다운로드: `cd qr_app && flutter pub get`

## 코드 스타일 및 지침

### 백엔드 (Go)
- **레이어별 책임 분리**:
  - `gin.Context`를 Service나 Repository 레이어로 넘기지 마세요. 프레임워크에 종속되지 않게 순수 Go 코드로 유지해야 합니다.
  - Handler는 Service에서 리턴된 에러를 받아서 적절한 HTTP 상태 코드(400, 404, 500 등)로 매핑해야 합니다.
- **에러 처리**: 에러는 항상 명시적으로 처리하세요. `panic`을 터트리는 것보다 에러를 감싸서(wrap) 리턴하는 방식을 선호합니다.

### 프론트엔드 (Flutter)
- UI 컴포넌트는 선언형 위젯(Stateless/StatefulWidget) 또는 사용하는 상태관리 라이브러리 규칙에 맞춰 작성하세요.
- 파일 및 폴더 이름은 `snake_case`, 클래스 이름은 `PascalCase` 등 다트(Dart) 스타일 가이드를 준수하세요.

### 공통 보안 사항
- 민감한 데이터나 API 키를 코드에 절대 하드코딩하지 마세요. 항상 환경 변수나 보안 설정 파일을 사용해야 합니다.