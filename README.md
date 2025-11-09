# Modsynth Shared Configurations

> Shared ESLint, Prettier, and TypeScript configurations for all Modsynth projects

## Overview

This repository contains standardized configuration files for code quality and consistency across the Modsynth ecosystem. All modules should use these configurations to maintain a uniform development experience.

## Contents

- **ESLint Configurations** - Code linting for JavaScript/TypeScript
- **Prettier Configurations** - Code formatting
- **TypeScript Configurations** - TypeScript compiler options

## Installation

### For TypeScript/Node.js Projects

```bash
npm install --save-dev @modsynth/shared-configs eslint prettier typescript
```

### For React Projects

```bash
npm install --save-dev @modsynth/shared-configs eslint prettier typescript \
  eslint-plugin-react eslint-plugin-react-hooks eslint-plugin-jsx-a11y
```

## Usage

### ESLint

#### For Node.js/Backend Projects

Create `.eslintrc.json` in your project root:

```json
{
  "extends": ["./node_modules/@modsynth/shared-configs/eslint/.eslintrc.node.json"],
  "parserOptions": {
    "project": "./tsconfig.json"
  }
}
```

#### For React/Frontend Projects

Create `.eslintrc.json` in your project root:

```json
{
  "extends": ["./node_modules/@modsynth/shared-configs/eslint/.eslintrc.react.json"],
  "parserOptions": {
    "project": "./tsconfig.json"
  }
}
```

#### For Base Projects (Generic TypeScript)

Create `.eslintrc.json` in your project root:

```json
{
  "extends": ["./node_modules/@modsynth/shared-configs/eslint/.eslintrc.base.json"],
  "parserOptions": {
    "project": "./tsconfig.json"
  }
}
```

### Prettier

Create `.prettierrc.json` in your project root:

```json
"@modsynth/shared-configs/prettier/.prettierrc.json"
```

Or extend and customize:

```json
{
  "...": "@modsynth/shared-configs/prettier/.prettierrc.json",
  "printWidth": 120
}
```

Copy `.prettierignore`:

```bash
cp node_modules/@modsynth/shared-configs/prettier/.prettierignore .prettierignore
```

### TypeScript

#### For Node.js/Backend Projects

Create `tsconfig.json` in your project root:

```json
{
  "extends": "@modsynth/shared-configs/typescript/tsconfig.node.json",
  "compilerOptions": {
    "outDir": "./dist",
    "rootDir": "./src"
  },
  "include": ["src/**/*"]
}
```

#### For React/Frontend Projects

Create `tsconfig.json` in your project root:

```json
{
  "extends": "@modsynth/shared-configs/typescript/tsconfig.react.json",
  "compilerOptions": {
    "baseUrl": "./src"
  },
  "include": ["src/**/*"]
}
```

## Configuration Details

### ESLint Rules

All configurations include:

- TypeScript strict type checking
- Unused variable detection (with `_` prefix exemption)
- Import ordering and grouping
- Console warning (allow only `warn` and `error`)
- No debugger statements

**React-specific additions:**
- React Hooks rules enforcement
- JSX accessibility checks
- Automatic React version detection

**Node-specific adjustments:**
- Allow `process.exit()`
- Allow CommonJS `require()`

### Prettier Rules

- **Semi**: Always use semicolons
- **Quotes**: Single quotes for strings
- **Print Width**: 100 characters
- **Tab Width**: 2 spaces
- **Trailing Commas**: ES5 style
- **Arrow Parens**: Always
- **End of Line**: LF (Unix style)

### TypeScript Rules

All configurations use:

- **Strict Mode**: Enabled
- **ES2022**: Target and lib
- **Source Maps**: Enabled
- **Declaration Files**: Generated
- **No Unused Locals/Parameters**: Enforced
- **No Implicit Returns**: Required

**React-specific:**
- JSX: `react-jsx` (no need to import React)
- Module Resolution: `bundler`
- Vite types included

**Node-specific:**
- Module: `CommonJS`
- Node types included

## Package Scripts

Add these scripts to your `package.json`:

```json
{
  "scripts": {
    "lint": "eslint . --ext .ts,.tsx",
    "lint:fix": "eslint . --ext .ts,.tsx --fix",
    "format": "prettier --write \"src/**/*.{ts,tsx,json,md}\"",
    "format:check": "prettier --check \"src/**/*.{ts,tsx,json,md}\"",
    "type-check": "tsc --noEmit"
  }
}
```

## VS Code Integration

Create `.vscode/settings.json` in your project:

```json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "eslint.validate": [
    "javascript",
    "javascriptreact",
    "typescript",
    "typescriptreact"
  ]
}
```

Recommended VS Code extensions:
- `dbaeumer.vscode-eslint`
- `esbenp.prettier-vscode`

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Lint and Format Check

on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: npm ci
      - run: npm run lint
      - run: npm run format:check
      - run: npm run type-check
```

## Migration Guide

### Migrating Existing Projects

1. **Install shared configs**:
   ```bash
   npm install --save-dev @modsynth/shared-configs
   ```

2. **Remove old config files**:
   ```bash
   rm .eslintrc.* .prettierrc.* tsconfig.json
   ```

3. **Create new configs** following the usage examples above

4. **Run auto-fix**:
   ```bash
   npm run lint:fix
   npm run format
   ```

5. **Review and commit**:
   ```bash
   git diff
   git add .
   git commit -m "chore: migrate to shared-configs"
   ```

## Development

To modify these configurations:

1. Clone this repository
2. Make changes to the appropriate config files
3. Update version in `package.json`
4. Test in a sample project
5. Commit and tag with new version

```bash
git tag v0.2.0
git push origin v0.2.0
```

## Configuration Files

```
shared-configs/
├── eslint/
│   ├── .eslintrc.base.json      # Base TypeScript rules
│   ├── .eslintrc.react.json     # React-specific rules
│   └── .eslintrc.node.json      # Node.js-specific rules
├── prettier/
│   ├── .prettierrc.json         # Formatting rules
│   └── .prettierignore          # Files to ignore
├── typescript/
│   ├── tsconfig.base.json       # Base compiler options
│   ├── tsconfig.react.json      # React compiler options
│   └── tsconfig.node.json       # Node.js compiler options
├── package.json
└── README.md
```

## Version History

- **v0.1.0** - Initial release with ESLint, Prettier, and TypeScript configs

## Contributing

When proposing configuration changes:

1. Explain the reasoning
2. Show examples of what it fixes/improves
3. Consider backward compatibility
4. Update this README

## License

MIT

## Links

- [Modsynth Organization](https://github.com/modsynth)
- [ESLint Documentation](https://eslint.org/)
- [Prettier Documentation](https://prettier.io/)
- [TypeScript Documentation](https://www.typescriptlang.org/)
