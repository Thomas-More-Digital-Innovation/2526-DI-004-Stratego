/**
 * GoStrategy Authentication and Password Validation Rules
 * Aligned precisely with the Go backend:
 * - Minimum length: 8 characters
 * - Must contain at least one uppercase letter (A-Z)
 * - Must contain at least one lowercase letter (a-z)
 * - Must contain at least one number (0-9)
 * - Allowed characters: Alphanumeric, spaces, and !@#$%^&*()_+=-.
 * - Leading or trailing spaces are forbidden
 */

export const PASSWORD_RULES = {
    minLength: 8,
    requireUppercase: true,
    requireLowercase: true,
    requireNumber: true,
    forbiddenPattern: /[^a-zA-Z0-9!@#$%^&*()_+=\-\. ]/
};

export interface PasswordChecks {
    isMinLength: boolean;
    hasUppercase: boolean;
    hasLowercase: boolean;
    hasNumber: boolean;
    isValidFormat: boolean;
}

export function validatePasswordInput(password: string): PasswordChecks {
    const hasEdgeSpace = password.startsWith(" ") || password.endsWith(" ");
    return {
        isMinLength: password.length >= PASSWORD_RULES.minLength,
        hasUppercase: /[A-Z]/.test(password),
        hasLowercase: /[a-z]/.test(password),
        hasNumber: /[0-9]/.test(password),
        isValidFormat: !hasEdgeSpace && !PASSWORD_RULES.forbiddenPattern.test(password)
    };
}

export function isPasswordStrong(password: string): boolean {
    const checks = validatePasswordInput(password);
    return (
        checks.isMinLength &&
        checks.hasUppercase &&
        checks.hasLowercase &&
        checks.hasNumber &&
        checks.isValidFormat
    );
}

export function getPasswordStrengthScore(password: string): number {
    if (password.length === 0) return 0;
    const checks = validatePasswordInput(password);
    let score = 0;
    
    // 1. Length check
    if (checks.isMinLength) score++;
    // 2. Case check (both upper and lower)
    if (checks.hasUppercase && checks.hasLowercase) score++;
    // 3. Number check
    if (checks.hasNumber) score++;
    // 4. Clean formatting check
    if (checks.isValidFormat) score++;
    
    // Returns 0 to 4 score
    return score;
}
