# Testing Framework

Modsynth 프로젝트를 위한 통합 테스트 프레임워크입니다.

## 개요

이 디렉토리는 Go 및 TypeScript/React 프로젝트를 위한 테스트 유틸리티와 설정을 제공합니다:

- **Go 테스트**: Testcontainers를 사용한 통합 테스트
- **React 테스트**: React Testing Library + Jest/Vitest
- **E2E 테스트**: Playwright 설정

## Go 테스트

### 설치

```bash
go get github.com/testcontainers/testcontainers-go
go get github.com/stretchr/testify
```

### 사용법

#### PostgreSQL 테스트

```go
package service_test

import (
    "testing"
    "github.com/modsynth/shared-configs/testing"
)

func TestUserService(t *testing.T) {
    // PostgreSQL 컨테이너 시작
    postgres := testing.SetupPostgres(t)

    // 테이블 마이그레이션
    postgres.DB.AutoMigrate(&User{})

    // 서비스 테스트
    service := NewUserService(postgres.DB)
    user, err := service.Create("test@example.com")

    testing.AssertNoError(t, err)
    testing.AssertNotEqual(t, user.ID, "")

    // Cleanup은 자동으로 처리됨 (t.Cleanup)
}
```

#### Redis 테스트

```go
func TestCacheService(t *testing.T) {
    // Redis 컨테이너 시작
    redis := testing.SetupRedis(t)

    // 서비스 테스트
    service := NewCacheService(redis.Client)
    err := service.Set("key", "value")

    testing.AssertNoError(t, err)

    value, err := service.Get("key")
    testing.AssertNoError(t, err)
    testing.AssertEqual(t, value, "value")
}
```

#### 트랜잭션 테스트

```go
func TestWithTransaction(t *testing.T) {
    postgres := testing.SetupPostgres(t)
    postgres.DB.AutoMigrate(&User{})

    testing.RunInTransaction(t, postgres.DB, func(tx *gorm.DB) {
        // 트랜잭션 내에서 테스트
        user := &User{Email: "test@example.com"}
        tx.Create(user)

        // 여기서 실패하면 자동 롤백
        testing.AssertNotEqual(t, user.ID, uint(0))
    })

    // 트랜잭션이 롤백되었는지 확인
    var count int64
    postgres.DB.Model(&User{}).Count(&count)
    testing.AssertEqual(t, count, int64(0))
}
```

#### 헬퍼 함수

```go
// 에러 검증
testing.AssertNoError(t, err)
testing.AssertError(t, err)

// 값 비교
testing.AssertEqual(t, got, want)
testing.AssertNotEqual(t, got, want)

// 조건 검증
testing.AssertTrue(t, condition, "message")
testing.AssertFalse(t, condition, "message")

// 타이밍
testing.WaitFor(t, 5*time.Second, func() bool {
    return service.IsReady()
})

// 데이터 정리
testing.TruncateTables(t, db, "users", "posts")
testing.FlushRedis(t, client)
```

### 통합 테스트 예제

```go
func TestUserRegistrationFlow(t *testing.T) {
    // 인프라 설정
    postgres := testing.SetupPostgres(t)
    redis := testing.SetupRedis(t)

    // 마이그레이션
    postgres.DB.AutoMigrate(&User{}, &Session{})

    // 서비스 초기화
    userService := NewUserService(postgres.DB)
    authService := NewAuthService(postgres.DB, redis.Client)

    // 1. 사용자 등록
    user, err := userService.Register("test@example.com", "password")
    testing.AssertNoError(t, err)
    testing.AssertNotEqual(t, user.ID, uint(0))

    // 2. 로그인
    token, err := authService.Login("test@example.com", "password")
    testing.AssertNoError(t, err)
    testing.AssertNotEqual(t, token, "")

    // 3. 토큰 검증
    valid, err := authService.ValidateToken(token)
    testing.AssertNoError(t, err)
    testing.AssertTrue(t, valid, "token should be valid")

    // 4. 로그아웃
    err = authService.Logout(token)
    testing.AssertNoError(t, err)

    // 5. 토큰 무효화 확인
    valid, err = authService.ValidateToken(token)
    testing.AssertNoError(t, err)
    testing.AssertFalse(t, valid, "token should be invalid after logout")
}
```

---

## TypeScript/React 테스트

### 설치

```bash
npm install --save-dev \
  @testing-library/react \
  @testing-library/jest-dom \
  @testing-library/user-event \
  vitest \
  jsdom
```

### 설정

**package.json:**
```json
{
  "scripts": {
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage"
  }
}
```

**vitest.config.ts:**
```typescript
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/setupTests.ts',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/setupTests.ts',
      ],
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
});
```

### 사용법

#### 컴포넌트 테스트

```typescript
import { render, screen, userEvent } from '@/test-utils';
import { Button } from './Button';

describe('Button', () => {
  it('renders with text', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick when clicked', async () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click me</Button>);

    await userEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('disables when loading', () => {
    render(<Button loading>Click me</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });
});
```

#### Redux 연동 테스트

```typescript
import { render, screen } from '@/test-utils';
import { UserProfile } from './UserProfile';

describe('UserProfile', () => {
  it('displays user information', () => {
    const mockUser = {
      id: '1',
      name: 'Test User',
      email: 'test@example.com',
    };

    render(<UserProfile />, {
      withRedux: true,
      preloadedState: {
        auth: {
          user: mockUser,
          token: 'mock-token',
        },
      },
    });

    expect(screen.getByText('Test User')).toBeInTheDocument();
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
  });
});
```

#### Router 연동 테스트

```typescript
import { render, screen, userEvent } from '@/test-utils';
import { Navigation } from './Navigation';

describe('Navigation', () => {
  it('navigates to profile page', async () => {
    render(<Navigation />, { withRouter: true });

    await userEvent.click(screen.getByText('Profile'));
    expect(window.location.pathname).toBe('/profile');
  });
});
```

#### API Mock 테스트

```typescript
import { render, screen, waitFor } from '@/test-utils';
import { UserList } from './UserList';
import { mockApiResponse } from '@/test-utils';

describe('UserList', () => {
  it('loads and displays users', async () => {
    const mockUsers = [
      { id: '1', name: 'User 1' },
      { id: '2', name: 'User 2' },
    ];

    // API 모킹
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockUsers),
      } as Response)
    );

    render(<UserList />);

    // 로딩 상태 확인
    expect(screen.getByText(/loading/i)).toBeInTheDocument();

    // 데이터 로드 대기
    await waitFor(() => {
      expect(screen.getByText('User 1')).toBeInTheDocument();
      expect(screen.getByText('User 2')).toBeInTheDocument();
    });
  });
});
```

#### 폼 테스트

```typescript
import { render, screen, userEvent } from '@/test-utils';
import { fillInput, submitForm } from '@/test-utils';
import { LoginForm } from './LoginForm';

describe('LoginForm', () => {
  it('submits form with valid data', async () => {
    const handleSubmit = vi.fn();
    render(<LoginForm onSubmit={handleSubmit} />);

    await fillInput('Email', 'test@example.com');
    await fillInput('Password', 'password123');
    await submitForm();

    expect(handleSubmit).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'password123',
    });
  });

  it('shows validation errors', async () => {
    render(<LoginForm onSubmit={vi.fn()} />);

    await submitForm();

    expect(screen.getByText(/email is required/i)).toBeInTheDocument();
    expect(screen.getByText(/password is required/i)).toBeInTheDocument();
  });
});
```

#### WebSocket 테스트

```typescript
import { render, screen } from '@/test-utils';
import { MockWebSocket } from '@/test-utils';
import { ChatRoom } from './ChatRoom';

describe('ChatRoom', () => {
  it('receives and displays messages', () => {
    const mockWs = new MockWebSocket('ws://localhost:8080');
    global.WebSocket = vi.fn(() => mockWs) as any;

    render(<ChatRoom roomId="123" />);

    // 메시지 시뮬레이션
    mockWs.simulateMessage({
      type: 'MESSAGE',
      payload: {
        id: '1',
        text: 'Hello World',
        user: 'Test User',
      },
    });

    expect(screen.getByText('Hello World')).toBeInTheDocument();
  });
});
```

### 커버리지

```bash
# 커버리지 리포트 생성
npm run test:coverage

# 커버리지 확인
open coverage/index.html
```

---

## E2E 테스트 (Playwright)

### 설치

```bash
npm install --save-dev @playwright/test
npx playwright install
```

### 설정

**playwright.config.ts:**
```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !process.env.CI,
  },
});
```

### 예제

```typescript
import { test, expect } from '@playwright/test';

test('user can login', async ({ page }) => {
  await page.goto('/login');

  await page.fill('input[name="email"]', 'test@example.com');
  await page.fill('input[name="password"]', 'password123');
  await page.click('button[type="submit"]');

  await expect(page).toHaveURL('/dashboard');
  await expect(page.locator('h1')).toContainText('Dashboard');
});

test('user can create a task', async ({ page }) => {
  await page.goto('/projects/123');

  await page.click('button:has-text("New Task")');
  await page.fill('input[name="title"]', 'Test Task');
  await page.fill('textarea[name="description"]', 'Test Description');
  await page.click('button:has-text("Create")');

  await expect(page.locator('text=Test Task')).toBeVisible();
});
```

---

## 모범 사례

### 1. 테스트 구조
```typescript
describe('ComponentName', () => {
  describe('when condition', () => {
    it('should do something', () => {
      // Arrange
      const props = { ... };

      // Act
      render(<Component {...props} />);

      // Assert
      expect(...).toBe(...);
    });
  });
});
```

### 2. 테스트 격리
- 각 테스트는 독립적이어야 함
- 테스트 간 상태 공유 금지
- `beforeEach`/`afterEach` 사용

### 3. 의미 있는 테스트
- 구현이 아닌 동작 테스트
- 사용자 관점에서 테스트
- Edge case 고려

### 4. 성능
- 불필요한 렌더링 최소화
- Mock 적절히 사용
- 통합 테스트 수 제한

## 참고 자료

- [Testing Library](https://testing-library.com/)
- [Vitest](https://vitest.dev/)
- [Playwright](https://playwright.dev/)
- [Testcontainers](https://www.testcontainers.org/)
