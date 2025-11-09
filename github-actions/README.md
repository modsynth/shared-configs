# GitHub Actions Workflows

Modsynth 프로젝트를 위한 재사용 가능한 GitHub Actions 워크플로우 템플릿입니다.

## 워크플로우 목록

### 1. Go CI (`go-ci.yml`)

Go 프로젝트를 위한 CI 파이프라인입니다.

**기능:**
- 멀티 버전 테스트 (Go 1.21, 1.22)
- `go vet` 정적 분석
- `staticcheck` 린팅
- 테스트 실행 (race detector 포함)
- 코드 커버리지 업로드 (Codecov)
- 빌드 검증
- `golangci-lint` 린팅

**사용 방법:**

```bash
# 프로젝트 루트에 .github/workflows/ 디렉토리 생성
mkdir -p .github/workflows

# 워크플로우 복사
cp /path/to/shared-configs/github-actions/go-ci.yml .github/workflows/

# Git에 추가 및 커밋
git add .github/workflows/go-ci.yml
git commit -m "ci: add Go CI workflow"
git push
```

**필요 조건:**
- `go.mod` 파일이 프로젝트 루트에 존재
- 테스트 파일 (`*_test.go`) 존재

### 2. Node.js CI (`node-ci.yml`)

TypeScript/Node.js 프로젝트를 위한 CI 파이프라인입니다.

**기능:**
- 멀티 버전 테스트 (Node 18.x, 20.x, 21.x)
- ESLint 린팅
- TypeScript 타입 체크
- 테스트 실행 및 커버리지
- 빌드 검증
- npm 자동 배포 (태그 푸시 시)

**사용 방법:**

```bash
# 워크플로우 복사
cp /path/to/shared-configs/github-actions/node-ci.yml .github/workflows/

# package.json에 필수 스크립트 추가
{
  "scripts": {
    "lint": "eslint . --ext .ts,.tsx",
    "type-check": "tsc --noEmit",
    "test": "jest",
    "build": "tsc"
  }
}

# Git에 추가 및 커밋
git add .github/workflows/node-ci.yml package.json
git commit -m "ci: add Node.js CI workflow"
git push
```

**필요 조건:**
- `package.json` 파일 존재
- `tsconfig.json` 파일 존재 (TypeScript 프로젝트)
- `npm run lint`, `npm run type-check`, `npm test`, `npm run build` 스크립트 정의
- npm 배포를 위해 `NPM_TOKEN` secret 설정 필요

**npm 배포 설정:**

1. npm 토큰 발급: https://www.npmjs.com/settings/YOUR_USERNAME/tokens
2. GitHub Repository Settings → Secrets and variables → Actions
3. New repository secret: `NPM_TOKEN`

### 3. Release (`release.yml`)

태그 푸시 시 자동으로 GitHub Release를 생성합니다.

**기능:**
- 태그 기반 릴리스 생성
- 자동 변경 로그 생성
- 설치 가이드 자동 생성
- Pre-release 자동 감지 (alpha, beta, rc)

**사용 방법:**

```bash
# 워크플로우 복사
cp /path/to/shared-configs/github-actions/release.yml .github/workflows/

# Git에 추가 및 커밋
git add .github/workflows/release.yml
git commit -m "ci: add release workflow"
git push

# 릴리스 태그 생성 및 푸시
git tag v0.2.0
git push origin v0.2.0
```

**태그 명명 규칙:**
- `v0.1.0` - 정식 릴리스
- `v0.1.0-alpha.1` - Alpha 릴리스 (pre-release)
- `v0.1.0-beta.1` - Beta 릴리스 (pre-release)
- `v0.1.0-rc.1` - Release Candidate (pre-release)

## 프로젝트별 워크플로우 설정

### Go 프로젝트 (백엔드 모듈)

```bash
cd /path/to/your-go-module
mkdir -p .github/workflows
cp /path/to/shared-configs/github-actions/go-ci.yml .github/workflows/
cp /path/to/shared-configs/github-actions/release.yml .github/workflows/
git add .github/workflows/
git commit -m "ci: add GitHub Actions workflows"
git push
```

### TypeScript 프로젝트 (프론트엔드 모듈)

```bash
cd /path/to/your-ts-module
mkdir -p .github/workflows
cp /path/to/shared-configs/github-actions/node-ci.yml .github/workflows/
cp /path/to/shared-configs/github-actions/release.yml .github/workflows/
git add .github/workflows/
git commit -m "ci: add GitHub Actions workflows"
git push
```

## Codecov 설정 (선택사항)

코드 커버리지를 추적하려면 Codecov를 설정하세요:

1. https://codecov.io/ 방문
2. GitHub Organization 연동
3. 각 레포지토리 활성화
4. Codecov 토큰 확인 (public repo는 불필요)

## 배지 추가

README.md에 배지를 추가하여 빌드 상태를 표시하세요:

```markdown
# Your Module

![Go CI](https://github.com/modsynth/your-module/workflows/Go%20CI/badge.svg)
![Node.js CI](https://github.com/modsynth/your-module/workflows/Node.js%20CI/badge.svg)
[![codecov](https://codecov.io/gh/modsynth/your-module/branch/main/graph/badge.svg)](https://codecov.io/gh/modsynth/your-module)
```

## 고급 설정

### 조건부 실행

특정 경로의 파일만 변경되었을 때 워크플로우 실행:

```yaml
on:
  push:
    branches: [ main ]
    paths:
      - 'src/**'
      - 'go.mod'
      - 'go.sum'
```

### 매트릭스 확장

더 많은 버전 테스트:

```yaml
strategy:
  matrix:
    go-version: ['1.19', '1.20', '1.21', '1.22']
    os: [ubuntu-latest, windows-latest, macos-latest]
```

### 캐싱 최적화

의존성 캐싱으로 빌드 속도 향상:

```yaml
- name: Cache Go modules
  uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

## 문제 해결

### 워크플로우가 실행되지 않는 경우

1. `.github/workflows/` 디렉토리가 정확한지 확인
2. YAML 문법 오류 확인: https://www.yamllint.com/
3. 브랜치 이름 확인 (main vs master)

### 테스트 실패 시

1. 로컬에서 테스트 실행: `go test ./...` 또는 `npm test`
2. GitHub Actions 로그 확인
3. 환경 변수 설정 확인

### npm 배포 실패 시

1. `NPM_TOKEN` secret 설정 확인
2. npm 토큰 권한 확인 (Automation 또는 Publish)
3. package.json의 `name` 필드가 `@modsynth/...` 형식인지 확인

## 참고 자료

- [GitHub Actions 문서](https://docs.github.com/en/actions)
- [Codecov 문서](https://docs.codecov.io/)
- [Semantic Versioning](https://semver.org/)
- [npm Publishing](https://docs.npmjs.com/cli/v10/commands/npm-publish)
