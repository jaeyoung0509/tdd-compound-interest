## 복리 BNPL 도메인

Go로 작성된 연체/복리 BNPL 시스템을 위한 도메인 주도 설계 스케치입니다.

### 패키지 구조
- `domain/payment`: Payment 애그리게이트, 생명주기 상태 (`Status*`), 연체 정보, 도메인 에러.
- `domain/money`: 통화 스케일 처리 및 BPS 헬퍼가 포함된 Decimal 기반 Money 값 객체.
- `domain/user`: 스코프된 ID와 유효성 검사가 포함된 User 애그리게이트 스텁.
- `domain/shared`: 공통 ID 헬퍼 (ULID).

### 주요 설계 포인트
- Payment는 불변식(중복 결제 방지, 결제 후 연체 처리 불가)을 보호하기 위해 상태 전이(`Pay`, `MarkOverdue`)를 캡슐화합니다.
- Money는 정밀도 보존을 위해 `shopspring/decimal`과 통화별 스케일(KRW:0, USD:2)을 사용하며, BPS 헬퍼는 이자 계산을 지원합니다.
- ID는 ULID를 통해 생성되며, 제로 값 누수를 방지하기 위해 애그리게이트별로 래핑됩니다.

### 테스트 실행
```bash
go test ./...
```

### 다음 단계
- 시계/이자율 인터페이스를 도입하고 Payment에 복리 이자 발생을 구현합니다.
- 도메인을 전송/타입으로부터 독립적으로 유지하면서 영속성/어댑터를 추가합니다.
- 필요에 따라 더 많은 통화와 반올림 정책으로 Money를 확장합니다.
