import React, { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import { Provider } from 'react-redux';
import { BrowserRouter } from 'react-router-dom';
import { configureStore, PreloadedState } from '@reduxjs/toolkit';

// Example Redux store setup - customize for your app
interface AppStore {
  // Add your store slices here
  // e.g., auth: AuthState;
}

export const createMockStore = (preloadedState?: PreloadedState<AppStore>) => {
  return configureStore({
    reducer: {
      // Add your reducers here
      // e.g., auth: authReducer,
    },
    preloadedState,
  });
};

interface ExtendedRenderOptions extends Omit<RenderOptions, 'queries'> {
  preloadedState?: PreloadedState<AppStore>;
  store?: ReturnType<typeof createMockStore>;
  withRouter?: boolean;
  withRedux?: boolean;
}

export function renderWithProviders(
  ui: ReactElement,
  {
    preloadedState,
    store = createMockStore(preloadedState),
    withRouter = false,
    withRedux = false,
    ...renderOptions
  }: ExtendedRenderOptions = {}
) {
  let Wrapper = ({ children }: { children: React.ReactNode }) => <>{children}</>;

  if (withRedux) {
    const ReduxWrapper = ({ children }: { children: React.ReactNode }) => (
      <Provider store={store}>{children}</Provider>
    );
    Wrapper = ReduxWrapper;
  }

  if (withRouter) {
    const PreviousWrapper = Wrapper;
    Wrapper = ({ children }: { children: React.ReactNode }) => (
      <BrowserRouter>
        <PreviousWrapper>{children}</PreviousWrapper>
      </BrowserRouter>
    );
  }

  return { store, ...render(ui, { wrapper: Wrapper, ...renderOptions }) };
}

// Re-export everything from React Testing Library
export * from '@testing-library/react';
export { renderWithProviders as render };

// Utility functions for common test scenarios

export const waitForLoadingToFinish = async () => {
  const { waitForElementToBeRemoved, screen } = await import('@testing-library/react');
  await waitForElementToBeRemoved(() => screen.queryByText(/loading/i), {
    timeout: 3000,
  });
};

export const mockApiResponse = <T,>(data: T, delay = 0): Promise<T> => {
  return new Promise((resolve) => {
    setTimeout(() => resolve(data), delay);
  });
};

export const mockApiError = (message: string, delay = 0): Promise<never> => {
  return new Promise((_, reject) => {
    setTimeout(() => reject(new Error(message)), delay);
  });
};

// User event helpers
export const createMockUser = () => ({
  id: '1',
  email: 'test@example.com',
  name: 'Test User',
  role: 'user',
});

export const createMockAuthState = (isAuthenticated = true) => ({
  user: isAuthenticated ? createMockUser() : null,
  token: isAuthenticated ? 'mock-token' : null,
  isLoading: false,
  error: null,
});

// Form helpers
export const fillInput = async (name: string, value: string) => {
  const { screen, userEvent } = await import('@testing-library/react');
  const input = screen.getByLabelText(name);
  await userEvent.clear(input);
  await userEvent.type(input, value);
};

export const submitForm = async () => {
  const { screen, userEvent } = await import('@testing-library/react');
  const submitButton = screen.getByRole('button', { name: /submit/i });
  await userEvent.click(submitButton);
};

// Async helpers
export const flushPromises = () => {
  return new Promise((resolve) => setImmediate(resolve));
};

// Local storage helpers
export const mockLocalStorage = () => {
  const store: Record<string, string> = {};

  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => {
      store[key] = value;
    },
    removeItem: (key: string) => {
      delete store[key];
    },
    clear: () => {
      Object.keys(store).forEach((key) => delete store[key]);
    },
  };
};

// WebSocket mock
export class MockWebSocket {
  url: string;
  readyState = WebSocket.CONNECTING;
  onopen: ((event: Event) => void) | null = null;
  onclose: ((event: CloseEvent) => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((event: Event) => void) | null = null;

  constructor(url: string) {
    this.url = url;
    setTimeout(() => {
      this.readyState = WebSocket.OPEN;
      this.onopen?.(new Event('open'));
    }, 0);
  }

  send(data: string) {
    // Mock send
  }

  close() {
    this.readyState = WebSocket.CLOSED;
    this.onclose?.(new CloseEvent('close'));
  }

  simulateMessage(data: any) {
    this.onmessage?.(new MessageEvent('message', { data: JSON.stringify(data) }));
  }
}

// API mock helpers
export const createMockFetch = (responses: Record<string, any>) => {
  return vi.fn((url: string) => {
    const response = responses[url];
    if (!response) {
      return Promise.reject(new Error(`No mock for URL: ${url}`));
    }

    return Promise.resolve({
      ok: true,
      json: () => Promise.resolve(response),
      text: () => Promise.resolve(JSON.stringify(response)),
    });
  });
};

// Intersection Observer mock
export const createMockIntersectionObserver = () => {
  const observers = new Map();

  return class {
    constructor(callback: IntersectionObserverCallback) {
      observers.set(this, callback);
    }

    observe(element: Element) {
      const callback = observers.get(this);
      if (callback) {
        callback(
          [
            {
              target: element,
              isIntersecting: true,
              intersectionRatio: 1,
            } as IntersectionObserverEntry,
          ],
          this as any
        );
      }
    }

    disconnect() {
      observers.delete(this);
    }

    unobserve() {
      // Mock unobserve
    }
  };
};
