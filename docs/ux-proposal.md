# mcrpc 공개 API UX 개선 방향성 제안

대상: `github.com/fi-xz/mcrpc` v0 공개 API
작성 시점 기준 커밋: `7112592`

이 문서는 라이브러리를 처음 쓰는 외부 사용자 관점에서 발견되는 문제를 정리하고,
각각에 대해 방향성과 대략적인 호환성 비용을 기록한다. 우선순위는
**정확성 문제 > 사용 불가 문제 > 관용성(idiomatic) 문제** 순이다.

각 항목의 상태는 다음과 같이 표기한다.

- **[완료]** — 반영됨. 아래 본문은 반영 전 상태에 대한 기록이다.

17개 항목 전부 반영 완료. 연결 수립 방식이 바뀌었으므로 breaking change이며,
마이그레이션 표는 [README](../README.md#compatibility-notes)에 있다.

실서버 검증까지 마쳤고, 그 과정에서 제안서에 없던 결함 셋을 더 찾았다.
아래 "실서버 검증" 절에 정리했다.

---

## P0 — 정확성 문제 (버그에 가까움)

### 1. 핸드셰이크 실패 시 `nil` 에러 반환 — **[완료]**

```go
// mcrpc.go
if response.StatusCode != http.StatusSwitchingProtocols {
    return nil, err   // err은 이 시점에 반드시 nil
}
```

이 분기는 `err == nil`이 보장된 이후이므로 `(nil, nil)`을 반환한다. 호출자는
`if err != nil` 검사를 통과한 뒤 nil 클라이언트를 역참조하고 패닉한다.

**방향:** 상태 코드와 서버 응답 본문을 담은 전용 에러를 반환한다.

```go
return nil, fmt.Errorf("mcrpc: handshake failed: unexpected status %s", response.Status)
```

인증 실패(401/403)는 별도 센티넬(`ErrUnauthorized`)로 노출하면 호출자가
"시크릿이 틀렸다"와 "서버가 죽었다"를 구분할 수 있다.

> **반영:** 상태 코드를 담은 에러를 반환하도록 수정. 다이얼 실패도
> `fmt.Errorf("mcrpc: dial %s: %w", ...)`로 래핑해 어느 주소에서 실패했는지 남긴다.
> `ErrUnauthorized` 센티넬은 8번의 에러 타입 정리와 함께 다루는 편이 낫다고 보아
> 이번에는 넣지 않았다.

### 2. TLS 사용 여부가 클라이언트 인증서 유무에 묶여 있음 — **[완료]**

```go
protocol := "ws"
if cert != nil {
    protocol = "wss"
}
```

- 서버가 TLS를 켰지만 **클라이언트 인증서를 요구하지 않는** 구성(가장 흔한 구성)에서
  `CreateWithTLS(..., nil, true)`를 호출하면 조용히 평문 `ws://`로 연결을 시도한다.
- 사용자는 TLS를 요청했는데 실패도 경고도 없이 다운그레이드된다.

**방향:** 스킴 결정을 인증서가 아니라 명시적 의도에 연결한다. 아래 5번의 옵션 패턴과
합쳐 `mcrpc.WithTLS(cfg *tls.Config)`가 있으면 `wss`, 없으면 `ws`로 결정한다.
인증서는 그 `tls.Config` 안의 선택 항목이 된다.

> **반영:** 5번의 옵션 패턴을 기다리지 않고, 내부 `createClient`에 `useTLS` 인자를
> 추가해 스킴을 호출자의 의도에 직접 연결했다. `CreateWithTLS`는 인증서가 nil이어도
> 항상 `wss`로 붙고 `InsecureSkipVerify`도 반영된다. 공개 시그니처는 그대로다.
> 5번을 진행하면 `useTLS`는 `WithTLS` 옵션의 존재 여부로 대체된다.
>
> 이 변경은 동작 변경이다. `CreateWithTLS(..., nil, ...)`로 호출하면서 실제로는
> 평문 연결에 의존하고 있던 코드가 있다면 `Create`로 바꿔야 한다.

### 3. 핸들러 필드 대입은 데이터 레이스 — **[완료]**

`Create`가 반환되는 시점에 `jsonrpc2.NewConn`은 이미 수신 고루틴을 돌리고 있다.
사용자가 그 뒤에 `client.OnPlayerJoined = ...`를 대입하는 README 권장 패턴은
수신 고루틴의 필드 읽기와 경합한다. 서버가 연결 직후 `server/status`를 보내는
구성에서는 실제로 알림 유실도 발생한다.

**방향(택1, 아래로 갈수록 침습적):**

- **(a) 연결 시작을 분리한다.** `Create`는 다이얼만 하고, `client.Start(ctx)` 또는
  `client.Listen(ctx)` 호출 시점에 `jsonrpc2.NewConn`을 만든다. 핸들러 등록 구간이
  레이스 없이 확보된다. 기존 코드와의 차이는 한 줄 추가라 마이그레이션이 가장 싸다.
- **(b) 생성 시 주입한다.** `mcrpc.WithHandler(h Handler)` 옵션으로 넘긴다.
- **(c) setter로 감싼다.** `SetOnPlayerJoined(fn)` + 내부 `sync.RWMutex`.
  API 표면이 25개 이상 늘어나므로 권장하지 않는다.

권장: **(a) + (b) 조합**. 필드 대입은 당분간 유지하되 문서에서 "Start 이전에만
대입하라"고 못 박는다.

> **반영:** (a)와 (b)를 합치되, 검토 과정에서 (a)를 더 밀어붙였다. `Create`가
> 다이얼까지 하고 `Start`가 jsonrpc2만 붙이는 형태는 **리더 없는 소켓을 쥔
> 반쯤 살아있는 객체**를 만든다. `Start`를 안 부르면 소켓이 새고, 생성자에는
> 여전히 ctx와 error가 필요하다.
>
> 그래서 다이얼까지 `Start`로 옮겼다. `New`는 I/O도 고루틴도 없이 값만 만든다:
>
> ```go
> client := mcrpc.New(addr, secret, mcrpc.WithHandler(h))
> err := client.Start(ctx)
> ```
>
> `Start` 이전에는 수신 고루틴이 **존재하지 않으므로** 경합이 성립 불가다. 문서
> 규칙이 아니라 구조로 막힌다. (b)의 `WithHandler`가 유일한 등록 경로이고,
> `Handler`는 값으로 복사되므로 호출자가 자기 사본을 나중에 건드려도 영향 없다.
>
> 부수 효과로 서버 없이 도는 테스트가 가능해졌다 — `New`가 연결을 요구하지
> 않으므로. `client_test.go`의 `newFakeServer`(httptest WebSocket)로 연결·종료·
> 재시작·알림 전달·에러 매핑을 `-race`로 검증한다.

### 4. `Message` / 밴 구조체의 빈 필드가 그대로 직렬화됨 — **[완료]**

```go
type Message struct {
    Translatable       string   `json:"translatable"`
    TranslatableParams []string `json:"translatableParams"`
    Literal            string   `json:"literal"`
}
```

`omitempty`가 없어 리터럴 메시지 하나를 보내도
`{"translatable":"","translatableParams":null,"literal":"hi"}`가 전송된다.
MSMP는 `translatable`과 `literal` 중 하나만 오는 것을 기대하므로, 서버 구현에 따라
거절되거나 빈 번역 키로 해석될 수 있다. `UserBan.Expires`(영구 밴이면 빈 문자열),
`IPBan.Expires`도 같은 문제를 가진다.

**방향:** 해당 필드에 `omitempty`를 붙이고, 상호 배타적인 조합은 생성자로 강제한다.

```go
func LiteralMessage(text string) Message
func TranslatableMessage(key string, params ...string) Message
```

> **반영:** `Message`의 세 필드와 `UserBan.Expires`, `IPBan.Expires`에 `omitempty`를
> 붙였다. `SystemMessage.Overlay` 같은 bool은 `false`가 유의미한 값이라 제외했다.
> 생성자 `LiteralMessage` / `TranslatableMessage`는 [helpers.go](../helpers.go)에 추가.
> 직렬화 결과는 `TestMessageConstructorsOmitUnusedFields`,
> `TestPermanentBanOmitsExpires`가 검증한다.
>
> `Expires`에 대한 가정: MSMP가 영구 밴을 "필드 부재"로 받는다고 전제했다.
> 실서버로 확인됨 — 아래 "실서버 검증" 절 참조.

---

## P1 — 사용성: "쓸 수는 있으나 매번 불편함"

### 5. 생성자 시그니처가 위치 인자 위주라 호출부가 읽히지 않음 — **[완료]**

```go
client, err := mcrpc.CreateWithTLS(ctx, "localhost", 8080, "your-secret", &cert, true)
```

마지막 `true`가 `insecure`(인증서 검증 생략)라는 사실이 호출부에 드러나지 않는다.
보안에 영향을 주는 플래그가 가장 눈에 안 띄는 자리에 있다.

**방향:** 함수형 옵션 + 관용적 이름.

```go
client, err := mcrpc.Dial(ctx, "localhost:8080",
    mcrpc.WithSecret(secret),
    mcrpc.WithTLS(&tls.Config{Certificates: []tls.Certificate{cert}}),
    mcrpc.WithInsecureSkipVerify(), // 이름 자체가 경고로 동작
)
```

- `host string, port int` → `addr string`으로 합치면 `net.JoinHostPort` 관례와 맞고
  IPv6 리터럴도 자연스럽게 처리된다.
- `Create`/`CreateWithTLS`는 `Dial`을 호출하는 얇은 래퍼로 남겨 deprecated 표기한다.

> **반영:** 이름은 `Dial`이 아니라 `New`로 했다 — 더 이상 다이얼하지 않으므로
> `Dial`은 거짓말이 된다. 옵션은 `WithTLS`, `WithClientCertificate`,
> `WithInsecureSkipVerify`, `WithHandshakeTimeout`, `WithHandler`.
>
> TLS 관련 옵션은 서로 순서 무관하다. 각자 플래그만 세우고, 실제 `tls.Config`는
> `Start` 시점에 `tlsSettings()`가 조립한다. `WithInsecureSkipVerify()` 뒤에
> `WithTLS(cfg)`를 써도 결과가 같다.
>
> `addr`은 `"host:port"` 문자열 하나. 분리해서 들고 있는 호출자를 위해
> `NewHostPort(host, port, secret, opts...)`를 함께 둔다.
>
> `Create` / `CreateWithTLS`는 deprecated 래퍼로 남기지 않고 **삭제**했다. 래퍼는
> 핸들러를 등록할 방법이 없으므로, 남겨두면 "핸들러 없이 이미 시작된 클라이언트를
> 받아서 나중에 대입한다"는 경합 경로만 계속 살아있게 된다. 올바른 경로가 생긴
> 이상 그건 편의가 아니라 함정이다.

### 6. `MCRPCClient` 이름이 중복된다 — **[완료]**

`mcrpc.MCRPCClient`는 패키지명을 두 번 읽게 만든다. Go 관례는 `mcrpc.Client`다.

**방향:** `type Client struct{...}` 로 개명하고 `type MCRPCClient = Client` 별칭을
남긴다. 별칭이므로 기존 코드는 그대로 컴파일된다.

> **반영:** `Client`로 개명. 처음에는 `MCRPCClient = Client` 별칭을 deprecated로
> 남겼으나, 4·5단계에서 breaking change가 승인되면서 별칭도 제거했다.

### 7. 내부 커넥션이 공개 필드로 노출됨 — **[완료]**

```go
JSONRPCConn   *jsonrpc2.Conn
WebsocketConn *websocket.Conn
```

- `jsonrpc2`, `gorilla/websocket`이 공개 API의 일부가 되어, 의존성 교체나 major
  업그레이드가 곧 mcrpc의 breaking change가 된다.
- 사용자가 `client.JSONRPCConn.Close()`를 직접 부르면 `closed` 플래그가 어긋난다.

**방향:** 비공개 필드로 내리고 필요한 것만 메서드로 노출한다
(`Ping`, `RemoteAddr`, 그리고 이미 있는 `Close`/`IsClosed`/`DisconnectNotify`).
탈출구가 필요하면 `Underlying() *jsonrpc2.Conn` 한 개만 남기고 "지원 범위 밖" 명시.

> **반영:** 두 필드 모두 비공개. `jsonrpc2`와 `gorilla/websocket`이 공개 API에서
> 사라졌다. 재시작을 지원하려면 세션 상태가 mutable해야 하므로 `mu sync.RWMutex`로
> 보호하고 `conn()` 접근자를 통해서만 읽는다.
>
> `IsClosed()`는 `IsRunning()`으로 대체했다. 재시작 가능한 객체에서 "closed"는
> 종결을 암시해 오해를 부른다.
>
> `Underlying()` 탈출구는 넣지 않았다. 필요해지면 그때 추가하는 편이,
> 넣어두고 못 빼는 것보다 낫다.

### 8. 에러가 원본 `jsonrpc2.Error` 그대로 올라온다 — **[완료]**

호출자가 "플레이어가 존재하지 않음"과 "권한 없음"과 "연결 끊김"을 구분하려면
`jsonrpc2`를 직접 import해서 코드 번호를 봐야 한다.

**방향:** 패키지 에러 타입 하나를 정의하고 모든 `Call`을 감싼다.

```go
type Error struct {
    Method  string
    Code    int64
    Message string
    wrapped error
}
func (e *Error) Error() string { ... }
func (e *Error) Unwrap() error { return e.wrapped }

var ErrNotConnected = errors.New("mcrpc: not connected")
```

`errors.As`로 코드 분기가 가능해지고, `%w` 래핑으로 AGENTS.md의 에러 규칙과도 맞는다.
현재 모든 메서드가 `err`을 맨몸으로 반환하고 있어 이 변경은 한 곳(`call` 헬퍼 도입)
에서 처리 가능하다.

> **반영:** [errors.go](../errors.go)에 `*Error`, `ErrNotConnected`,
> `(*Client).call` 추가. 모든 메서드(68개 호출부)가 `c.JSONRPCConn.Call` 대신
> `c.call`을 거친다. `Unwrap`이 원본을 노출하므로 `errors.As(err, &rpcErr)`로
> `*jsonrpc2.Error`에도 여전히 도달 가능하다.
>
> `call`은 `error` 인터페이스를 반환하고 성공 시 명시적으로 `nil`을 돌려준다
> (typed-nil 함정 회피). 반환 타입이 `error`로 동일해 호출부는 breaking 없음.

### 9. 반복 파라미터가 슬라이스 고정 — **[완료]**

```go
client.AddAllowlist(ctx, []mcrpc.Player{p})
client.RemoveIPBanlist(ctx, []string{"1.2.3.4"})
```

한 명만 다루는 경우가 압도적으로 많은데 매번 슬라이스 리터럴을 쓴다.

**방향:** 가변 인자로 바꾼다. `AddAllowlist(ctx context.Context, players ...Player)`.
`[]Player`를 넘기던 코드는 `add...`만 붙이면 되고, 슬라이스 하나를 넘기던 호출은
`AddAllowlist(ctx, list...)`로 컴파일된다.
플레이어 이름만으로 지정할 수 있는 헬퍼(`mcrpc.PlayerByName("fi_xz")`)도 같이 두면
UUID를 모르는 흔한 상황이 해결된다.

> **반영:** `Add*` / `Remove*` / `KickPlayers`만 가변 인자로 바꿨다. `Set*`은
> 슬라이스를 유지한다 — "목록 전체를 이걸로 교체"는 의미상 리스트 연산이고,
> `SetAllowlist(ctx)`가 "전부 지움"으로 읽히는 건 `ClearAllowlist`와 겹쳐 혼란스럽다.
>
> 가변 인자를 0개로 호출하면 nil 슬라이스가 되어 JSON `null`로 직렬화된다.
> 이걸 막으려고 모든 리스트 파라미터를 `nonNilSlice`로 감쌌다.
> `TestListParametersSerialiseAsArrays`가 검증한다.

### 10. 시간 값이 문자열 — **[완료]**

`UserBan.Expires`, `IPBan.Expires`가 ISO 8601 `string`이라 사용자가 포맷을 직접 안다.

**방향:** 필드는 프로토콜 호환을 위해 `string`으로 두되, 헬퍼를 제공한다.

```go
func (b UserBan) ExpiresAt() (time.Time, bool)   // false면 영구 밴
func BanUntil(t time.Time) string
```

> **반영:** `ExpiresAt`, `IsPermanent`를 `UserBan` / `IPBan` 양쪽에 추가하고,
> `BanUntil`과 `BanTimeFormat`(= RFC3339) 상수를 공개했다.
>
> 메서드는 [internal/types/accessors.go](../internal/types/accessors.go)에 있다.
> 별칭은 원본과 동일 타입이라 루트 패키지에서는 메서드를 붙일 수 없기 때문이다.
> 별칭을 통해 그대로 승격되므로 사용자에게는 `mcrpc.UserBan`의 메서드로 보인다.

### 11. 게임룰 값이 `any` — **[완료]**

`UntypedGameRule.Value any` + `TypedGameRule.Type string`이라 읽는 쪽에서 매번
타입 스위치를 쓴다.

**방향:** 접근자와 생성자를 붙인다.

```go
func (g TypedGameRule) Bool() (bool, bool)
func (g TypedGameRule) Int() (int, bool)
func BoolRule(key string, v bool) UntypedGameRule
func IntRule(key string, v int) UntypedGameRule
```

> **반영:** 접근자 `Bool` / `Int` / `StringValue`를 `UntypedGameRule`에 추가했다
> (`TypedGameRule`은 임베딩이라 자동 승격). 생성자는 `BoolRule` / `IntRule` /
> `StringRule`.
>
> `Int`는 JSON 디코딩 결과가 `float64`로 들어오는 점, 일부 룰이 문자열로 전송되는
> 점을 모두 처리한다. `1.5` 같은 비정수 float는 거부한다.
> 접근자 이름이 `String`이 아니라 `StringValue`인 이유는 `String() (string, bool)`이
> `fmt.Stringer` 규약과 어긋나 `go vet`의 stdmethods 검사에 걸리기 때문이다.

---

## P2 — 관용성 / 정리

### 12. `context context.Context` 파라미터 섀도잉 — **[완료]**

모든 메서드가 파라미터 이름으로 `context`를 쓴다. 함수 본문 안에서 `context` 패키지를
쓸 수 없고(`context.WithTimeout` 등), Go 컨벤션(`ctx`)과도 어긋난다.
AGENTS.md 예제조차 `ctx`를 쓰고 있어 문서와 구현이 불일치한다.

**방향:** 전부 `ctx`로 개명한다. 파라미터 이름은 API 호환성에 영향이 없다.

> **반영:** 전 메서드 개명 완료. AGENTS.md 예제와 구현이 일치한다.

### 13. 오타 및 섀도잉된 빌트인 — **[완료]**

- `createMCRCPClient` → `createClient` (MCRPC가 MCRCP로 뒤집힘)
- `SetMaxPlayers(context context.Context, max int)` → `max`는 Go 1.21+ 빌트인 섀도잉.
  `count` 또는 `limit`으로.

> **반영:** 둘 다 처리. `max` → `limit`. 둘 다 비공개 식별자/파라미터 이름이라
> API 영향 없음.

### 14. 알림 파싱 실패가 조용히 삼켜짐 — **[완료]**

```go
if err := json.Unmarshal(params, &player); err == nil {
    c.OnPlayerJoined(player)
}
```

서버 스키마가 바뀌면 알림이 아무 흔적 없이 사라진다. 디버깅 단서가 0이다.

**방향:** `OnError func(method string, err error)` 훅을 하나 추가하고, 미설정 시에는
`log.Printf` 대신 무시하되 최소한 훅 자체는 제공한다. `OnNotification`이 이미 raw
`json.RawMessage`를 주므로 그것과 짝을 이룬다.

> **반영:** `Handler.OnError func(method string, err error)` 추가. 디코딩은 제네릭
> 헬퍼 `decodeNotification`로 일원화해, 14군데에 흩어져 있던
> `if err := json.Unmarshal(...); err == nil` 패턴을 한 곳으로 모았다.
> 미설정 시 동작은 종전대로 무시.

### 15. 재연결 경로 부재 — **[완료]**

`DisconnectNotify`만 있고, 끊긴 뒤 상태를 유지한 채 다시 붙는 수단이 없다.
사용자는 `Dial`부터 핸들러 25개 재등록까지 전부 다시 짠다.

**방향:** 옵션 기반 자동 재연결(`WithReconnect(backoff)`)이 이상적이나 범위가 크다.
최소 단계로, 핸들러 묶음을 값으로 들고 다닐 수 있는 `Handler` 구조체를 도입하면
(3-(b)) 재연결 코드가 몇 줄로 줄어든다.

> **반영:** 자동 재연결은 넣지 않았지만, 수동 재연결이 `Start` 재호출 한 줄이 됐다.
> `Close`는 설정과 `Handler`를 지우지 않고 세션만 해제하므로 같은 값을 다시 쓸 수
> 있다. `TestRestartAfterClose`가 검증한다.
>
> `Start`가 띄우는 감시 고루틴이 서버 hangup과 ctx 취소를 모두 `Close`로 변환하므로,
> 재연결 루프는 `<-client.DisconnectNotify()`를 기다렸다가 다시 `Start`하면 된다.
> 자동화는 이 위에 사용자 코드로 얹을 수 있어, 라이브러리에 백오프 정책을 박아넣는
> 것보다 낫다고 판단했다.

---

### 16. 와이어에 실제로 무엇이 오가는지 볼 수 없다 — **[완료]**

응답이 성공했다고 해서 서버가 우리가 보낸 JSON을 의도대로 해석했다는 뜻은 아니다.
구조체 태그가 실제로 무엇을 만드는지, 생략된 필드가 와이어에서 어떻게 보이는지,
untyped 값이 어떤 JSON 타입으로 도착하는지 확인할 방법이 없었다.

**방향:** `jsonrpc2`의 `OnSend`/`OnRecv` ConnOpt를 옵션으로 노출한다.

> **반영:** [trace.go](../trace.go)에 `WithTrace(func(TraceMessage))` 추가.
> 더불어 서버가 `rpc.discover`로 자기 API 스키마(OpenRPC)를 돌려준다는 사실을
> 활용해, 추측 대신 스키마를 대조하는 테스트를 두었다:
> `TestWireNotificationsTakeOneParameter`(모든 알림이 파라미터 1개 이하 — 디코딩
> 모델의 전제), `TestWireRequestParameterNames`(setter 36개의 JSON 필드명 대조).
> 이후의 프로토콜 변경은 사용자 제보가 아니라 테스트 실패로 드러난다.
> `TraceMessage`는 방향·메서드·ID·params·result·에러 코드를 담고, `String()`이
> 한 줄로 렌더한다.
>
> `jsonrpc2.Error`를 그대로 노출하지 않고 `ErrorCode`/`ErrorMessage`로 평탄화했다.
> #7에서 숨긴 의존성을 진단 타입으로 되살리면 의미가 없기 때문이다.
>
> 프레임 전체 raw 바이트 대신 params/result를 쓴다. `ObjectStream`을 감싸면 쓰기는
> 이중 marshal이 되고, 읽기는 `ReadJSON`이 내부 타입으로 바로 디코딩해 raw 확보가
> 안 된다. 검증 대상이 어차피 params/result JSON이라 이쪽이 정확하고 단순하다.
>
> 통합 테스트에서는 `TEST_TRACE=1 go test -v ...`로 켜진다. 아래 결함 셋은 전부
> 이걸 켜고서야 보였다.

### 17. 서버 버전을 알 수 없다 — **[완료]**

동작이 버전에 따라 갈리는데(게임룰 표현, 알림 종류), 클라이언트가 기준으로 삼을
버전을 노출하지 않았다. `GetServerStatus`가 주는 것은 **마인크래프트** 버전이고,
분기 기준이 되어야 하는 것은 **관리 API 버전**이다.

`rpc.discover`가 돌려주는 OpenRPC 문서의 `info.version`이 그 값이다. 실측:

| API | 마인크래프트 | 이 라이브러리에 영향 |
|---|---|---|
| 1.0.0 | 1.21.9 | 게임룰 값이 문자열, 키가 camelCase |
| 3.0.0 | 26.2 | 게임룰 값이 네이티브 JSON, 키가 namespaced. `notification/server/activity` 추가 |
| 3.1.0 | 26.3 | `notification/world/upgrade_*` 추가 |

> **반영:** `(*Client).APIVersion(ctx)` 추가. `rpc.discover`를 호출해 `info.version`만
> 꺼낸다. 상수는 `protocol.MethodDiscover` (JSON-RPC 예약 `rpc.` 네임스페이스라
> `minecraft:` 접두사가 없다).
>
> 이 발견으로 `world/upgrade_*` 4개 상수의 정체도 해명됐다. 스키마에 없었던 것은
> 오류가 아니라 **26.3(API 3.1.0)에서 추가**되기 때문이다(실서버로 확인, 아래 참조).
> 파라미터도
> 문서와 일치한다 — `upgrade_progress`는 0~1 숫자(초당 1회로 rate limit),
> `upgrade_failed`는 reason 문자열, 나머지 둘은 파라미터 없음.

## 권장 실행 순서

| 단계 | 내용 | 호환성 | 상태 |
|---|---|---|---|
| 1 | P0 #1, #2, #4 (핸드셰이크 에러, TLS 스킴, `omitempty`) | 동작 수정, 시그니처 불변 | 완료 |
| 2 | #12, #13 (`ctx` 개명, 오타) | 완전 호환 | 완료 |
| 3 | #6 `Client` 개명, #8 에러 타입, #10/#11 헬퍼 | 추가만 함, 완전 호환 | 완료 |
| 4 | #3 `New`/`Start` 분리, #5 옵션, #14 `OnError`, #15 재연결 | breaking | 완료 |
| 5 | #7 커넥션 필드 비공개, #9 가변 인자 | breaking | 완료 |
| 6 | #16 `WithTrace`, #17 `APIVersion`, 실서버 검증 | 추가만 함 | 완료 |

전 단계 완료. 4·5단계는 deprecated 유예 없이 한 번에 처리했다 — 라이브러리 사용자가
작성자 본인뿐이라 유예 기간이 비용만 되고, `Create`를 남겨두면 "핸들러 없이 이미
시작된 클라이언트"라는 위험한 경로만 살아남기 때문이다.

남은 것은 v1.0.0 태깅이며, 그 전에 아래 미검증 항목을 실서버로 확인하는 것을 권한다.

## 실트래픽 확인

`players/joined` / `players/left`는 실제 클라이언트 접속 없이는 확인할 수 없어
마지막까지 남아 있었다. 양 버전 모두 실접속으로 확인 완료:

| 서버 | joined | left | 비고 |
|---|---|---|---|
| 1.21.9 (API 1.0.0) | 전달됨, UUID 포함 | 전달됨 | 디코딩 실패 0건 |
| 26.2 (API 3.0.0) | 전달됨, UUID 포함 | 전달됨 | `server/activity`도 함께 발화 |

두 알림은 스키마로 shape(`[player]` positional)를 확정한 뒤 코드를 고쳤고, 실트래픽이
그 수정을 확인해 준 형태다. `server/activity`는 26.2 전용이며 스키마상 `params: []`,
30초당 1회로 rate limit — `OnServerActivity func()`와 일치한다.

프로브(`TestProbePlayerNotifications`)는 사람이 접속해야 하므로 `PROBE_WINDOW`가
설정된 경우에만 동작하고, 그렇지 않으면 skip한다.

## world/upgrade 확인

26.3 스냅샷(API 3.1.0)에 변환이 필요한 월드를 물려 부팅시키고, 프로브를 먼저 띄워
포트가 열리는 즉시 붙는 방식으로 확인했다. 4개 중 3개 실증:

| 알림 | 실제 params | 결과 |
|---|---|---|
| `upgrade_started` | 없음 | 확인 |
| `upgrade_progress` | `[0.0]` | 확인 |
| `upgrade_finished` | 없음 | 확인 |
| `upgrade_failed` | `["reason"]` (문서 기준) | **미확인** — 변환 실패를 유도하지 않음 |

`[0.0]`이 중요하다. positional 인자 배열 가정이 **객체가 아닌 스칼라 payload에서도**
성립함을 보여준다 — `decodeParam`이 기대는 바로 그 지점이고, 이제 숫자 타입까지
실측으로 뒷받침된다.

첫 시도는 실패했다. 작은 월드라 변환이 관리 서버 바인딩보다 먼저 끝나, 붙자마자
`upgrade_finished`만 받았다. 큰 월드로 바꾸고 재시도 간격을 250ms에서 20ms로 줄이니
전 구간이 잡혔다. 변환은 약 1.6초 걸렸고, 초당 1회 rate limit 때문에 `progress`는
한 번만 왔다.

이 알림들은 **서버가 완전히 뜨기 전에** 발화하므로, 핸들러를 생성 시점에 등록하고
일찍 접속해야 받을 수 있다. `New`/`Start` 분리(#3)가 없었다면 이 확인 자체가
불가능했을 것이다.

## 함께 반영된 항목 (별도 번호 없음)

- **내부 타입이 공개 시그니처에 노출되던 문제**: `types.SystemMessage`,
  `types.UntypedGameRule`이 `SendSystemMessage`, `UpdateGameRule`의 파라미터
  타입이라 외부 사용자가 값을 만들 수 없었다. 모든 별칭을 `types.go` 한 곳으로
  모으고 `Message`, `Version`, `SystemMessage`, `UntypedGameRule` 별칭을 추가했다.
  재발 방지는 외부 테스트 패키지인 `api_public_test.go`가 컴파일 타임에 막는다.
- **IPv6 주소 처리**: 다이얼 문자열 조립을 `fmt.Sprintf("%s://%s:%d")`에서
  `net.JoinHostPort`로 교체했다. 기존 방식은 `::1` 같은 리터럴에서 대괄호를 빠뜨려
  잘못된 URL을 만들었다.
- **단위 테스트 추가**: 기존 테스트는 전부 실서버가 필요해 CI에서 `t.Skip`된다.
  새 헬퍼·접근자·에러 타입은 서버 없이 도는 테스트로 덮었다
  ([helpers_test.go](../helpers_test.go), [errors_test.go](../errors_test.go),
  [api_public_test.go](../api_public_test.go)).

## 실서버 검증

실서버 **두 대**로 확인 완료 — Minecraft **1.21.9** (protocol 773)과 **26.2**
(protocol 776), 둘 다 TLS. 양쪽 모두 전체 스위트 통과, 커버리지 95.2%.

두 버전을 함께 돌린 것이 결정적이었다. 26.2만 봤다면 게임룰 값 표현 차이와
하트비트 회귀를 놓쳤을 것이고, `allow_flight`를 버전 차이로 오진했을 것이다.

| 항목 | 결과 |
|---|---|
| 1. `expires` 생략 = 영구 밴 (#4) | **확인.** 서버가 수용하고 응답에도 `expires` 없음. `BanUntil`의 RFC3339도 그대로 왕복 |
| 2. 클라이언트 인증서 없는 `wss` (#2) | **확인.** `InsecureSkipVerify` 없이 정상 핸드셰이크 |
| 3. 게임룰 값 JSON 타입 (#11) | **확인.** `"type":"boolean"`→`true`, `"type":"integer"`→`128`. 문자열로 오지 않음 |
| 4. 리스트 `[]` 수용 (#9) | **확인.** `{"bans":[]}`, `{"add":[]}` 모두 정상 |

TLS 검증은 `InsecureSkipVerify` 없이 했다. 서버가 loopback에만 바인드되어 있고
인증서는 Tailscale MagicDNS 이름으로 발급된 실제 Let's Encrypt 인증서라,
`tls.Config.ServerName`만 그 이름으로 지정해 SNI와 호스트명 검증을 통과시켰다.
`TEST_TLS_SERVER_NAME` 환경변수가 이 역할을 한다.

## 실서버에서 발견된 결함

트레이스(#16)를 켜지 않았으면 넷 다 발견 불가능했다. 서버가 200을 돌려준다고
해서 우리가 보낸 JSON을 의도대로 해석했다는 뜻이 아니기 때문이다.

### 알림 params는 positional 인자 배열이다 — 페이로드 핸들러 12개 전멸

`minecraft:notification/allowlist/added`가 `[{"name":"fi_xz",…}]`로 오는 걸 보고
처음에는 "변경 항목 리스트"로 읽고 `decodeEach`로 배열을 풀었다. **틀린 해석이다.**

`rpc.discover`가 돌려주는 OpenRPC 스키마가 확정했다:

```json
"minecraft:notification/allowlist/added":
  params: [ { "name": "player", "schema": {"$ref": "player"} } ]
```

파라미터는 **하나**이고 타입은 단일 `player` 객체다. 바깥 배열은 JSON-RPC의
**positional 인자 리스트**이고, 페이로드는 그 0번 원소다. 스키마상 모든 알림이
파라미터 1개(또는 0개)라 두 해석이 데이터로는 구분되지 않았다.

클라이언트는 params를 페이로드 타입으로 바로 언마샬하고 있었으므로 실패했다.
`decodeEach`로 고쳤던 9개뿐 아니라 **`players/joined`, `players/left`,
`server/status`까지 총 12개**가 같은 이유로 죽어 있었다.

**수정:** `decodeParam`이 positional 배열을 풀고 0번 원소를 페이로드로 디코딩한다.
12개 전부 같은 경로를 탄다.

두 가지 교훈이 남았다. 기존 단위 테스트가 **틀린 shape를 인코딩**하고 있어서 버그를
감추고 있었다 — 가짜 서버 테스트는 우리가 쓴 가정을 되풀이할 뿐이다. 그리고 와이어
관측만으로는 해석이 갈릴 수 있으므로, **스키마가 있으면 스키마를 봐야 한다.**

### `allow_flight/set`의 파라미터 이름이 틀렸다 — 버전 무관

`{"allowed":true}`가 `-32602 Invalid params`로 거절되고 `{"allow":true}`만 통과한다.
**1.21.9와 26.2 양쪽 모두** 동일하므로 버전 차이가 아니라 처음부터 버그였다.
`SetAllowFlight`는 한 번도 성공한 적이 없다.

스키마가 원인을 설명한다:

```json
"minecraft:serversettings/allow_flight/set":
  params: [ { "name": "allow", ... } ]
  result:   { "name": "allowed", ... }
```

`allowed`는 **결과**의 이름이다. 요청 필드에 결과 이름을 쓴 것이다.

부수적으로 확인된 것: 서버는 **미지의 필드를 거절한다**
(`{"definitely_not_a_field":…}` → -32602). 따라서 거절 여부로 파라미터 이름의
정오를 판정할 수 있다. 지금은 그 판정을 `TestWireRequestParameterNames`가
스키마와 대조해 자동으로 한다 — setter 36개 전부 일치.

### 게임룰 표기가 서버 버전에 따라 다르다 — 라이브러리 결함 아님

트레이스에 `minecraft:keep_inventory`, `minecraft:random_tick_speed`가 찍혀
처음에는 예제가 틀린 줄 알았으나, **버전 차이**다. 1.21.11에서 게임룰이 레지스트리로
옮겨지면서 전부 namespaced resource location + snake_case로 개명됐다. 일부는
"Edit Game Rules" 메뉴 이름을 따라 아예 다른 이름이 됐고, **값 의미가 반전된 것도
있다** — 단순 표기 변환으로 처리할 수 없다는 뜻이다.

| 서버 | 키 표기 |
|---|---|
| 1.21.9 – 1.21.10 | `keepInventory` |
| 1.21.11 이상 | `minecraft:keep_inventory` |

MSMP 자체가 1.21.9에 도입됐으므로 이 라이브러리의 지원 범위가 그 경계를
가로지른다. 어느 한쪽을 정답으로 문서에 박으면 반대쪽 사용자가 틀리게 된다.

게다가 **값 표현도 다르다.** 스키마 대조 결과:

| | 1.21.9 – 1.21.10 | 26.2 |
|---|---|---|
| 키 | `keepInventory` | `minecraft:keep_inventory` |
| `untyped_game_rule.value` | `{"type":"string"}` | `{"type":["boolean","integer"]}` |

구 버전은 `"true"`, `"3"`처럼 문자열로 보내고 신 버전은 네이티브 JSON 타입을 보낸다.
`BoolRule(key, true)`가 만드는 `{"value":true}`는 1.21.9에서 `-32602`로 거절된다.

이 사실이 `Bool()`/`Int()`의 문자열 수용 분기를 정당화한다. 실서버 검증 초기에
26.2만 보고 "문자열 경로는 방어용일 뿐"이라고 주석을 고쳤는데, 그게 오히려 틀렸다.

**대응:** 읽기는 접근자가 양쪽을 수용하므로 이미 버전 독립적이다. 쓰기는 버전 감지
대신 **서버가 쓴 표현을 미러링**하는 방식을 택했다:

```go
updated, err := client.UpdateGameRule(ctx, rule.WithBool(true))
```

`WithBool` / `WithInt` / `WithString`은 원본 rule의 표현을 따라간다
(`UsesStringValues()`가 어느 쪽인지 알려준다). 서버 버전을 조회해 분기하는 상태
저장 로직 없이, 값이 온 그대로 돌려보내면 되므로 항상 맞는다.

키는 라이브러리가 그대로 통과시키므로 코드 변경이 없다. README에 버전별 표를 넣고,
양쪽을 지원해야 하는 사용자는 하드코딩 대신 `GetGameRules`가 돌려주는 키를 쓰도록
안내했다.

### `difficulty` / `status_heartbeat_interval` — 서버 쪽 문제

둘 다 파라미터 이름이 맞고(스키마 대조 통과) 호출도 성공하지만 값이 바뀌지 않는다.
클라이언트 결함이 아니라는 근거:

- **난이도**: 유효 enum(`peaceful`/`easy`/`normal`/`hard`)만 수용하고 `HARD`나 오타는
  -32602로 거절한다. 서버가 값을 실제로 파싱해 대조하고 있다는 뜻이다. 그런데도
  `"easy"`를 유지한다. level.dat는 `locked: 0`, `difficulty: hard`인데 MSMP는
  `"easy"`(= `server.properties` 값)를 보고한다. **월드 상태가 아니라
  server.properties를 읽고 있는 것으로 보인다.** 콘솔의 `/difficulty hard`는 정상
  동작하므로 월드가 잠긴 것도 아니다. 1.21.9와 26.2 양쪽 동일.
- **하트비트**: `-1`까지 수용한다 — 값을 검증조차 하지 않는다. **26.2에서만** 값이
  적용되지 않으며, interval=2로 설정하고 7초를 기다려도 `server/status` 알림이 하나도
  오지 않는다(getter가 아니라 setter가 no-op이라는 증거). **1.21.9에서는 정상
  동작한다**(`set 2` → `reads back 2`). 26.2의 회귀로 보인다.

두 테스트는 "요청대로 바뀌었는가" 대신 **"setter 반환값이 서버의 실제 상태와
일치하는가"**를 단언하도록 바꿨다. 서버 정책과 무관하게 클라이언트 회귀는 잡힌다.
