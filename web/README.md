# Spark Park Cricket - Frontend

A modern Next.js frontend for the Spark Park Cricket tournament management system, built with TypeScript, Redux Toolkit, and Redux Saga.

## 🚀 Tech Stack

- **Framework**: Next.js 15 with App Router
- **Language**: TypeScript with strict configuration
- **State Management**: Redux Toolkit + Redux Saga
- **Styling**: Tailwind CSS
- **Testing**: Jest, React Testing Library, Cypress
- **Code Quality**: ESLint, Prettier, Husky

## 📁 Project Structure

```
web/
├── src/
│   ├── app/                    # Next.js App Router pages
│   ├── components/             # Reusable React components
│   ├── store/                  # Redux store configuration
│   │   ├── reducers/           # Redux slices
│   │   ├── sagas/              # Redux Saga middleware
│   │   └── hooks.ts            # Typed Redux hooks
│   ├── providers/              # React context providers
│   └── types/                  # TypeScript type definitions
├── cypress/                    # E2E tests
├── .cursor/rules/              # Cursor IDE rules
└── __tests__/                  # Test utilities and setup
```

## 🛠️ Development Setup

### Prerequisites

- Node.js 20.17+
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The application will be available at `http://localhost:3000`.

## 🧪 Testing

### Unit & Integration Tests

```bash
# Run all tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage

# Run tests for CI
npm run test:ci
```

### End-to-End Tests

```bash
# Run E2E tests headlessly
npm run e2e

# Open Cypress test runner
npm run e2e:open
```

## 🔧 Code Quality

### Linting & Formatting

```bash
# Check for linting errors
npm run lint

# Fix linting errors
npm run lint:fix

# Format code
npm run format

# Check formatting
npm run format:check
```

### Type Checking

```bash
# Run TypeScript type checking
npm run type-check
```

## 📋 Available Scripts

| Script          | Description                             |
| --------------- | --------------------------------------- |
| `dev`           | Start development server with Turbopack |
| `build`         | Build production bundle                 |
| `start`         | Start production server                 |
| `lint`          | Run ESLint                              |
| `lint:fix`      | Fix ESLint errors                       |
| `format`        | Format code with Prettier               |
| `test`          | Run Jest tests                          |
| `test:watch`    | Run tests in watch mode                 |
| `test:coverage` | Run tests with coverage report          |
| `e2e`           | Run Cypress E2E tests                   |
| `e2e:open`      | Open Cypress test runner                |
| `type-check`    | Run TypeScript compiler check           |

## 🏗️ Architecture

### State Management

The application uses Redux Toolkit for state management with Redux Saga for side effects:

- **Reducers**: Pure functions that specify how state changes
- **Sagas**: Handle async operations and side effects
- **Slices**: Combine reducers and actions using Redux Toolkit

### Component Structure

- **Pages**: Next.js App Router pages
- **Components**: Reusable UI components
- **Providers**: React context providers (Redux, Theme, etc.)

### API Integration

- RESTful API calls using fetch
- WebSocket connections for real-time updates
- Error handling and loading states

## 🎯 Features

- **Series Management**: Create and manage cricket series/tournaments
- **Match Management**: Schedule and track matches
- **Team Management**: Manage teams and players
- **Live Scoring**: Real-time scoreboard updates
- **Responsive Design**: Mobile-first responsive UI

## 🔒 Code Quality Standards

### TypeScript

- Strict type checking enabled
- No `any` types allowed
- Explicit return types for functions
- Proper interface definitions

### Testing Requirements

- **Unit Tests**: 90% code coverage minimum
- **Integration Tests**: Component and Redux flow testing
- **E2E Tests**: Complete user workflow testing

### Code Style

- ESLint with TypeScript rules
- Prettier for consistent formatting
- Husky pre-commit hooks
- Conventional commit messages

## 🚀 Deployment

### Production Build

```bash
# Build for production
npm run build

# Start production server
npm start
```

### Environment Variables

Create a `.env.local` file with:

```env
# Backend API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_GRAPHQL_URL=http://localhost:8080/api/v1/graphql
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws

# Authentication Configuration (if needed)
# NEXT_PUBLIC_AUTH_DOMAIN=your-auth-domain
# NEXT_PUBLIC_AUTH_CLIENT_ID=your-client-id
```

**Note**: Copy `.env.example` to `.env.local` and modify the values as needed for your environment.

## 📚 Documentation

- [Next.js Documentation](https://nextjs.org/docs)
- [Redux Toolkit Documentation](https://redux-toolkit.js.org/)
- [Redux Saga Documentation](https://redux-saga.js.org/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [Cypress Documentation](https://docs.cypress.io/)

## 🤝 Contributing

1. Follow the established code style and patterns
2. Write tests for all new functionality
3. Ensure all tests pass before submitting PR
4. Update documentation as needed

## 📄 License

This project is part of the Spark Park Cricket system.
