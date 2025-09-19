import '@testing-library/jest-dom';

declare global {
    namespace jest {
        interface Matchers<R> {
            toBeInTheDocument(): R;
            toBeDisabled(): R;
            toBeEnabled(): R;
            toHaveClass(className: string): R;
            toHaveTextContent(text: string | RegExp): R;
            toHaveValue(value: string | string[] | number): R;
            toBeChecked(): R;
            toBePartiallyChecked(): R;
            toHaveDescription(text?: string | RegExp): R;
            toHaveAccessibleDescription(text?: string | RegExp): R;
            toHaveAccessibleName(text?: string | RegExp): R;
            toHaveFormValues(expectedValues: Record<string, unknown>): R;
            toHaveDisplayValue(value: string | RegExp | (string | RegExp)[]): R;
            toBeRequired(): R;
            toBeInvalid(): R;
            toBeValid(): R;
            toHaveFocus(): R;
            toHaveAttribute(attr: string, value?: string): R;
            toHaveStyle(css: string | Record<string, unknown>): R;
        }
    }
}
