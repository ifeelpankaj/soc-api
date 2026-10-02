export interface paths {
    "/health": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Health check
         * @description Checks if the API server is running.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description API server is healthy */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.HealthCheckResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/health/live": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Liveness check
         * @description Checks if the API process is alive.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description API process is alive */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.LivenessResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/health/ready": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Readiness check
         * @description Checks if the API server and database are ready.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description API and database are ready */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ReadinessResponseDoc"];
                    };
                };
                /** @description Database connection failed */
                503: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.HealthErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/change-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Change password
         * @description Changes the authenticated user's password after verifying the current password, invalidates all access and refresh sessions, then clears auth cookies.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Change password payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.ChangePasswordRequest"];
                };
            };
            responses: {
                /** @description Password changed successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ChangePasswordAPIResponse"];
                    };
                };
                /** @description CURRENT_PASSWORD_INCORRECT, PASSWORD_NOT_SET, PASSWORD_TOO_LONG, mismatch, reuse, or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/forgot-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Request password reset OTP
         * @description Always returns a generic success message. If the email exists, sends a password reset OTP.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Forgot password payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.ForgotPasswordRequest"];
                };
            };
            responses: {
                /** @description Password reset instructions response */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ForgotPasswordAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/login": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Login
         * @description Authenticates a verified active user and sets access_token and refresh_token HTTP-only cookies.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Login payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.LoginRequest"];
                };
            };
            responses: {
                /** @description Login successful */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.LoginAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invalid credentials */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Email not verified or account disabled */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/logout": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Logout
         * @description Clears access_token and refresh_token cookies. This route is public and succeeds even if cookies are absent.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Logout successful */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.LogoutAPIResponse"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/profile": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get current user profile
         * @description Returns the authenticated user's profile. Does not refresh or rotate tokens.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Profile fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.GetProfileAPIResponse"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Account disabled */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Update current user profile
         * @description Updates authenticated user profile fields such as name, phone, and date of birth.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Profile update payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.UpdateUserRequest"];
                };
            };
            responses: {
                /** @description Profile updated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.GetProfileAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Account disabled */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/auth/profile/avatar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * GetAvatar
         * @description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Server-selected image view */
                    variant?: "list" | "avatar" | "detail" | "original";
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ImageViewAPIResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Request Entity Too Large */
                413: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Bad Gateway */
                502: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Service Unavailable */
                503: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        /**
         * PutAvatar
         * @description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
         */
        put: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody: components["requestBodies"]["putV1AuthProfileAvatar"];
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ImageViewAPIResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Request Entity Too Large */
                413: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Bad Gateway */
                502: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Service Unavailable */
                503: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        post?: never;
        /**
         * DeleteAvatar
         * @description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
         */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Image reference removed */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Request Entity Too Large */
                413: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Bad Gateway */
                502: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Service Unavailable */
                503: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/refresh": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Refresh access token
         * @description Validates the refresh token from either the refresh_token cookie or the request body.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Refresh token payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.RefreshTokenRequest"];
                };
            };
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.RefreshTokenAPIResponse"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal Server Error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/register": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Register a new user
         * @description Creates a user account and sends an email verification OTP.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Registration payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.RegisterRequest"];
                };
            };
            responses: {
                /** @description Account created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.RegisterAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Email or phone already exists */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/resend-otp": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Resend email verification OTP
         * @description Sends a fresh email verification OTP for an unverified user.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Resend OTP payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.ResendOTPRequest"];
                };
            };
            responses: {
                /** @description OTP sent successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ResendOTPAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description User not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Email already verified */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/reset-password": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reset password
         * @description Resets a user's password using a valid password reset OTP and clears auth cookies on success.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Reset password payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.ResetPasswordRequest"];
                };
            };
            responses: {
                /** @description Password reset successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ResetPasswordAPIResponse"];
                    };
                };
                /** @description Invalid OTP, password reuse, or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description User or active OTP not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Too many OTP attempts */
                429: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/resident/register": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Register resident without OTP
         * @description Creates a resident user for the public flat claim flow, marks the email verified, and sets auth cookies immediately. This does not create a society membership or flat resident.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Resident registration payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.ResidentRegisterRequest"];
                };
            };
            responses: {
                /** @description Resident account created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.LoginAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Email or phone already exists */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/verify-otp": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Verify email OTP
         * @description Verifies the registration email OTP and marks the user's email as verified.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Email verification payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.VerifyOTPRequest"];
                };
            };
            responses: {
                /** @description Email verified successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VerifyOTPAPIResponse"];
                    };
                };
                /** @description Invalid OTP or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description User or OTP not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Too many OTP attempts */
                429: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/bootstrap": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Bootstrap authenticated session
         * @description Returns the authenticated user, memberships, residences, and the default dashboard destination.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Bootstrap fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.BootstrapAPIResponse"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Account disabled or blocked */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/developer/dashboard/bootstrap": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get developer dashboard bootstrap
         * @description [Developer] Returns platform summary data for the developer dashboard.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Developer dashboard bootstrap fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.DeveloperDashboardBootstrapAPIResponse"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Developer access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/flat-claims": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List flat claims
         * @description [Owner/Admin/Staff/Developer] Lists flat claims with flexible filters and rich joined data.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Claim ID */
                    id?: number;
                    /** @description Society ID */
                    society_id?: number;
                    /** @description Flat ID */
                    flat_id?: number;
                    /** @description User ID */
                    user_id?: number;
                    /** @description Claim status: pending, approved, rejected, cancelled */
                    status?: string;
                    /** @description Search user, contact, flat, block, or status */
                    search?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat claims fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatClaimsAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        /**
         * Submit flat claim
         * @description [User] Submits a pending claim for a flat.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Submit flat claim payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.SubmitFlatClaimRequest"];
                };
            };
            responses: {
                /** @description Flat claim submitted successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatClaimAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Duplicate pending claim or flat unavailable */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/flat-claims/{claimId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat claim
         * @description [User/Developer] Fetches one flat claim with joined user, flat, and society data.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Claim ID */
                    claimId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat claim fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatClaimAPIResponse"];
                    };
                };
                /** @description Invalid claim ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat claim not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/flat-claims/{claimId}/cancel": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Cancel my flat claim
         * @description [User] Cancels the authenticated user's pending flat claim.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Claim ID */
                    claimId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat claim cancelled successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatClaimAPIResponse"];
                    };
                };
                /** @description Invalid claim ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat claim not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invalid claim transition */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/flat-residents": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List flat residents
         * @description [Owner/Admin/Staff/Developer] Lists flat residents with flexible filters and joined user, flat, and society data.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Resident ID */
                    id?: number;
                    /** @description Society ID */
                    society_id?: number;
                    /** @description Flat ID */
                    flat_id?: number;
                    /** @description User ID */
                    user_id?: number;
                    /** @description Resident role: owner, tenant, family */
                    role?: string;
                    /** @description Resident status: active, inactive, moved_out */
                    status?: string;
                    /** @description Primary resident flag */
                    is_primary?: boolean;
                    /** @description Search user, contact, flat, block, role, or status */
                    search?: string;
                    /** @description Search mode */
                    search_mode?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Residents fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatResidentsAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/flats": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List flats
         * @description [User/Developer] Lists flats with flexible admin/developer filters.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Flat ID */
                    id?: number;
                    /** @description Society ID */
                    society_id?: number;
                    /** @description Block */
                    block?: string;
                    /** @description Floor */
                    floor?: string;
                    /** @description Flat number */
                    flat_number?: string;
                    /** @description Flat status: vacant, occupied, blocked */
                    status?: string;
                    /** @description Active state */
                    is_active?: boolean;
                    /** @description Search flat number, block, floor, status, society name/code */
                    search?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flats fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatsAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/me/device-tokens": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Register or refresh a device push token
         * @description Stores an FCM token. Web requires an installation device_id and active resident access; the server binds it to the current authenticated session version. Unbound/revoked web tokens never receive delivery.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Device token payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.RegisterDeviceTokenRequest"];
                };
            };
            responses: {
                /** @description Device token registered successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.RegisterDeviceTokenAPIResponse"];
                    };
                };
                /** @description Invalid device token request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        /**
         * Unregister a device push token
         * @description Removes the authenticated user's FCM device token, typically on logout.
         */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Device token payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.UnregisterDeviceTokenRequest"];
                };
            };
            responses: {
                /** @description Device token removed successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MessageAPIResponse"];
                    };
                };
                /** @description Invalid device token request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/me/flat-claims": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List my flat claims
         * @description [User] Lists flat claims submitted by the authenticated user.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Society ID */
                    society_id?: number;
                    /** @description Flat ID */
                    flat_id?: number;
                    /** @description Claim status */
                    status?: string;
                    /** @description Search text */
                    search?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description My flat claims fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatClaimsAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/me/notifications": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List my notifications
         * @description Lists durable notifications for the authenticated user.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Page size */
                    limit?: number;
                    /** @description Opaque pagination cursor */
                    cursor?: string;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Notifications fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.NotificationsAPIResponse"];
                    };
                };
                /** @description Invalid notification request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/me/notifications/read-all": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Mark all notifications read
         * @description Marks all notifications for the authenticated user as read.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Notifications marked read */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.NotificationUnreadCountAPIResponse"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/me/notifications/unread-count": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get notification unread count
         * @description Returns the authenticated user's unread notification count.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Unread count fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.NotificationUnreadCountAPIResponse"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/me/notifications/{notificationId}/read": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Mark notification read
         * @description Marks one owned notification as read.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Notification ID */
                    notificationId: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Notification marked read */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.NotificationReadAPIResponse"];
                    };
                };
                /** @description Invalid notification request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Notification not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/me/residences": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List my residences
         * @description [User] Lists flat residences for the authenticated user with joined flat and society data.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Society ID */
                    society_id?: number;
                    /** @description Flat ID */
                    flat_id?: number;
                    /** @description Resident role */
                    role?: string;
                    /** @description Resident status */
                    status?: string;
                    /** @description Primary resident flag */
                    is_primary?: boolean;
                    /** @description Search text */
                    search?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description My residences fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MyResidencesAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/plans": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List plans
         * @description [Developer/User] Lists plans with flexible filters.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Plan code */
                    code?: string;
                    /** @description Billing cycle */
                    billing_cycle?: string;
                    /** @description Active state */
                    is_active?: boolean;
                    /** @description Search text */
                    search?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Plans fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PlansAPIResponse"];
                    };
                };
            };
        };
        put?: never;
        /**
         * Create plan
         * @description [Developer] Creates a subscription plan.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Create plan payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.CreatePlanRequest"];
                };
            };
            responses: {
                /** @description Plan created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PlanAPIResponse"];
                    };
                };
                /** @description Invalid request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Duplicate plan */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/plans/lookup": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get plan
         * @description [Developer/User] Fetches a plan by flexible filter.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Plan ID */
                    id?: number;
                    /** @description Plan code */
                    code?: string;
                    /** @description Plan name */
                    name?: string;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Plan fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PlanAPIResponse"];
                    };
                };
                /** @description Plan not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/plans/{planId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Update plan
         * @description [Developer] Updates a subscription plan.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Plan ID */
                    planId: number;
                };
                cookie?: never;
            };
            /** @description Update plan payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.UpdatePlanRequest"];
                };
            };
            responses: {
                /** @description Plan updated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PlanAPIResponse"];
                    };
                };
                /** @description Invalid request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Plan not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/plans/{planId}/activate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Activate plan
         * @description [Developer] Activates a subscription plan.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Plan ID */
                    planId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Plan activated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PlanAPIResponse"];
                    };
                };
                /** @description Plan not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/plans/{planId}/deactivate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Deactivate plan
         * @description [Developer] Deactivates a subscription plan.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Plan ID */
                    planId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Plan deactivated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PlanAPIResponse"];
                    };
                };
                /** @description Plan not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/flat-member-invites/{token}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat member invite by token
         * @description [Public] Fetches an active flat member invite from its public token before acceptance.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Member invite token */
                    token: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member invite fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PublicFlatMemberInviteAPIResponse"];
                    };
                };
                /** @description Invalid token */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite is expired, accepted, or cancelled */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/flat-member-invites/{token}/accept": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Accept flat member invite
         * @description [User] Accepts a flat member invite using the public token. The authenticated user becomes an active flat resident.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Member invite token */
                    token: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member invite accepted successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.AcceptFlatMemberInviteAPIResponse"];
                    };
                };
                /** @description Invalid token */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite unavailable or resident conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/flat-member-invites/{token}/join": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Join flat via member invite
         * @description [Public] Creates or authenticates a user and accepts a flat member invite without setting auth cookies.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Member invite token */
                    token: string;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.JoinFlatMemberInviteRequest"];
            responses: {
                /** @description Member joined successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.JoinFlatMemberInviteAPIResponse"];
                    };
                };
                /** @description Invalid request or token */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite unavailable or resident conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/short-links/{shortCode}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Resolve short link
         * @description [Public] Resolves a reusable Apna Gate short link into the current invite workflow state.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Short link code */
                    shortCode: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Short link resolved successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ShortLinkResolveAPIResponse"];
                    };
                };
                /** @description Invalid short code */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Short link not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Short link unavailable */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/short-links/{shortCode}/member-accept": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Accept member invite via short link
         * @description [User] Accepts a flat member invite resolved by short code.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Short link code */
                    shortCode: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member invite accepted successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.AcceptFlatMemberInviteAPIResponse"];
                    };
                };
                /** @description Invalid short code */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Short link or member invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite unavailable or resident conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/short-links/{shortCode}/member-join": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Join member invite via short link
         * @description [Public] Creates or authenticates a user and accepts a flat member invite resolved by short code.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Short link code */
                    shortCode: string;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.JoinFlatMemberInviteRequest"];
            responses: {
                /** @description Member joined successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.JoinFlatMemberInviteAPIResponse"];
                    };
                };
                /** @description Invalid request or short code */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Short link or member invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite unavailable or resident conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/short-links/{shortCode}/visitor-submit": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Submit visitor invite via short link
         * @description [Public] Submits visitor details for an active visitor invite resolved by short code.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Short link code */
                    shortCode: string;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.VisitorFormRequest"];
            responses: {
                /** @description Visitor entry created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryMutationAPIResponse"];
                    };
                };
                /** @description Invalid request or short code */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Short link or invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite unavailable */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/societies/{societyCode}/claim-options": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get public resident claim options
         * @description Public QR flow endpoint that resolves an active society code and returns safe society details plus active flats for resident claims.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society code */
                    societyCode: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Claim options fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PublicClaimOptionsAPIResponse"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/societies/{societyCode}/visitor-entries/public-qr": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Create public QR visitor entry
         * @description [Public] Creates a visitor entry from a society public QR flow. Approval may be required depending on visitor settings.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society code */
                    societyCode: string;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.VisitorFormRequest"];
            responses: {
                /** @description Visitor entry created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryMutationAPIResponse"];
                    };
                };
                /** @description Invalid request, society code, or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Public QR visitor entry is disabled */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society, flat, or visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry cannot be created in current state */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/societies/{societyCode}/visitor-entries/quick-link": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Create quick-link visitor entry
         * @description [Public] Creates a visitor entry from a quick-link flow. Approval may be required depending on visitor settings.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society code */
                    societyCode: string;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.VisitorFormRequest"];
            responses: {
                /** @description Visitor entry created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryMutationAPIResponse"];
                    };
                };
                /** @description Invalid request, society code, or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Quick-link visitor entry is disabled */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society, flat, or visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry cannot be created in current state */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/societies/{societyCode}/visitor-entry-options": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get visitor entry options
         * @description [Public] Fetches allowed visitor purposes and flat/block options for the public visitor entry form.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society code */
                    societyCode: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor entry options fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryOptionsAPIResponse"];
                    };
                };
                /** @description Invalid society code */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society or flats not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/societies/{societyCode}/visitor-invites/{token}/submit": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Submit visitor invite form
         * @description [Public] Submits visitor details for an active resident invite and creates the visitor entry.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society code */
                    societyCode: string;
                    /** @description Visitor invite token */
                    token: string;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.VisitorFormRequest"];
            responses: {
                /** @description Visitor entry created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryMutationAPIResponse"];
                    };
                };
                /** @description Invalid request, token, or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite is expired, used, or cancelled */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/visitor-entries/qr/validate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Validate visitor QR token
         * @description [Public] Validates a visitor QR token and returns the matching approved visitor entry.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.QRTokenRequest"];
            responses: {
                /** @description Visitor QR validated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryAPIResponse"];
                    };
                };
                /** @description Invalid request, token, or visitor entry state */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description QR token is expired, unavailable, or already used */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/public/visitor-invites/{token}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get visitor invite by token
         * @description [Public] Fetches a visitor invite for the web form, or recovers entry/QR when the invite was already used.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Visitor invite token */
                    token: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor invite fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorInviteAPIResponse"];
                    };
                };
                /** @description Invalid token */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite is expired, used without entry, or cancelled */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List societies
         * @description [Developer] Lists societies with flexible filters for admin/developer panels.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Society status */
                    status?: string;
                    /** @description Search text */
                    search?: string;
                    /** @description Search mode */
                    search_mode?: string;
                    /** @description Society name */
                    name?: string;
                    /** @description Society code */
                    code?: string;
                    /** @description City */
                    city?: string;
                    /** @description State */
                    state?: string;
                    /** @description Country */
                    country?: string;
                    /** @description Pincode */
                    pincode?: string;
                    /** @description Created by user ID */
                    created_by?: number;
                    /** @description Approved by user ID */
                    approved_by?: number;
                    /** @description Rejected by user ID */
                    rejected_by?: number;
                    /** @description Suspended by user ID */
                    suspended_by?: number;
                    /** @description Created from RFC3339 timestamp */
                    created_from?: string;
                    /** @description Created to RFC3339 timestamp */
                    created_to?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                    /** @description Sort by: created_at, updated_at, name, city, status */
                    sort_by?: string;
                    /** @description Sort order: asc, desc */
                    sort_order?: string;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Societies fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PaginatedSocietiesAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        /**
         * Create society request
         * @description [User] Creates a pending society request and assigns the requester as owner.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            /** @description Create society payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.CreateSocietyRequest"];
                };
            };
            responses: {
                /** @description Society request created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Duplicate society request */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/my": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List my societies
         * @description [User] Lists societies where the authenticated user has a membership.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description My societies fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MySocietiesAPIResponse"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get society
         * @description [User/Developer] Fetches one society by ID. Developer routes may use the same handler with broader guards.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Society fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyDetailAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /**
         * Soft delete society
         * @description [User] Soft-deletes a society. Requires active owner membership.
         */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Society deleted successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MessageAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /**
         * Update society
         * @description [User] Updates society profile fields. Requires active owner/admin membership.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Update society payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.UpdateSocietyRequest"];
                };
            };
            responses: {
                /** @description Society updated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/societies/{societyId}/allmember": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List all society members for developer
         * @description [Developer] Lists society members with filters without requiring owner/admin society membership.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Search full name, email, phone, role, or status */
                    search?: string;
                    /** @description Member role */
                    role?: string;
                    /** @description Member status */
                    status?: string;
                    /** @description User ID */
                    user_id?: number;
                    /** @description Invited by user ID */
                    invited_by?: number;
                    /** @description Removed by user ID */
                    removed_by?: number;
                    /** @description Joined from RFC3339 timestamp */
                    joined_from?: string;
                    /** @description Joined to RFC3339 timestamp */
                    joined_to?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                    /** @description Sort by: joined_at, role, status */
                    sort_by?: string;
                    /** @description Sort order: asc, desc */
                    sort_order?: string;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Members fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PaginatedMembersAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/approve": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Approve society
         * @description [Developer] Approves a pending society request. Transition allowed: pending -> active.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Society approved successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invalid society lifecycle transition */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/dashboard/bootstrap": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get society dashboard bootstrap
         * @description [Owner/Admin] Returns dashboard summary data for a society.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Society dashboard bootstrap fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyDashboardBootstrapAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flat-claims": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List society flat claims
         * @description [Owner/Admin/Staff] Lists flat claims inside one society with flexible filters and rich joined data.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Flat ID */
                    flat_id?: number;
                    /** @description User ID */
                    user_id?: number;
                    /** @description Claim status: pending, approved, rejected, cancelled */
                    status?: string;
                    /** @description Search user, contact, flat, block, or status */
                    search?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat claims fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatClaimsAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flat-claims/stats": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat claim stats
         * @description [Owner/Admin/Staff] Returns flat claim counts by status for a society.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat claim stats fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatClaimStatsAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flat-claims/{claimId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get society flat claim
         * @description [Owner/Admin/Staff] Fetches one flat claim inside a society with joined user, flat, and society data.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Claim ID */
                    claimId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat claim fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatClaimAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat claim not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flat-claims/{claimId}/approve": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Approve flat claim
         * @description [Owner/Admin/Staff] Transactionally approves a pending flat claim, activates society membership, creates flat resident, marks flat occupied, and returns rich joined data.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Claim ID */
                    claimId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat claim approved successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatApprovalAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat claim or flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat unavailable, duplicate resident, or primary resident conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flat-claims/{claimId}/reject": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reject flat claim
         * @description [Owner/Admin/Staff] Rejects a pending flat claim with a reason.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Claim ID */
                    claimId: number;
                };
                cookie?: never;
            };
            /** @description Reject flat claim payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.RejectFlatClaimRequest"];
                };
            };
            responses: {
                /** @description Flat claim rejected successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatClaimAPIResponse"];
                    };
                };
                /** @description Invalid request or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat claim not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invalid claim transition */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List society flats
         * @description [Owner/Admin/Staff] Lists flats inside one society with paginated metadata.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Block */
                    block?: string;
                    /** @description Floor */
                    floor?: string;
                    /** @description Flat number */
                    flat_number?: string;
                    /** @description Flat status: vacant, occupied, blocked */
                    status?: string;
                    /** @description Active state */
                    is_active?: boolean;
                    /** @description Search flat number, block, floor, or status */
                    search?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flats fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PaginatedFlatsAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        /**
         * Create flat
         * @description [Owner/Admin/Staff] Creates a flat inside a society.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Create flat payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.CreateFlatRequest"];
                };
            };
            responses: {
                /** @description Flat created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Duplicate flat */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/bulk": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Bulk create flats
         * @description [Owner/Admin/Staff] Creates multiple flats inside a society in one request.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Bulk create flats payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.BulkCreateFlatsRequest"];
                };
            };
            responses: {
                /** @description Flats created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.BulkFlatsAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Duplicate flat */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/generate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Generate and create flats
         * @description [Owner/Admin] Generates and persists a batch of flats atomically from a continuous or floor-based numbering pattern.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Generate flats payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.GenerateFlatsRequest"];
                };
            };
            responses: {
                /** @description Flats generated and created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.BulkFlatsAPIResponse"];
                    };
                };
                /** @description Invalid request or generation configuration */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Active subscription required or flat quota exceeded */
                402: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Duplicate flat */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/stats": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat stats
         * @description [Owner/Admin/Staff/Developer] Returns flat counts by status and active state for a society.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat stats fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatStatsAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat
         * @description [User/Developer] Fetches one flat with society context.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /**
         * Deactivate flat
         * @description [Owner/Admin/Staff] Deactivates a flat. This is not a hard delete.
         */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat deleted successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MessageAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        options?: never;
        head?: never;
        /**
         * Update flat
         * @description [Owner/Admin/Staff] Updates flat details, status, active state, or metadata.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            /** @description Update flat payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.UpdateFlatRequest"];
                };
            };
            responses: {
                /** @description Flat updated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/block": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Block flat
         * @description [Owner/Admin/Staff] Blocks an active flat so it cannot receive claims or residents.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat blocked successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invalid flat transition */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/member-invites": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List flat member invite history
         * @description [Resident] Lists member invite history for a flat. Requires flat owner or primary resident access.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Member invite status */
                    status?: "pending" | "accepted" | "expired" | "cancelled";
                    /** @description Search name, contact, role, status, or short code */
                    search?: string;
                    /** @description Maximum records to return */
                    limit?: number;
                    /** @description Records to skip */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member invites fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatMemberInvitesAPIResponse"];
                    };
                };
                /** @description Invalid query or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat member management access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        /**
         * Create flat member invite
         * @description [Resident] Creates a member invite for a flat and returns the invite plus shareable token details.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            /** @description Member invite request */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.CreateFlatMemberInviteRequest"];
                };
            };
            responses: {
                /** @description Member invite created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatMemberInviteTokenAPIResponse"];
                    };
                };
                /** @description Invalid request or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat member management access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/member-invites/{inviteId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat member invite detail
         * @description [Resident] Fetches one member invite for a flat with short-link details. Requires flat owner or primary resident access.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Member invite ID */
                    inviteId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member invite fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatMemberInviteAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat member management access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/member-invites/{inviteId}/cancel": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Cancel flat member invite
         * @description [Resident] Cancels a pending member invite for a flat.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Member invite ID */
                    inviteId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member invite cancelled successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MessageAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat member management access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member invite cannot be cancelled */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/members": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List flat residents (resident)
         * @description [Resident] Lists active members for a flat. Requires permission to view flat visitor data.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Resident role: owner, tenant, family */
                    role?: string;
                    /** @description Primary resident flag */
                    is_primary?: boolean;
                    /** @description Search user, contact, flat, block, role, or status */
                    search?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Residents fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatResidentsAPIResponse"];
                    };
                };
                /** @description Invalid path or query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/residents": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List flat residents for a society flat
         * @description [Owner/Admin/Staff] Lists residents for the flat identified in the path.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Resident role: owner, tenant, family */
                    role?: string;
                    /** @description Resident status: active, inactive, moved_out */
                    status?: string;
                    /** @description Primary resident flag */
                    is_primary?: boolean;
                    /** @description Search user, contact, flat, block, role, or status */
                    search?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Residents fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatResidentsAPIResponse"];
                    };
                };
                /** @description Invalid path or query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/residents/users/{userId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Add resident to flat
         * @description [Owner/Admin/Staff] Manually adds a user as an active resident of a flat.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description User ID */
                    userId: number;
                };
                cookie?: never;
            };
            /** @description Add flat resident payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.AddFlatResidentRequest"];
                };
            };
            responses: {
                /** @description Resident added successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatResidentAPIResponse"];
                    };
                };
                /** @description Invalid request or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat or user not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Duplicate resident, blocked flat, or primary resident conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/residents/{residentId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat resident
         * @description [User/Developer] Fetches one flat resident with joined user, flat, and society data.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Resident ID */
                    residentId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Resident fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatResidentAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        /**
         * Remove resident from flat
         * @description [Owner/Admin/Staff] Marks an active resident inactive and marks the flat vacant if this was the last active resident.
         */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Resident ID */
                    residentId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Resident removed successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MessageAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/residents/{residentId}/move-out": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Move out resident
         * @description [Owner/Admin/Staff] Marks a resident moved out and marks the flat vacant if this was the last active resident.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Resident ID */
                    residentId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Resident moved out successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatResidentAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/residents/{residentId}/primary": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Change primary resident
         * @description [Owner/Admin/Staff] Transactionally changes the active primary resident for a flat.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Resident ID */
                    residentId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Primary resident changed successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatResidentAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Primary resident conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/residents/{residentId}/role": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Update flat resident role
         * @description [Owner/Admin/Staff] Updates an active flat resident role.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Resident ID */
                    residentId: number;
                };
                cookie?: never;
            };
            /** @description Update flat resident role payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.UpdateFlatResidentRoleRequest"];
                };
            };
            responses: {
                /** @description Resident role updated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatResidentAPIResponse"];
                    };
                };
                /** @description Invalid request or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/unblock": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Unblock flat
         * @description [Owner/Admin/Staff] Unblocks a blocked flat and returns it to vacant state.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat unblocked successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invalid flat transition */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-context": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat visitor context for resident
         * @description [Resident] Returns occupancy, visitor settings summary, and recent visitors for a resident flat.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat visitor context fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatVisitorContextAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat or visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-entries": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List flat visitor entries
         * @description [Resident] Lists visitor entries for a resident flat with optional status, purpose, and date filters. created_from/created_to are legacy creation-time filters; event/event_from/event_to filter by visitor lifecycle timestamp.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Visitor status */
                    status?: "waiting_approval" | "approved" | "rejected" | "checked_in" | "checked_out" | "cancelled" | "expired" | "auto_closed";
                    /** @description Visitor purpose */
                    purpose?: "guest" | "delivery" | "cab" | "service" | "maintenance" | "staff" | "other";
                    /** @description Created from (RFC3339) */
                    created_from?: string;
                    /** @description Created to (RFC3339) */
                    created_to?: string;
                    /** @description Lifecycle timestamp to filter by */
                    event?: "created" | "checked_in" | "checked_out" | "expected" | "activity";
                    /** @description Lifecycle filter start (RFC3339) */
                    event_from?: string;
                    /** @description Lifecycle filter end (RFC3339) */
                    event_to?: string;
                    /** @description Maximum records to return */
                    limit?: number;
                    /** @description Records to skip */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor entries fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntriesAPIResponse"];
                    };
                };
                /** @description Invalid query or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-entries/pending": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List pending visitor approvals
         * @description [Resident] Lists visitor entries waiting for approval for a resident flat.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Pending visitor approvals fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntriesAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-entries/{entryId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat visitor entry
         * @description [Resident] Fetches a single visitor entry for a resident flat.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor entry fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry or flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-entries/{entryId}/photo": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * GetResidentVisitorPhoto
         * @description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Server-selected image view */
                    variant?: "list" | "avatar" | "detail" | "original";
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ImageViewAPIResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Request Entity Too Large */
                413: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Bad Gateway */
                502: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Service Unavailable */
                503: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-invites": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List visitor invite history
         * @description [Resident] Lists visitor invites for a resident flat with status and pagination filters.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Visitor invite status */
                    status?: "active" | "used" | "expired" | "cancelled";
                    /** @description Search purpose, status, visitor, contact, or short code */
                    search?: string;
                    /** @description Maximum records to return */
                    limit?: number;
                    /** @description Records to skip */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor invites fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorInviteHistoryAPIResponse"];
                    };
                };
                /** @description Invalid query or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        /**
         * Create visitor invite
         * @description [Resident] Creates a pre-approved visitor invite for a flat and returns the invite plus QR token details.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.CreateVisitorInviteRequest"];
            responses: {
                /** @description Visitor invite created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorInviteTokenAPIResponse"];
                    };
                };
                /** @description Invalid request, validation error, or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society, flat, or visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite cannot be created in current state */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-invites/staff": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Create visitor invite (staff)
         * @description [Staff] Creates a visitor invite for a flat and returns the invite plus shareable token details.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.CreateVisitorInviteRequest"];
            responses: {
                /** @description Visitor invite created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorInviteTokenAPIResponse"];
                    };
                };
                /** @description Invalid request, validation error, or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society, flat, or visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite cannot be created in current state */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-invites/{inviteId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get visitor invite detail
         * @description [Resident] Fetches one visitor invite for a resident flat with linked visitor, entry, and short-link details.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Visitor invite ID */
                    inviteId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor invite fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorInviteHistoryItemAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-settings": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get flat visitor settings
         * @description [Owner/Admin/Flat Resident] Fetches all visitor purpose settings for a flat.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat visitor settings fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatVisitorSettingsAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner/admin or flat resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat or visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-settings/reset": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reset flat visitor settings
         * @description [Owner/Admin/Flat Resident] Resets a flat's visitor purpose settings back to the production defaults.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Flat visitor settings reset successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MessageAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner/admin or flat resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat or visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/flats/{flatId}/visitor-settings/{purpose}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Update flat visitor purpose setting
         * @description [Owner/Admin/Flat Resident] Updates approval, duration override, or enabled state for one visitor purpose on a flat.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Flat ID */
                    flatId: number;
                    /** @description Visitor purpose */
                    purpose: "guest" | "delivery" | "cab" | "service" | "maintenance" | "staff" | "other";
                };
                cookie?: never;
            };
            /** @description Flat visitor purpose setting update payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.UpdateFlatVisitorSettingRequest"];
                };
            };
            responses: {
                /** @description Flat visitor setting updated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.FlatVisitorSettingAPIResponse"];
                    };
                };
                /** @description Invalid request, visitor purpose, or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner/admin or flat resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Flat or visitor setting not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/societies/{societyId}/guard-desk/bootstrap": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get guard desk bootstrap
         * @description [Owner/Admin/Staff] Returns aggregated guard desk dashboard data for a society.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Guard desk bootstrap fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.GuardDeskBootstrapAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/guards": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Create guard
         * @description [Owner/Admin] Creates a verified staff user for guard access without OTP or invite token.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Create guard payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.CreateGuardRequest"];
                };
            };
            responses: {
                /** @description Guard created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.GuardAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Duplicate email, phone, or member */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/members": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List society members
         * @description [User] Lists society members with filters for owner/admin member management screens.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Search full name, email, phone, role, or status */
                    search?: string;
                    /** @description Member role */
                    role?: string;
                    /** @description Member status */
                    status?: string;
                    /** @description User ID */
                    user_id?: number;
                    /** @description Invited by user ID */
                    invited_by?: number;
                    /** @description Removed by user ID */
                    removed_by?: number;
                    /** @description Joined from RFC3339 timestamp */
                    joined_from?: string;
                    /** @description Joined to RFC3339 timestamp */
                    joined_to?: string;
                    /** @description Limit */
                    limit?: number;
                    /** @description Offset */
                    offset?: number;
                    /** @description Sort by: joined_at, role, status */
                    sort_by?: string;
                    /** @description Sort order: asc, desc */
                    sort_order?: string;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Members fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.PaginatedMembersAPIResponse"];
                    };
                };
                /** @description Invalid query parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        /**
         * Add society member
         * @description [User] Adds a member to a society. Requires active owner/admin membership.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Add member payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.AddSocietyMemberRequest"];
                };
            };
            responses: {
                /** @description Member added successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyMemberAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Duplicate member or owner protection */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/members/summary": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get society member summary
         * @description [User] Returns member counts for a society.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member summary fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyMemberSummaryAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/members/{memberId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get society member
         * @description [User] Fetches one society member by member ID inside a society.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Member ID */
                    memberId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyMemberAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/members/{memberId}/visitor-approval-stats": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get member visitor approval stats
         * @description [Owner/Admin] Returns how many visitor entries a member has approved or rejected.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Member ID */
                    memberId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member visitor approval stats fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MemberVisitorApprovalStatsAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/members/{userId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /**
         * Remove member
         * @description [User] Removes a member from the society. Requires active owner/admin membership and protects the last owner.
         */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description User ID */
                    userId: number;
                };
                cookie?: never;
            };
            /** @description Removal reason */
            requestBody?: {
                content: {
                    "application/json": components["schemas"]["models.SocietyReasonRequest"];
                };
            };
            responses: {
                /** @description Member removed successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MessageAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner protection conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/members/{userId}/reactivate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reactivate member
         * @description [User] Reactivates a suspended or removed member. Requires active owner/admin membership.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description User ID */
                    userId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Member reactivated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyMemberAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/members/{userId}/role": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Change member role
         * @description [User] Changes a member role. Requires active owner/admin membership and protects the last owner.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description User ID */
                    userId: number;
                };
                cookie?: never;
            };
            /** @description Change role payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.ChangeSocietyMemberRoleRequest"];
                };
            };
            responses: {
                /** @description Member role changed successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyMemberAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner protection conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/societies/{societyId}/members/{userId}/suspend": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Suspend member
         * @description [User] Suspends a member. Requires active owner/admin membership and protects the last owner.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description User ID */
                    userId: number;
                };
                cookie?: never;
            };
            /** @description Suspension reason */
            requestBody?: {
                content: {
                    "application/json": components["schemas"]["models.SocietyReasonRequest"];
                };
            };
            responses: {
                /** @description Member suspended successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyMemberAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Member not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner protection conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/onboarding/bootstrap": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get society onboarding bootstrap
         * @description [Owner/Admin] Returns lifecycle, subscription, flat, and staff readiness before dashboard access.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Society onboarding bootstrap fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyOnboardingBootstrapAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/reactivate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reactivate society
         * @description [Developer] Reactivates a suspended society. Transition allowed: suspended -> active.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Society reactivated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invalid society lifecycle transition */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/reject": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reject society
         * @description [Developer] Rejects a pending society request. Transition allowed: pending -> rejected.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Rejection reason */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.SocietyReasonRequest"];
                };
            };
            responses: {
                /** @description Society rejected successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invalid society lifecycle transition */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/restore": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Restore soft-deleted society
         * @description [Developer] Restores a soft-deleted society to pending status.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Society restored successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society is not soft-deleted */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/subscriptions/plans/{planId}/pending": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Create pending subscription */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Plan ID */
                    planId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Subscription created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionAPIResponse"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/subscriptions/plans/{planId}/trial": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Create trial subscription */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Plan ID */
                    planId: number;
                };
                cookie?: never;
            };
            /** @description Trial subscription payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.CreateTrialSubscriptionRequest"];
                };
            };
            responses: {
                /** @description Trial subscription created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionAPIResponse"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/suspend": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Suspend society
         * @description [Developer] Suspends an active society. Transition allowed: active -> suspended.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Suspension reason */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.SocietyReasonRequest"];
                };
            };
            responses: {
                /** @description Society suspended successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invalid society lifecycle transition */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/transfer-ownership": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Transfer society ownership
         * @description [User] Transfers ownership to another user and leaves exactly one active owner.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Transfer ownership payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.TransferOwnershipRequest"];
                };
            };
            responses: {
                /** @description Ownership transferred successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyMemberAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner protection conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List visitor entries
         * @description [Owner/Admin/Staff] Lists visitor entries for a society with optional flat, status, source, purpose, limit, and offset filters.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Flat ID */
                    flat_id?: number;
                    /** @description Visitor status */
                    status?: "waiting_approval" | "approved" | "rejected" | "checked_in" | "checked_out" | "cancelled" | "expired" | "auto_closed";
                    /** @description Visitor entry source */
                    source?: "resident_link" | "public_qr" | "guard_entry" | "quick_link";
                    /** @description Visitor purpose */
                    purpose?: "guest" | "delivery" | "cab" | "service" | "maintenance" | "staff" | "other";
                    /** @description Maximum records to return */
                    limit?: number;
                    /** @description Records to skip */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor entries fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntriesAPIResponse"];
                    };
                };
                /** @description Invalid query or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/check-in": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Check in visitor
         * @description [Owner/Admin/Staff] Checks in a visitor using a valid QR token.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.QRTokenRequest"];
            responses: {
                /** @description Visitor checked in successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryAPIResponse"];
                    };
                };
                /** @description Invalid request, token, or visitor entry state */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description QR token is expired, unavailable, or already used */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/expected-guests": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List expected guest entries
         * @description [Owner/Admin/Staff] Lists pre-invited resident-link visitor entries expected today and ready for arrival.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Search by name, phone, flat, vehicle, or purpose */
                    search?: string;
                    /** @description Expected window start (RFC3339, defaults to IST today start) */
                    event_from?: string;
                    /** @description Expected window end (RFC3339, defaults to IST tomorrow start) */
                    event_to?: string;
                    /** @description Maximum records to return */
                    limit?: number;
                    /** @description Records to skip */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Expected guest entries fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntriesAPIResponse"];
                    };
                };
                /** @description Invalid query or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/guard": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Create guard visitor entry
         * @description [Owner/Admin/Staff] Creates a visitor entry from the guard desk for a society. Approval may be required depending on visitor settings.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["models.VisitorFormRequest"];
            responses: {
                /** @description Visitor entry created successfully */
                201: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryMutationAPIResponse"];
                    };
                };
                /** @description Invalid request, validation error, or society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Society, flat, or visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry cannot be created in current state */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/pending": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List society pending visitor approvals
         * @description [Owner/Admin/Staff] Lists visitor entries waiting for approval across the society.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Flat ID */
                    flat_id?: number;
                    /** @description Block */
                    block?: string;
                    /** @description Maximum records to return */
                    limit?: number;
                    /** @description Records to skip */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Pending visitor approvals fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorPendingEntriesAPIResponse"];
                    };
                };
                /** @description Invalid query or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/stats": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get visitor entry stats
         * @description [Owner/Admin/Staff] Returns visitor dashboard statistics for a society.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor entry stats fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryStatsAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/stats/daily": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get visitor entry daily stats
         * @description [Owner/Admin/Staff] Returns daily visitor counts for charting.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Number of days (default 7, max 90) */
                    days?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor entry daily stats fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryDailyStatsAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/waiting-at-gate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List visitors waiting at gate
         * @description [Owner/Admin/Staff] Lists approved non-invite visitor entries awaiting manual check-in. Excludes resident-link QR guests (see Expected Guests).
         */
        get: {
            parameters: {
                query?: {
                    /** @description Search by name, phone, flat, vehicle, delivery partner, or purpose */
                    search?: string;
                    /** @description Maximum records to return */
                    limit?: number;
                    /** @description Records to skip */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Waiting at gate entries fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntriesAPIResponse"];
                    };
                };
                /** @description Invalid query or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get visitor entry
         * @description [Owner/Admin/Staff] Fetches one visitor entry by ID.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor entry fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Update visitor entry details
         * @description [Owner/Admin/Staff] Updates visitor and entry details before check-in for approved or pending entries.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            /** @description Visitor entry update payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.UpdateGuardVisitorEntryRequest"];
                };
            };
            responses: {
                /** @description Visitor entry updated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryMutationAPIResponse"];
                    };
                };
                /** @description Invalid request or validation error */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry is not in an editable state */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}/approve": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Approve visitor entry
         * @description [Resident] Approves a waiting visitor entry and returns QR token details for the visitor.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor entry approved successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryMutationAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident approval access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry is not waiting for approval */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}/approve-and-check-in": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Guard approve and check in visitor
         * @description [Owner/Admin/Staff] Atomically approves and checks in a pending visitor entry.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: components["requestBodies"]["models.GuardApproveEntryRequest"];
            responses: {
                /** @description Visitor approved and checked in successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryAPIResponse"];
                    };
                };
                /** @description Invalid path parameter or visitor entry state */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}/check-in": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Check in visitor by entry ID
         * @description [Owner/Admin/Staff] Checks in an approved visitor entry without scanning QR.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor checked in successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryAPIResponse"];
                    };
                };
                /** @description Invalid path parameter or visitor entry state */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}/check-out": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Check out visitor
         * @description [Owner/Admin/Staff] Checks out a checked-in visitor entry.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor checked out successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryAPIResponse"];
                    };
                };
                /** @description Invalid path parameter or visitor entry state */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor is not checked in or already checked out */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}/events": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List visitor entry events
         * @description [Owner/Admin/Staff] Lists audit events recorded for a visitor entry.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor entry events fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryEventsAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}/guard-approve": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Guard approve visitor entry
         * @description [Owner/Admin/Staff] Approves a pending visitor entry without check-in.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: components["requestBodies"]["models.GuardApproveEntryRequest"];
            responses: {
                /** @description Visitor entry approved successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryMutationAPIResponse"];
                    };
                };
                /** @description Invalid path parameter or visitor entry state */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}/notify": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Notify resident about pending visitor
         * @description [Owner/Admin/Staff] Re-sends a pending approval notification to flat residents.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Resident notified successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.APIResponse"];
                    };
                };
                /** @description Invalid path parameter or visitor entry state */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Owner, admin, or staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}/photo": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * GetVisitorPhoto
         * @description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Server-selected image view */
                    variant?: "list" | "avatar" | "detail" | "original";
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ImageViewAPIResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Request Entity Too Large */
                413: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Bad Gateway */
                502: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Service Unavailable */
                503: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        /**
         * PutVisitorPhoto
         * @description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
         */
        put: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody: components["requestBodies"]["putV1AuthProfileAvatar"];
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ImageViewAPIResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Request Entity Too Large */
                413: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Bad Gateway */
                502: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Service Unavailable */
                503: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        post?: never;
        /**
         * DeleteVisitorPhoto
         * @description Managed private image. Viewing URLs expire after 15 minutes. Visitor images belong to the visitor record.
         */
        delete: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Image reference removed */
                204: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content?: never;
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Forbidden */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Conflict */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Request Entity Too Large */
                413: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Bad Gateway */
                502: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Service Unavailable */
                503: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-entries/{entryId}/reject": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reject visitor entry
         * @description [Resident] Rejects a waiting visitor entry with a reason.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor entry ID */
                    entryId: number;
                };
                cookie?: never;
            };
            /** @description Visitor rejection reason */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.RejectVisitorEntryRequest"];
                };
            };
            responses: {
                /** @description Visitor entry rejected successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorEntryMutationAPIResponse"];
                    };
                };
                /** @description Invalid request, validation error, or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Resident rejection access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor entry is not waiting for approval */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-invites/{inviteId}/cancel": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Cancel visitor invite
         * @description [Resident] Cancels a visitor invite created for a society before it is used.
         */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor invite ID */
                    inviteId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor invite cancelled successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.MessageAPIResponse"];
                    };
                };
                /** @description Invalid path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Invite owner or resident access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite cannot be cancelled */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-invites/{token}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get visitor invite by token (staff)
         * @description [Owner/Admin/Staff] Fetches an active invite only when the invite belongs to the authenticated staff member's society.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                    /** @description Visitor invite token */
                    token: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor invite fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.VisitorInviteAPIResponse"];
                    };
                };
                /** @description Invalid society ID or token */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Staff access required */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor invite is expired, used, or cancelled */
                409: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-settings": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get society visitor settings
         * @description [Owner/Admin] Fetches society-level visitor configuration for approval mode, QR entry, guard entry, pre-approval, durations, and active state.
         */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Visitor settings fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyVisitorSettingsAPIResponse"];
                    };
                };
                /** @description Invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Update society visitor settings
         * @description [Owner/Admin] Updates society-level visitor configuration. Boolean false values are accepted and persisted.
         */
        patch: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            /** @description Society visitor settings update payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.UpdateSocietyVisitorSettingsRequest"];
                };
            };
            responses: {
                /** @description Visitor settings updated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyVisitorSettingsAPIResponse"];
                    };
                };
                /** @description Invalid request, validation error, or invalid society ID */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Visitor settings not found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        trace?: never;
    };
    "/v1/societies/{societyId}/visitor-settings/flats": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * List society flat visitor settings
         * @description [Owner/Admin] Lists flat-level visitor purpose settings across the society with pagination.
         */
        get: {
            parameters: {
                query?: {
                    /** @description Flat ID */
                    flat_id?: number;
                    /** @description Block */
                    block?: string;
                    /** @description Visitor purpose */
                    purpose?: "guest" | "delivery" | "cab" | "service" | "maintenance" | "staff" | "other";
                    /** @description Maximum records to return */
                    limit?: number;
                    /** @description Records to skip */
                    offset?: number;
                };
                header?: never;
                path: {
                    /** @description Society ID */
                    societyId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Society flat visitor settings fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SocietyFlatVisitorSettingsAPIResponse"];
                    };
                };
                /** @description Invalid query or path parameter */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Missing, invalid, or expired access token */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Insufficient society role */
                403: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Internal server error */
                500: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/subscriptions": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** List subscriptions */
        get: {
            parameters: {
                query?: {
                    /** @description Society ID */
                    society_id?: number;
                    /** @description Plan ID */
                    plan_id?: number;
                    /** @description Subscription status */
                    status?: string;
                    /** @description Search text */
                    search?: string;
                    /** @description Search mode */
                    search_mode?: string;
                    /** @description Only active subscriptions */
                    is_active_only?: boolean;
                    /** @description Only expired subscriptions */
                    expired_only?: boolean;
                    /** @description Subscriptions ending before RFC3339 timestamp */
                    expiring_before?: string;
                    /** @description Subscriptions ending within N days */
                    expiring_days?: number;
                    /** @description Page size */
                    limit?: number;
                    /** @description Page offset */
                    offset?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Subscriptions fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionsAPIResponse"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/subscriptions/lookup": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get subscription */
        get: {
            parameters: {
                query?: {
                    /** @description Subscription ID */
                    id?: number;
                    /** @description Society ID */
                    society_id?: number;
                };
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Subscription fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionAPIResponse"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/subscriptions/stats": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get subscription stats */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path?: never;
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Subscription stats fetched successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionStatsAPIResponse"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/subscriptions/{subscriptionId}/activate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Activate subscription */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Subscription ID */
                    subscriptionId: number;
                };
                cookie?: never;
            };
            /** @description Activate subscription payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.ActivateSubscriptionRequest"];
                };
            };
            responses: {
                /** @description Subscription activated successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionAPIResponse"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/subscriptions/{subscriptionId}/cancel": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Cancel subscription */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Subscription ID */
                    subscriptionId: number;
                };
                cookie?: never;
            };
            /** @description Cancel subscription payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.CancelSubscriptionRequest"];
                };
            };
            responses: {
                /** @description Subscription cancelled successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionAPIResponse"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/subscriptions/{subscriptionId}/expire": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Expire subscription */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Subscription ID */
                    subscriptionId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Subscription expired successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionAPIResponse"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/subscriptions/{subscriptionId}/plans/{planId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Change subscription plan */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Subscription ID */
                    subscriptionId: number;
                    /** @description New plan ID */
                    planId: number;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description Subscription plan changed successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionAPIResponse"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/subscriptions/{subscriptionId}/renew": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Renew subscription */
        post: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Subscription ID */
                    subscriptionId: number;
                };
                cookie?: never;
            };
            /** @description Renew subscription payload */
            requestBody: {
                content: {
                    "application/json": components["schemas"]["models.RenewSubscriptionRequest"];
                };
            };
            responses: {
                /** @description Subscription renewed successfully */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.SubscriptionAPIResponse"];
                    };
                };
            };
        };
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/me/notifications/{notificationId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get one owned notification for authenticated tap routing */
        get: {
            parameters: {
                query?: never;
                header?: never;
                path: {
                    /** @description Notification ID */
                    notificationId: string;
                };
                cookie?: never;
            };
            requestBody?: never;
            responses: {
                /** @description OK */
                200: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.NotificationDetailAPIResponse"];
                    };
                };
                /** @description Bad Request */
                400: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Unauthorized */
                401: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
                /** @description Not Found */
                404: {
                    headers: {
                        [name: string]: unknown;
                    };
                    content: {
                        "application/json": components["schemas"]["models.ErrorResponseDoc"];
                    };
                };
            };
        };
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        "models.APIResponse": {
            data?: unknown;
            error?: components["schemas"]["models.ErrorData"];
            message?: string;
            success?: boolean;
        };
        "models.AcceptFlatMemberInviteAPIResponse": {
            data?: components["schemas"]["models.AcceptFlatMemberInviteData"];
            /** @example Member invite accepted successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.AcceptFlatMemberInviteData": {
            acceptance?: components["schemas"]["models.AcceptFlatMemberInviteResponse"];
        };
        "models.AcceptFlatMemberInviteResponse": {
            invite?: components["schemas"]["models.FlatMemberInviteResponse"];
            resident?: components["schemas"]["models.FlatResidentResponse"];
        };
        "models.ActivateSubscriptionRequest": {
            ends_at: string;
            metadata?: {
                [key: string]: unknown;
            };
            starts_at: string;
        };
        "models.AddFlatResidentRequest": {
            is_primary?: boolean;
            metadata?: {
                [key: string]: unknown;
            };
            role: components["schemas"]["models.FlatResidentRole"];
        };
        "models.AddSocietyMemberRequest": {
            metadata?: {
                [key: string]: unknown;
            };
            role: components["schemas"]["models.SocietyMemberRole"];
            society_id: number;
            user_id: number;
        };
        "models.ApproveFlatClaimResponse": {
            claim?: components["schemas"]["models.FlatClaimResponse"];
            flat?: components["schemas"]["models.FlatResponse"];
            resident?: components["schemas"]["models.FlatResidentResponse"];
        };
        /** @enum {string} */
        "models.AuthProvider": "email" | "google" | "apple" | "phone";
        "models.AuthRefreshData": {
            access_token?: string;
            access_token_expires_at?: string;
            /** @example Access token refreshed successfully */
            message?: string;
        };
        "models.AuthSessionData": {
            access_token?: string;
            access_token_expires_at?: string;
            refresh_token?: string;
            refresh_token_expires_at?: string;
            user?: components["schemas"]["models.UserResponse"];
        };
        /** @enum {string} */
        "models.BillingCycle": "monthly" | "yearly";
        "models.BootstrapAPIResponse": {
            data?: components["schemas"]["models.BootstrapData"];
            /** @example Bootstrap fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.BootstrapData": {
            defaultDashboard?: components["schemas"]["models.DefaultDashboardResponse"];
            memberships?: components["schemas"]["models.SocietyMemberResponse"][];
            residences?: components["schemas"]["models.FlatResidentResponse"][];
            user?: components["schemas"]["models.UserResponse"];
        };
        "models.BulkCreateFlatsRequest": {
            flats: components["schemas"]["models.CreateFlatRequest"][];
        };
        "models.BulkCreateFlatsResponse": {
            items?: components["schemas"]["models.FlatResponse"][];
            total?: number;
        };
        "models.BulkFlatsAPIResponse": {
            data?: components["schemas"]["models.BulkFlatsData"];
            /** @example Flats created successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.BulkFlatsData": {
            flats?: components["schemas"]["models.BulkCreateFlatsResponse"];
        };
        "models.CancelSubscriptionRequest": {
            metadata?: {
                [key: string]: unknown;
            };
            reason: string;
        };
        "models.ChangePasswordAPIResponse": {
            data?: components["schemas"]["models.MessageData"];
            /** @example Password changed successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.ChangePasswordRequest": {
            confirm_password: string;
            current_password: string;
            new_password: string;
        };
        "models.ChangeSocietyMemberRoleRequest": {
            role: components["schemas"]["models.SocietyMemberRole"];
            society_id: number;
            user_id: number;
        };
        "models.CreateFlatMemberInviteRequest": {
            email?: string;
            full_name: string;
            phone?: string;
            role: components["schemas"]["models.FlatMemberInviteRole"];
        };
        "models.CreateFlatRequest": {
            block?: string;
            flat_number: string;
            floor?: string;
            metadata?: {
                [key: string]: unknown;
            };
        };
        "models.CreateGuardRequest": {
            email: string;
            first_name: string;
            last_name?: string;
            password: string;
            phone_number: string;
        };
        "models.CreateGuardResponse": {
            member?: components["schemas"]["models.SocietyMemberResponse"];
            user?: components["schemas"]["models.UserResponse"];
        };
        "models.CreatePlanRequest": {
            billing_cycle: components["schemas"]["models.BillingCycle"];
            code: string;
            currency: string;
            description?: string;
            features?: {
                [key: string]: unknown;
            };
            max_admins?: number;
            max_flats: number;
            max_residents: number;
            max_staff?: number;
            name: string;
            price_amount_paise?: number;
        };
        "models.CreateSocietyRequest": {
            address_line1?: string;
            address_line2?: string;
            city?: string;
            country?: string;
            email?: string;
            landmark?: string;
            metadata?: {
                [key: string]: unknown;
            };
            name: string;
            phone_number?: string;
            pincode?: string;
            society_code?: string;
            state?: string;
            total_blocks?: number;
            total_flats?: number;
        };
        "models.CreateTrialSubscriptionRequest": {
            ends_at?: string;
            metadata?: {
                [key: string]: unknown;
            };
            starts_at: string;
            trial_ends_at: string;
        };
        "models.CreateVisitorInviteRequest": {
            expires_at?: string;
            purpose?: components["schemas"]["models.VisitorPurpose"];
        };
        /** @enum {string} */
        "models.DashboardKind": "developer" | "society_admin" | "select_society" | "onboarding";
        "models.DefaultDashboardResponse": {
            kind?: components["schemas"]["models.DashboardKind"];
            path?: string;
            society_id?: number;
        };
        "models.DevOTPMessageData": {
            /** @example 123456 */
            dev_otp?: string;
            /** @example Please verify your email using the OTP sent to your email address */
            message?: string;
        };
        "models.DeveloperDashboardBootstrapAPIResponse": {
            data?: components["schemas"]["models.DeveloperDashboardBootstrapData"];
            /** @example Developer dashboard bootstrap fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.DeveloperDashboardBootstrapData": {
            dashboard?: components["schemas"]["models.DeveloperDashboardBootstrapResponse"];
        };
        "models.DeveloperDashboardBootstrapResponse": {
            plan_stats?: components["schemas"]["models.DeveloperDashboardPlanStatsResponse"];
            recent_pending_societies?: components["schemas"]["models.SocietyResponse"][];
            recent_subscriptions?: components["schemas"]["models.SocietySubscriptionResponse"][];
            residence_stats?: components["schemas"]["models.DeveloperDashboardResidenceStatsResponse"];
            society_stats?: components["schemas"]["models.DeveloperDashboardSocietyStatsResponse"];
            subscription_stats?: components["schemas"]["models.SubscriptionStatsResponse"];
        };
        "models.DeveloperDashboardPlanStatsResponse": {
            active?: number;
            inactive?: number;
            total?: number;
        };
        "models.DeveloperDashboardResidenceStatsResponse": {
            active_residents?: number;
            total_residents?: number;
        };
        "models.DeveloperDashboardSocietyStatsResponse": {
            active?: number;
            pending?: number;
            rejected?: number;
            suspended?: number;
            total?: number;
        };
        /** @enum {string} */
        "models.DevicePlatform": "ios" | "android" | "web";
        "models.DeviceToken": {
            created_at?: string;
            device_id?: string;
            id?: number;
            last_seen_at?: string;
            platform?: components["schemas"]["models.DevicePlatform"];
            token?: string;
            updated_at?: string;
            user_id?: number;
        };
        "models.DeviceTokenData": {
            device_token?: components["schemas"]["models.DeviceToken"];
        };
        "models.ErrorData": {
            code?: string;
            details?: {
                [key: string]: unknown;
            };
            message?: string;
        };
        "models.ErrorResponseDoc": {
            error?: components["schemas"]["models.ErrorData"];
            /** @example false */
            success?: boolean;
        };
        "models.FlatAPIResponse": {
            data?: components["schemas"]["models.FlatData"];
            /** @example Flat fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatApprovalAPIResponse": {
            data?: components["schemas"]["models.FlatApprovalData"];
            /** @example Flat claim approved successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatApprovalData": {
            approval?: components["schemas"]["models.ApproveFlatClaimResponse"];
        };
        "models.FlatClaimAPIResponse": {
            data?: components["schemas"]["models.FlatClaimData"];
            /** @example Flat claim fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatClaimData": {
            claim?: components["schemas"]["models.FlatClaimResponse"];
        };
        "models.FlatClaimResponse": {
            block?: string;
            cancelled_at?: string;
            created_at?: string;
            flat_id?: number;
            flat_number?: string;
            flat_status?: components["schemas"]["models.FlatStatus"];
            floor?: string;
            id?: number;
            note?: string;
            rejection_reason?: string;
            requested_primary?: boolean;
            requested_role?: components["schemas"]["models.FlatResidentRole"];
            reviewed_at?: string;
            reviewed_by?: number;
            reviewer_email?: string;
            reviewer_name?: string;
            reviewer_phone?: string;
            society_code?: string;
            society_id?: number;
            society_name?: string;
            status?: components["schemas"]["models.FlatClaimStatus"];
            updated_at?: string;
            user_email?: string;
            user_id?: number;
            user_name?: string;
            user_phone?: string;
        };
        "models.FlatClaimStatsAPIResponse": {
            data?: components["schemas"]["models.FlatClaimStatsData"];
            /** @example Flat claim stats fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatClaimStatsData": {
            stats?: components["schemas"]["models.FlatClaimStatsResponse"];
        };
        "models.FlatClaimStatsResponse": {
            approved_claims?: number;
            cancelled_claims?: number;
            pending_claims?: number;
            rejected_claims?: number;
            total_claims?: number;
        };
        /** @enum {string} */
        "models.FlatClaimStatus": "pending" | "approved" | "rejected" | "cancelled";
        "models.FlatClaimsAPIResponse": {
            data?: components["schemas"]["models.FlatClaimsData"];
            /** @example Flat claims fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatClaimsData": {
            claims?: components["schemas"]["models.FlatClaimResponse"][];
        };
        "models.FlatData": {
            flat?: components["schemas"]["models.FlatResponse"];
        };
        "models.FlatMemberInviteAPIResponse": {
            data?: components["schemas"]["models.FlatMemberInviteData"];
            /** @example Member invite fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatMemberInviteData": {
            invite?: components["schemas"]["models.FlatMemberInviteResponse"];
        };
        "models.FlatMemberInviteResponse": {
            created_at?: string;
            email?: string;
            expires_at?: string;
            flat_id?: number;
            full_name?: string;
            id?: number;
            invited_by?: number;
            phone?: string;
            role?: components["schemas"]["models.FlatMemberInviteRole"];
            short_link?: components["schemas"]["models.ShortLinkResponse"];
            society_id?: number;
            status?: components["schemas"]["models.FlatMemberInviteStatus"];
            updated_at?: string;
        };
        /** @enum {string} */
        "models.FlatMemberInviteRole": "family" | "tenant";
        /** @enum {string} */
        "models.FlatMemberInviteStatus": "pending" | "accepted" | "expired" | "cancelled";
        "models.FlatMemberInviteTokenAPIResponse": {
            data?: components["schemas"]["models.FlatMemberInviteTokenData"];
            /** @example Member invite created successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatMemberInviteTokenData": {
            invite?: components["schemas"]["models.FlatMemberInviteResponse"];
            link?: components["schemas"]["models.ShortLinkResponse"];
            short_link?: components["schemas"]["models.ShortLinkResponse"];
            token?: components["schemas"]["models.FlatMemberInviteTokenResponse"];
        };
        "models.FlatMemberInviteTokenResponse": {
            expires_at?: string;
            token?: string;
        };
        "models.FlatMemberInvitesAPIResponse": {
            data?: components["schemas"]["models.FlatMemberInvitesData"];
            /** @example Member invites fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatMemberInvitesData": {
            invites?: components["schemas"]["models.FlatMemberInviteResponse"][];
            limit?: number;
            offset?: number;
            total?: number;
        };
        /** @enum {string} */
        "models.FlatNumberingMode": "continuous" | "floor_based";
        "models.FlatRecentVisitorSummary": {
            entry_id?: number;
            full_name?: string;
            purpose?: components["schemas"]["models.VisitorPurpose"];
            status?: components["schemas"]["models.VisitorStatus"];
            visited_on?: string;
        };
        "models.FlatResidentAPIResponse": {
            data?: components["schemas"]["models.FlatResidentData"];
            /** @example Resident fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatResidentData": {
            resident?: components["schemas"]["models.FlatResidentResponse"];
        };
        "models.FlatResidentResponse": {
            block?: string;
            created_at?: string;
            created_by?: number;
            flat_id?: number;
            flat_number?: string;
            floor?: string;
            id?: number;
            is_primary?: boolean;
            moved_in_at?: string;
            moved_out_at?: string;
            role?: components["schemas"]["models.FlatResidentRole"];
            society_code?: string;
            society_id?: number;
            society_name?: string;
            status?: components["schemas"]["models.FlatResidentStatus"];
            updated_at?: string;
            user_email?: string;
            user_id?: number;
            user_name?: string;
            user_phone?: string;
        };
        /** @enum {string} */
        "models.FlatResidentRole": "owner" | "tenant" | "family";
        /** @enum {string} */
        "models.FlatResidentStatus": "active" | "inactive" | "moved_out";
        "models.FlatResidentsAPIResponse": {
            data?: components["schemas"]["models.FlatResidentsData"];
            /** @example Residents fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatResidentsData": {
            residents?: components["schemas"]["models.FlatResidentResponse"][];
        };
        "models.FlatResponse": {
            block?: string;
            created_at?: string;
            created_by?: number;
            flat_number?: string;
            floor?: string;
            id?: number;
            is_active?: boolean;
            primary_resident_name?: string;
            society_code?: string;
            society_id?: number;
            society_name?: string;
            status?: components["schemas"]["models.FlatStatus"];
            updated_at?: string;
        };
        "models.FlatStatsAPIResponse": {
            data?: components["schemas"]["models.FlatStatsData"];
            /** @example Flat stats fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatStatsData": {
            stats?: components["schemas"]["models.FlatStatsResponse"];
        };
        "models.FlatStatsResponse": {
            active_flats?: number;
            blocked_flats?: number;
            inactive_flats?: number;
            occupied_flats?: number;
            society_id?: number;
            total_flats?: number;
            vacant_flats?: number;
        };
        /** @enum {string} */
        "models.FlatStatus": "vacant" | "occupied" | "blocked";
        "models.FlatVisitorContextAPIResponse": {
            data?: components["schemas"]["models.FlatVisitorContextData"];
            /** @example Flat visitor context fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatVisitorContextData": {
            context?: components["schemas"]["models.FlatVisitorContextResponse"];
        };
        "models.FlatVisitorContextResident": {
            full_name?: string;
            id?: number;
            user_id?: number;
        };
        "models.FlatVisitorContextResponse": {
            inherits_society_mode?: boolean;
            occupancy_status?: components["schemas"]["models.FlatStatus"];
            primary_resident?: components["schemas"]["models.FlatVisitorContextResident"];
            recent_visitors?: components["schemas"]["models.FlatRecentVisitorSummary"][];
            society_approval_mode?: components["schemas"]["models.VisitorApprovalMode"];
            total_residents?: number;
            visitor_settings?: components["schemas"]["models.FlatVisitorSettingsResponse"][];
        };
        "models.FlatVisitorSettingAPIResponse": {
            data?: components["schemas"]["models.FlatVisitorSettingData"];
            /** @example Flat visitor setting updated successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatVisitorSettingData": {
            visitor_setting?: components["schemas"]["models.FlatVisitorSettingsResponse"];
        };
        "models.FlatVisitorSettingsAPIResponse": {
            data?: components["schemas"]["models.FlatVisitorSettingsData"];
            /** @example Flat visitor settings fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatVisitorSettingsData": {
            visitor_settings?: components["schemas"]["models.FlatVisitorSettingsResponse"][];
        };
        "models.FlatVisitorSettingsResponse": {
            approval_required?: boolean;
            created_at?: string;
            default_visit_duration_minutes?: number;
            flat_id?: number;
            id?: number;
            is_enabled?: boolean;
            purpose?: components["schemas"]["models.VisitorPurpose"];
            society_id?: number;
            updated_at?: string;
            updated_by?: number;
        };
        "models.FlatsAPIResponse": {
            data?: components["schemas"]["models.FlatsData"];
            /** @example Flats fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.FlatsData": {
            flats?: components["schemas"]["models.FlatResponse"][];
        };
        "models.ForgotPasswordAPIResponse": {
            data?: components["schemas"]["models.DevOTPMessageData"];
            /** @example if this email exists, password reset instructions have been sent */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.ForgotPasswordRequest": {
            email: string;
        };
        "models.GenerateFlatsRequest": {
            block: string;
            flats_per_floor: number;
            number_of_floors?: number;
            numbering_mode: components["schemas"]["models.FlatNumberingMode"];
            sequence_start?: string;
            start_floor?: number;
            total_flats?: number;
            unit_start?: string;
        };
        "models.GetProfileAPIResponse": {
            data?: components["schemas"]["models.UserData"];
            /** @example Profile fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        /** @enum {string} */
        "models.GlobalRole": "user" | "developer" | "super_admin";
        "models.GuardAPIResponse": {
            data?: components["schemas"]["models.GuardData"];
            /** @example Guard created successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.GuardApproveEntryRequest": {
            on_behalf?: boolean;
            reason?: string;
        };
        "models.GuardData": {
            guard?: components["schemas"]["models.CreateGuardResponse"];
        };
        "models.GuardDeskBootstrapAPIResponse": {
            data?: components["schemas"]["models.GuardDeskBootstrapData"];
            /** @example Guard desk bootstrap fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.GuardDeskBootstrapData": {
            desk?: components["schemas"]["models.GuardDeskBootstrapResponse"];
        };
        "models.GuardDeskBootstrapResponse": {
            expected_guests_count?: number;
            pending_preview?: components["schemas"]["models.VisitorPendingEntry"][];
            society?: components["schemas"]["models.SocietyResponse"];
            stats?: components["schemas"]["models.VisitorEntryStatsResponse"];
            visitor_settings?: components["schemas"]["models.GuardDeskVisitorSettingsSummary"];
            waiting_at_gate_count?: number;
        };
        "models.GuardDeskVisitorSettingsSummary": {
            allow_guard_entry?: boolean;
            allow_guard_on_behalf_approval?: boolean;
        };
        "models.HealthCheckResponseDoc": {
            /** @example go-server */
            app?: string;
            /** @example development */
            environment?: string;
            /** @example healthy */
            status?: string;
            /** @example 2026-05-26T02:55:50Z */
            timestamp?: string;
            /** @example 1.0.0 */
            version?: string;
        };
        "models.HealthErrorResponseDoc": {
            /** @example database connection failed */
            error?: string;
            /** @example not_ready */
            status?: string;
        };
        /** @enum {string} */
        "models.ImageVariant": "original" | "list" | "avatar" | "detail";
        "models.ImageView": {
            expires_at?: string;
            url?: string;
            variant?: components["schemas"]["models.ImageVariant"];
        };
        "models.ImageViewAPIResponse": {
            data?: components["schemas"]["models.ImageView"];
            message?: string;
            success?: boolean;
        };
        "models.JoinFlatMemberInviteAPIResponse": {
            data?: components["schemas"]["models.JoinFlatMemberInviteData"];
            /** @example Member joined successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.JoinFlatMemberInviteData": {
            join?: components["schemas"]["models.JoinFlatMemberInviteResponse"];
        };
        "models.JoinFlatMemberInviteRequest": {
            email?: string;
            first_name?: string;
            identifier?: string;
            last_name?: string;
            password: string;
        };
        "models.JoinFlatMemberInviteResponse": {
            acceptance?: components["schemas"]["models.AcceptFlatMemberInviteResponse"];
            user?: components["schemas"]["models.UserResponse"];
        };
        "models.LivenessResponseDoc": {
            /** @example alive */
            status?: string;
            /** @example 2026-05-26T02:55:50Z */
            timestamp?: string;
        };
        "models.LoginAPIResponse": {
            data?: components["schemas"]["models.AuthSessionData"];
            /** @example Login successful */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.LoginRequest": {
            email?: string;
            password: string;
            phone_number?: string;
        };
        "models.LogoutAPIResponse": {
            data?: components["schemas"]["models.MessageData"];
            /** @example Logout successful */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.MemberVisitorApprovalStatsAPIResponse": {
            data?: components["schemas"]["models.MemberVisitorApprovalStatsData"];
            /** @example Member visitor approval stats fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.MemberVisitorApprovalStatsData": {
            stats?: components["schemas"]["models.MemberVisitorApprovalStatsResponse"];
        };
        "models.MemberVisitorApprovalStatsResponse": {
            approved_count?: number;
            rejected_count?: number;
        };
        "models.MessageAPIResponse": {
            data?: components["schemas"]["models.MessageData"];
            /** @example Request completed successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.MessageData": {
            /** @example Request completed successfully */
            message?: string;
        };
        "models.MyResidencesAPIResponse": {
            data?: components["schemas"]["models.MyResidencesData"];
            /** @example My residences fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.MyResidencesData": {
            residences?: components["schemas"]["models.FlatResidentResponse"][];
        };
        "models.MySocietiesAPIResponse": {
            data?: components["schemas"]["models.MySocietiesData"];
            /** @example My societies fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.MySocietiesData": {
            societies?: components["schemas"]["models.MySocietyResponse"][];
        };
        "models.MySocietyResponse": {
            member?: components["schemas"]["models.SocietyMemberResponse"];
            society?: components["schemas"]["models.SocietyResponse"];
        };
        "models.Notification": {
            body?: string;
            created_at?: string;
            data?: {
                [key: string]: unknown;
            };
            event_key?: string;
            flat_id?: number;
            id?: string;
            read_at?: string;
            society_id?: number;
            title?: string;
            type?: string;
            user_id?: number;
        };
        "models.NotificationListResult": {
            items?: components["schemas"]["models.Notification"][];
            next_cursor?: string;
        };
        "models.NotificationReadAPIResponse": {
            data?: components["schemas"]["models.NotificationReadResponse"];
            /** @example Notification marked read */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.NotificationReadResponse": {
            id?: string;
            read_at?: string;
        };
        "models.NotificationUnreadCountAPIResponse": {
            data?: components["schemas"]["models.NotificationUnreadCountResponse"];
            /** @example Unread count fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.NotificationUnreadCountResponse": {
            unread_count?: number;
        };
        "models.NotificationsAPIResponse": {
            data?: components["schemas"]["models.NotificationListResult"];
            /** @example Notifications fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.PaginatedFlatsAPIResponse": {
            data?: components["schemas"]["models.PaginatedFlatsData"];
            /** @example Flats fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.PaginatedFlatsData": {
            flats?: components["schemas"]["models.PaginatedFlatsResponse"];
        };
        "models.PaginatedFlatsResponse": {
            items?: components["schemas"]["models.FlatResponse"][];
            limit?: number;
            offset?: number;
            total?: number;
        };
        "models.PaginatedMembersAPIResponse": {
            data?: components["schemas"]["models.PaginatedMembersData"];
            /** @example Members fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.PaginatedMembersData": {
            members?: components["schemas"]["models.PaginatedMembersResponse"];
        };
        "models.PaginatedMembersResponse": {
            items?: components["schemas"]["models.SocietyMemberResponse"][];
            limit?: number;
            offset?: number;
            total?: number;
        };
        "models.PaginatedSocietiesAPIResponse": {
            data?: components["schemas"]["models.PaginatedSocietiesData"];
            /** @example Societies fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.PaginatedSocietiesData": {
            societies?: components["schemas"]["models.PaginatedSocietiesResponse"];
        };
        "models.PaginatedSocietiesResponse": {
            items?: components["schemas"]["models.SocietyResponse"][];
            limit?: number;
            offset?: number;
            total?: number;
        };
        "models.PlanAPIResponse": {
            data?: components["schemas"]["models.PlanData"];
            /** @example Plan fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.PlanData": {
            plan?: components["schemas"]["models.PlanResponse"];
        };
        "models.PlanResponse": {
            billing_cycle?: components["schemas"]["models.BillingCycle"];
            code?: string;
            created_at?: string;
            currency?: string;
            description?: string;
            features?: {
                [key: string]: unknown;
            };
            id?: number;
            is_active?: boolean;
            max_admins?: number;
            max_flats?: number;
            max_residents?: number;
            max_staff?: number;
            name?: string;
            price_amount_paise?: number;
            updated_at?: string;
        };
        "models.PlansAPIResponse": {
            data?: components["schemas"]["models.PlansData"];
            /** @example Plans fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.PlansData": {
            plans?: components["schemas"]["models.PlanResponse"][];
        };
        "models.PublicClaimFlatResponse": {
            block?: string;
            flat_number?: string;
            floor?: string;
            id?: number;
            status?: components["schemas"]["models.FlatStatus"];
        };
        "models.PublicClaimOptionsAPIResponse": {
            data?: components["schemas"]["models.PublicClaimOptionsData"];
            /** @example Claim options fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.PublicClaimOptionsData": {
            flats?: components["schemas"]["models.PublicClaimFlatResponse"][];
            society?: components["schemas"]["models.PublicClaimSocietyResponse"];
        };
        "models.PublicClaimSocietyResponse": {
            city?: string;
            country?: string;
            id?: number;
            name?: string;
            pincode?: string;
            society_code?: string;
            state?: string;
            total_flats?: number;
        };
        "models.PublicFlatMemberInviteAPIResponse": {
            data?: components["schemas"]["models.PublicFlatMemberInviteData"];
            /** @example Member invite fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.PublicFlatMemberInviteData": {
            invite?: components["schemas"]["models.PublicFlatMemberInviteView"];
        };
        "models.PublicFlatMemberInviteView": {
            block?: string;
            email?: string;
            expires_at?: string;
            flat_number?: string;
            floor?: string;
            full_name?: string;
            id?: number;
            phone?: string;
            role?: components["schemas"]["models.FlatMemberInviteRole"];
            society_name?: string;
            status?: components["schemas"]["models.FlatMemberInviteStatus"];
        };
        "models.PublicVisitorInviteView": {
            block?: string;
            expires_at?: string;
            flat_number?: string;
            floor?: string;
            id?: number;
            purpose?: components["schemas"]["models.VisitorPurpose"];
            society_name?: string;
            status?: components["schemas"]["models.VisitorInviteStatus"];
        };
        "models.QRTokenRequest": {
            token?: string;
        };
        "models.QRTokenResponse": {
            expires_at?: string;
            token?: string;
        };
        "models.ReadinessResponseDoc": {
            /** @example connected */
            database?: string;
            /** @example ready */
            status?: string;
            /** @example 2026-05-26T02:55:50Z */
            timestamp?: string;
        };
        "models.RefreshTokenAPIResponse": {
            data?: components["schemas"]["models.AuthRefreshData"];
            /** @example Access token refreshed successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.RefreshTokenRequest": {
            refresh_token?: string;
        };
        "models.RegisterAPIResponse": {
            data?: components["schemas"]["models.UserMessageData"];
            /** @example Account created successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.RegisterDeviceTokenAPIResponse": {
            data?: components["schemas"]["models.DeviceTokenData"];
            /** @example Device token registered successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.RegisterDeviceTokenRequest": {
            device_id?: string;
            platform: components["schemas"]["models.DevicePlatform"];
            token: string;
        };
        "models.RegisterRequest": {
            email: string;
            first_name: string;
            last_name?: string;
            password: string;
            phone_number: string;
        };
        "models.RejectFlatClaimRequest": {
            reason: string;
        };
        "models.RejectVisitorEntryRequest": {
            reason?: string;
        };
        "models.RenewSubscriptionRequest": {
            ends_at: string;
            metadata?: {
                [key: string]: unknown;
            };
            starts_at: string;
        };
        "models.ResendOTPAPIResponse": {
            data?: components["schemas"]["models.DevOTPMessageData"];
            /** @example verification OTP sent successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.ResendOTPRequest": {
            email: string;
        };
        "models.ResetPasswordAPIResponse": {
            data?: components["schemas"]["models.MessageData"];
            /** @example Password reset successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.ResetPasswordRequest": {
            confirm_password: string;
            email: string;
            new_password: string;
            otp: string;
        };
        "models.ResidentRegisterRequest": {
            email: string;
            first_name: string;
            last_name?: string;
            password: string;
            phone_number: string;
        };
        "models.ShortLinkResolveAPIResponse": {
            data?: components["schemas"]["models.ShortLinkResolveData"];
            /** @example Short link resolved successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.ShortLinkResolveData": {
            link?: components["schemas"]["models.ShortLinkResolveResponse"];
        };
        "models.ShortLinkResolveResponse": {
            entry?: components["schemas"]["models.VisitorEntry"];
            invite?: unknown;
            link?: components["schemas"]["models.ShortLinkResponse"];
            qr?: components["schemas"]["models.QRTokenResponse"];
            state?: components["schemas"]["models.ShortLinkState"];
            type?: components["schemas"]["models.ShortLinkResourceType"];
            visitor?: components["schemas"]["models.VisitorSummary"];
        };
        /** @enum {string} */
        "models.ShortLinkResourceType": "visitor_invite" | "member_invite";
        "models.ShortLinkResponse": {
            code?: string;
            expires_at?: string;
            url?: string;
        };
        /** @enum {string} */
        "models.ShortLinkState": "form" | "waiting_approval" | "approved" | "checked_in" | "completed" | "join" | "accepted" | "expired" | "cancelled" | "unavailable";
        "models.SocietyAPIResponse": {
            data?: components["schemas"]["models.SocietyData"];
            /** @example Society request created successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SocietyDashboardBootstrapAPIResponse": {
            data?: components["schemas"]["models.SocietyDashboardBootstrapData"];
            /** @example Society dashboard bootstrap fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SocietyDashboardBootstrapData": {
            dashboard?: components["schemas"]["models.SocietyDashboardBootstrapResponse"];
        };
        "models.SocietyDashboardBootstrapResponse": {
            claim_stats?: components["schemas"]["models.FlatClaimStatsResponse"];
            current_subscription?: components["schemas"]["models.SocietySubscriptionResponse"];
            flat_stats?: components["schemas"]["models.FlatStatsResponse"];
            member_stats?: components["schemas"]["models.SocietyDashboardMemberStatsResponse"];
            plan_ads?: components["schemas"]["models.PlanResponse"][];
            recent_pending_claims?: components["schemas"]["models.FlatClaimResponse"][];
            society?: components["schemas"]["models.SocietyResponse"];
            subscription_health?: components["schemas"]["models.SocietyDashboardSubscriptionHealthResponse"];
            subscription_usage?: components["schemas"]["models.SocietyDashboardSubscriptionUsageResponse"];
            visitor_daily_last_7_days?: components["schemas"]["models.VisitorDailyCountResponse"][];
            visitor_stats?: components["schemas"]["models.VisitorEntryStatsResponse"];
        };
        "models.SocietyDashboardMemberStatsResponse": {
            admins?: number;
            owners?: number;
            residents?: number;
            staff?: number;
            total_active_members?: number;
        };
        "models.SocietyDashboardQuotaUsageResponse": {
            limit?: number;
            percent?: number;
            remaining?: number;
            used?: number;
        };
        "models.SocietyDashboardSubscriptionHealthResponse": {
            days_until_expiry?: number;
            is_active?: boolean;
            is_expiring_soon?: boolean;
            lifecycle_label?: string;
        };
        "models.SocietyDashboardSubscriptionUsageResponse": {
            admins?: components["schemas"]["models.SocietyDashboardQuotaUsageResponse"];
            flats?: components["schemas"]["models.SocietyDashboardQuotaUsageResponse"];
            residents?: components["schemas"]["models.SocietyDashboardQuotaUsageResponse"];
            staff?: components["schemas"]["models.SocietyDashboardQuotaUsageResponse"];
        };
        "models.SocietyData": {
            society?: components["schemas"]["models.SocietyResponse"];
        };
        "models.SocietyDetailAPIResponse": {
            data?: components["schemas"]["models.SocietyDetailData"];
            /** @example Society fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SocietyDetailData": {
            society?: components["schemas"]["models.SocietyDetailResponse"];
        };
        "models.SocietyDetailResponse": {
            address_line1?: string;
            address_line2?: string;
            approved_at?: string;
            approved_by?: number;
            city?: string;
            country?: string;
            created_at?: string;
            created_by?: number;
            email?: string;
            id?: number;
            landmark?: string;
            members_count?: number;
            name?: string;
            phone_number?: string;
            pincode?: string;
            rejected_at?: string;
            rejected_by?: number;
            rejection_reason?: string;
            society_code?: string;
            state?: string;
            status?: components["schemas"]["models.SocietyStatus"];
            suspended_at?: string;
            suspended_by?: number;
            suspension_reason?: string;
            total_blocks?: number;
            total_flats?: number;
            updated_at?: string;
        };
        "models.SocietyFlatVisitorSettingRow": {
            approval_required?: boolean;
            block?: string;
            default_visit_duration_minutes?: number;
            flat_id?: number;
            flat_number?: string;
            is_enabled?: boolean;
            purpose?: components["schemas"]["models.VisitorPurpose"];
        };
        "models.SocietyFlatVisitorSettingsAPIResponse": {
            data?: components["schemas"]["models.SocietyFlatVisitorSettingsData"];
            /** @example Society flat visitor settings fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SocietyFlatVisitorSettingsData": {
            limit?: number;
            offset?: number;
            settings?: components["schemas"]["models.SocietyFlatVisitorSettingRow"][];
            total?: number;
        };
        "models.SocietyMemberAPIResponse": {
            data?: components["schemas"]["models.SocietyMemberData"];
            /** @example Member fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SocietyMemberData": {
            member?: components["schemas"]["models.SocietyMemberResponse"];
        };
        "models.SocietyMemberResponse": {
            created_at?: string;
            id?: number;
            invited_by?: number;
            joined_at?: string;
            remove_reason?: string;
            removed_at?: string;
            removed_by?: number;
            role?: components["schemas"]["models.SocietyMemberRole"];
            society_id?: number;
            status?: components["schemas"]["models.SocietyMemberStatus"];
            updated_at?: string;
            user_email?: string;
            user_full_name?: string;
            user_id?: number;
            user_phone?: string;
        };
        /** @enum {string} */
        "models.SocietyMemberRole": "owner" | "admin" | "staff" | "resident";
        /** @enum {string} */
        "models.SocietyMemberStatus": "pending" | "active" | "suspended" | "removed";
        "models.SocietyMemberSummaryAPIResponse": {
            data?: components["schemas"]["models.SocietyMemberSummaryData"];
            /** @example Member summary fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SocietyMemberSummaryData": {
            summary?: components["schemas"]["models.SocietyMemberSummaryResponse"];
        };
        "models.SocietyMemberSummaryResponse": {
            active_members?: number;
            admins?: number;
            owners?: number;
            pending_members?: number;
            removed_members?: number;
            residents?: number;
            staff?: number;
            suspended_members?: number;
            total_members?: number;
        };
        "models.SocietyOnboardingBootstrapAPIResponse": {
            data?: components["schemas"]["models.SocietyOnboardingBootstrapData"];
            /** @example Society onboarding bootstrap fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SocietyOnboardingBootstrapData": {
            onboarding?: components["schemas"]["models.SocietyOnboardingBootstrapResponse"];
        };
        "models.SocietyOnboardingBootstrapResponse": {
            flat_count?: number;
            has_active_subscription?: boolean;
            has_flats?: boolean;
            has_staff?: boolean;
            is_onboarded?: boolean;
            missing_steps?: string[];
            next_path?: string;
            society?: components["schemas"]["models.SocietyResponse"];
            staff_count?: number;
        };
        "models.SocietyReasonRequest": {
            reason: string;
        };
        "models.SocietyResponse": {
            address_line1?: string;
            address_line2?: string;
            approved_at?: string;
            approved_by?: number;
            city?: string;
            country?: string;
            created_at?: string;
            created_by?: number;
            email?: string;
            id?: number;
            landmark?: string;
            name?: string;
            phone_number?: string;
            pincode?: string;
            rejected_at?: string;
            rejected_by?: number;
            rejection_reason?: string;
            society_code?: string;
            state?: string;
            status?: components["schemas"]["models.SocietyStatus"];
            suspended_at?: string;
            suspended_by?: number;
            suspension_reason?: string;
            total_blocks?: number;
            total_flats?: number;
            updated_at?: string;
        };
        /** @enum {string} */
        "models.SocietyStatus": "pending" | "active" | "suspended" | "rejected";
        "models.SocietySubscriptionResponse": {
            activated_at?: string;
            activated_by?: number;
            billing_cycle?: components["schemas"]["models.BillingCycle"];
            cancellation_reason?: string;
            cancelled_at?: string;
            cancelled_by?: number;
            created_at?: string;
            created_by?: number;
            currency?: string;
            current_plan_code?: string;
            current_plan_name?: string;
            ends_at?: string;
            expired_at?: string;
            features?: {
                [key: string]: unknown;
            };
            id?: number;
            max_admins?: number;
            max_flats?: number;
            max_residents?: number;
            max_staff?: number;
            plan_code?: string;
            plan_id?: number;
            plan_name?: string;
            price_amount_paise?: number;
            society_code?: string;
            society_id?: number;
            society_name?: string;
            starts_at?: string;
            status?: components["schemas"]["models.SubscriptionStatus"];
            trial_ends_at?: string;
            updated_at?: string;
        };
        "models.SocietyVisitorSettingsAPIResponse": {
            data?: components["schemas"]["models.SocietyVisitorSettingsData"];
            /** @example Visitor settings fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SocietyVisitorSettingsData": {
            visitor_settings?: components["schemas"]["models.SocietyVisitorSettingsResponse"];
        };
        "models.SocietyVisitorSettingsResponse": {
            allow_guard_entry?: boolean;
            allow_guard_on_behalf_approval?: boolean;
            allow_public_qr_entry?: boolean;
            allow_resident_pre_approval?: boolean;
            approval_mode?: components["schemas"]["models.VisitorApprovalMode"];
            created_at?: string;
            default_visit_duration_minutes?: number;
            grace_period_minutes?: number;
            id?: number;
            is_active?: boolean;
            qr_expiry_minutes?: number;
            society_id?: number;
            updated_at?: string;
            updated_by?: number;
        };
        "models.SubmitFlatClaimRequest": {
            flat_id: number;
            metadata?: {
                [key: string]: unknown;
            };
            note?: string;
            requested_primary?: boolean;
            requested_role: components["schemas"]["models.FlatResidentRole"];
            society_id: number;
        };
        "models.SubscriptionAPIResponse": {
            data?: components["schemas"]["models.SubscriptionData"];
            /** @example Subscription fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SubscriptionData": {
            subscription?: components["schemas"]["models.SocietySubscriptionResponse"];
        };
        "models.SubscriptionStatsAPIResponse": {
            data?: components["schemas"]["models.SubscriptionStatsData"];
            /** @example Subscription stats fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SubscriptionStatsData": {
            stats?: components["schemas"]["models.SubscriptionStatsResponse"];
        };
        "models.SubscriptionStatsResponse": {
            active_subscriptions?: number;
            cancelled_subscriptions?: number;
            expired_subscriptions?: number;
            pending_subscriptions?: number;
            total_subscriptions?: number;
            trial_subscriptions?: number;
        };
        /** @enum {string} */
        "models.SubscriptionStatus": "pending" | "trial" | "active" | "expired" | "cancelled";
        "models.SubscriptionsAPIResponse": {
            data?: components["schemas"]["models.SubscriptionsData"];
            /** @example Subscriptions fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.SubscriptionsData": {
            subscriptions?: components["schemas"]["models.SocietySubscriptionResponse"][];
        };
        "models.TransferOwnershipRequest": {
            new_owner_user_id: number;
        };
        "models.UnregisterDeviceTokenRequest": {
            token: string;
        };
        "models.UpdateFlatRequest": {
            block?: string;
            flat_number?: string;
            floor?: string;
            is_active?: boolean;
            metadata?: {
                [key: string]: unknown;
            };
            status?: components["schemas"]["models.FlatStatus"];
        };
        "models.UpdateFlatResidentRoleRequest": {
            role: components["schemas"]["models.FlatResidentRole"];
        };
        "models.UpdateFlatVisitorSettingRequest": {
            approval_required?: boolean;
            default_visit_duration_minutes?: number;
            is_enabled?: boolean;
        };
        "models.UpdateGuardVisitorEntryRequest": {
            companion_details?: {
                [key: string]: unknown;
            }[];
            companions_count?: number;
            email?: string;
            flat_id?: number;
            full_name?: string;
            notes?: string;
            phone_number?: string;
            photo_url?: string;
            vehicle_number?: string;
            vehicle_type?: components["schemas"]["models.VisitorVehicleType"];
        };
        "models.UpdatePlanRequest": {
            billing_cycle?: components["schemas"]["models.BillingCycle"];
            code?: string;
            currency?: string;
            description?: string;
            features?: {
                [key: string]: unknown;
            };
            max_admins?: number;
            max_flats?: number;
            max_residents?: number;
            max_staff?: number;
            name?: string;
            price_amount_paise?: number;
        };
        "models.UpdateSocietyRequest": {
            address_line1?: string;
            address_line2?: string;
            city?: string;
            country?: string;
            email?: string;
            landmark?: string;
            metadata?: {
                [key: string]: unknown;
            };
            name?: string;
            phone_number?: string;
            pincode?: string;
            state?: string;
            total_blocks?: number;
            total_flats?: number;
        };
        "models.UpdateSocietyVisitorSettingsRequest": {
            allow_guard_entry?: boolean;
            allow_guard_on_behalf_approval?: boolean;
            allow_public_qr_entry?: boolean;
            allow_resident_pre_approval?: boolean;
            approval_mode?: components["schemas"]["models.VisitorApprovalMode"];
            default_visit_duration_minutes?: number;
            grace_period_minutes?: number;
            is_active?: boolean;
            qr_expiry_minutes?: number;
        };
        "models.UpdateUserRequest": {
            avatar_url?: string;
            date_of_birth?: string;
            first_name?: string;
            /** @enum {string} */
            gender?: "male" | "female" | "other" | "prefer_not_to_say";
            language?: string;
            last_name?: string;
            phone_number?: string;
            timezone?: string;
        };
        "models.UserData": {
            user?: components["schemas"]["models.UserResponse"];
        };
        "models.UserMessageData": {
            /** @example 123456 */
            dev_otp?: string;
            /** @example Please verify your email using the OTP sent to your email address */
            message?: string;
            user?: components["schemas"]["models.UserResponse"];
        };
        "models.UserResponse": {
            auth_provider?: components["schemas"]["models.AuthProvider"];
            avatar_url?: string;
            blocked_reason?: string;
            created_at?: string;
            date_of_birth?: string;
            email?: string;
            email_verified?: boolean;
            first_name?: string;
            full_name?: string;
            gender?: string;
            global_role?: components["schemas"]["models.GlobalRole"];
            id?: number;
            is_active?: boolean;
            is_blocked?: boolean;
            language?: string;
            last_login_at?: string;
            last_name?: string;
            phone_number?: string;
            phone_verified?: boolean;
            timezone?: string;
            updated_at?: string;
        };
        "models.VerifyOTPAPIResponse": {
            data?: components["schemas"]["models.UserData"];
            /** @example Email verified successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VerifyOTPRequest": {
            email: string;
            otp: string;
        };
        /** @enum {string} */
        "models.VisitorApprovalMode": "mandatory" | "optional" | "hybrid";
        "models.VisitorDailyCountResponse": {
            count?: number;
            date?: string;
        };
        "models.VisitorEntriesAPIResponse": {
            data?: components["schemas"]["models.VisitorEntriesData"];
            /** @example Visitor entries fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorEntriesData": {
            entries?: components["schemas"]["models.VisitorEntry"][];
            limit?: number;
            offset?: number;
            total?: number;
        };
        "models.VisitorEntry": {
            approved_at?: string;
            approved_by?: number;
            approver_name?: string;
            approver_phone_number?: string;
            auto_closed_at?: string;
            checked_in_at?: string;
            checked_out_at?: string;
            companion_details?: {
                [key: string]: unknown;
            }[];
            companions_count?: number;
            created_at?: string;
            created_by?: number;
            delivery_partner?: string;
            expected_at?: string;
            expected_checkout_at?: string;
            flat?: components["schemas"]["models.VisitorFlatSummary"];
            flat_id?: number;
            handled_by_guard_id?: number;
            handled_by_guard_name?: string;
            handled_by_guard_phone?: string;
            id?: number;
            invite_id?: number;
            metadata?: {
                [key: string]: unknown;
            };
            notes?: string;
            purpose?: components["schemas"]["models.VisitorPurpose"];
            qr_expires_at?: string;
            qr_used_at?: string;
            rejected_by?: number;
            rejection_reason?: string;
            service_provider?: string;
            society_id?: number;
            source?: components["schemas"]["models.VisitorEntrySource"];
            status?: components["schemas"]["models.VisitorStatus"];
            updated_at?: string;
            vehicle_number?: string;
            vehicle_type?: components["schemas"]["models.VisitorVehicleType"];
            visitor?: components["schemas"]["models.VisitorSummary"];
            visitor_id?: number;
        };
        "models.VisitorEntryAPIResponse": {
            data?: components["schemas"]["models.VisitorEntryData"];
            /** @example Visitor entry fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorEntryDailyStatsAPIResponse": {
            data?: components["schemas"]["models.VisitorEntryDailyStatsData"];
            /** @example Visitor entry daily stats fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorEntryDailyStatsData": {
            stats?: components["schemas"]["models.VisitorEntryDailyStatsResponse"];
        };
        "models.VisitorEntryDailyStatsResponse": {
            daily?: components["schemas"]["models.VisitorDailyCountResponse"][];
            days?: number;
            metric?: string;
            timezone?: string;
            total?: number;
        };
        "models.VisitorEntryData": {
            entry?: components["schemas"]["models.VisitorEntry"];
        };
        "models.VisitorEntryEvent": {
            actor_user_id?: number;
            created_at?: string;
            event_type?: components["schemas"]["models.VisitorEventType"];
            id?: number;
            message?: string;
            metadata?: {
                [key: string]: unknown;
            };
            society_id?: number;
            visitor_entry_id?: number;
        };
        "models.VisitorEntryEventsAPIResponse": {
            data?: components["schemas"]["models.VisitorEntryEventsData"];
            /** @example Visitor entry events fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorEntryEventsData": {
            events?: components["schemas"]["models.VisitorEntryEvent"][];
        };
        "models.VisitorEntryMutationAPIResponse": {
            data?: components["schemas"]["models.VisitorEntryMutationResponse"];
            /** @example Visitor entry created successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorEntryMutationResponse": {
            entry?: components["schemas"]["models.VisitorEntry"];
            qr?: components["schemas"]["models.QRTokenResponse"];
        };
        "models.VisitorEntryOptionsAPIResponse": {
            data?: components["schemas"]["models.VisitorEntryOptionsData"];
            /** @example Visitor entry options fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorEntryOptionsBlock": {
            block?: string;
            flats?: components["schemas"]["models.VisitorEntryOptionsFlat"][];
        };
        "models.VisitorEntryOptionsData": {
            options?: components["schemas"]["models.VisitorEntryOptionsResponse"];
        };
        "models.VisitorEntryOptionsFlat": {
            block?: string;
            flat_number?: string;
            floor?: string;
            id?: number;
        };
        "models.VisitorEntryOptionsResponse": {
            blocks?: components["schemas"]["models.VisitorEntryOptionsBlock"][];
            flats?: components["schemas"]["models.VisitorEntryOptionsFlat"][];
            has_more?: boolean;
            purposes?: components["schemas"]["models.VisitorPurpose"][];
        };
        /** @enum {string} */
        "models.VisitorEntrySource": "resident_link" | "public_qr" | "guard_entry" | "quick_link";
        "models.VisitorEntryStatsAPIResponse": {
            data?: components["schemas"]["models.VisitorEntryStatsData"];
            /** @example Visitor entry stats fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorEntryStatsData": {
            stats?: components["schemas"]["models.VisitorEntryStatsResponse"];
        };
        "models.VisitorEntryStatsResponse": {
            auto_closed_today?: number;
            checked_out_in_range?: number;
            checked_out_today?: number;
            pending_approvals?: number;
            rejected_today?: number;
            today_visitors?: number;
            visitors_inside?: number;
        };
        /** @enum {string} */
        "models.VisitorEventType": "created" | "approved" | "rejected" | "checked_in" | "checked_out" | "cancelled" | "expired" | "auto_closed" | "qr_generated" | "qr_used" | "guard_approved_on_behalf";
        "models.VisitorFlatSummary": {
            block?: string;
            flat_number?: string;
            floor?: string;
            id?: number;
        };
        "models.VisitorFormRequest": {
            companion_details?: {
                [key: string]: unknown;
            }[];
            companions_count?: number;
            delivery_partner?: string;
            email?: string;
            expected_at?: string;
            expected_checkout_at?: string;
            flat_id?: number;
            full_name?: string;
            metadata?: {
                [key: string]: unknown;
            };
            notes?: string;
            phone_number?: string;
            photo_url?: string;
            purpose?: components["schemas"]["models.VisitorPurpose"];
            service_provider?: string;
            vehicle_number?: string;
            vehicle_type?: components["schemas"]["models.VisitorVehicleType"];
        };
        "models.VisitorInvite": {
            created_at?: string;
            created_by?: number;
            expires_at?: string;
            flat_id?: number;
            id?: number;
            metadata?: {
                [key: string]: unknown;
            };
            purpose?: components["schemas"]["models.VisitorPurpose"];
            society_id?: number;
            status?: components["schemas"]["models.VisitorInviteStatus"];
            updated_at?: string;
            used_at?: string;
        };
        "models.VisitorInviteAPIResponse": {
            data?: components["schemas"]["models.VisitorInviteData"];
            /** @example Visitor invite fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorInviteData": {
            entry?: components["schemas"]["models.VisitorEntry"];
            invite?: components["schemas"]["models.PublicVisitorInviteView"];
            qr?: components["schemas"]["models.QRTokenResponse"];
            /** @example form */
            view?: string;
        };
        "models.VisitorInviteEntry": {
            checked_in_at?: string;
            checked_out_at?: string;
            id?: number;
            qr_expires_at?: string;
            qr_used_at?: string;
            status?: components["schemas"]["models.VisitorStatus"];
        };
        "models.VisitorInviteHistoryAPIResponse": {
            data?: components["schemas"]["models.VisitorInviteHistoryData"];
            /** @example Visitor invites fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorInviteHistoryData": {
            invites?: components["schemas"]["models.VisitorInviteHistoryItem"][];
            limit?: number;
            offset?: number;
            total?: number;
        };
        "models.VisitorInviteHistoryItem": {
            created_at?: string;
            created_by?: number;
            entry?: components["schemas"]["models.VisitorInviteEntry"];
            expires_at?: string;
            flat_id?: number;
            id?: number;
            metadata?: {
                [key: string]: unknown;
            };
            purpose?: components["schemas"]["models.VisitorPurpose"];
            short_link?: components["schemas"]["models.ShortLinkResponse"];
            society_id?: number;
            status?: components["schemas"]["models.VisitorInviteStatus"];
            updated_at?: string;
            used_at?: string;
            visitor?: components["schemas"]["models.VisitorSummary"];
        };
        "models.VisitorInviteHistoryItemAPIResponse": {
            data?: components["schemas"]["models.VisitorInviteHistoryItemData"];
            /** @example Visitor invite fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorInviteHistoryItemData": {
            invite?: components["schemas"]["models.VisitorInviteHistoryItem"];
        };
        /** @enum {string} */
        "models.VisitorInviteStatus": "active" | "used" | "expired" | "cancelled";
        "models.VisitorInviteTokenAPIResponse": {
            data?: components["schemas"]["models.VisitorInviteTokenData"];
            /** @example Visitor invite created successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorInviteTokenData": {
            invite?: components["schemas"]["models.VisitorInvite"];
            link?: components["schemas"]["models.ShortLinkResponse"];
            short_link?: components["schemas"]["models.ShortLinkResponse"];
            token?: components["schemas"]["models.QRTokenResponse"];
        };
        "models.VisitorPendingEntriesAPIResponse": {
            data?: components["schemas"]["models.VisitorPendingEntriesData"];
            /** @example Pending visitor approvals fetched successfully */
            message?: string;
            /** @example true */
            success?: boolean;
        };
        "models.VisitorPendingEntriesData": {
            entries?: components["schemas"]["models.VisitorPendingEntry"][];
            limit?: number;
            offset?: number;
            total?: number;
        };
        "models.VisitorPendingEntry": {
            approved_at?: string;
            approved_by?: number;
            approver_name?: string;
            approver_phone_number?: string;
            auto_closed_at?: string;
            checked_in_at?: string;
            checked_out_at?: string;
            companion_details?: {
                [key: string]: unknown;
            }[];
            companions_count?: number;
            created_at?: string;
            created_by?: number;
            delivery_partner?: string;
            expected_at?: string;
            expected_checkout_at?: string;
            flat?: components["schemas"]["models.VisitorFlatSummary"];
            flat_id?: number;
            handled_by_guard_id?: number;
            handled_by_guard_name?: string;
            handled_by_guard_phone?: string;
            id?: number;
            invite_id?: number;
            metadata?: {
                [key: string]: unknown;
            };
            notes?: string;
            primary_resident_id?: number;
            primary_resident_name?: string;
            purpose?: components["schemas"]["models.VisitorPurpose"];
            qr_expires_at?: string;
            qr_used_at?: string;
            rejected_by?: number;
            rejection_reason?: string;
            service_provider?: string;
            society_id?: number;
            source?: components["schemas"]["models.VisitorEntrySource"];
            status?: components["schemas"]["models.VisitorStatus"];
            updated_at?: string;
            vehicle_number?: string;
            vehicle_type?: components["schemas"]["models.VisitorVehicleType"];
            visitor?: components["schemas"]["models.VisitorSummary"];
            visitor_id?: number;
            waiting_since?: string;
        };
        /** @enum {string} */
        "models.VisitorPurpose": "guest" | "delivery" | "cab" | "service" | "maintenance" | "staff" | "other";
        /** @enum {string} */
        "models.VisitorStatus": "waiting_approval" | "approved" | "rejected" | "checked_in" | "checked_out" | "cancelled" | "expired" | "auto_closed";
        "models.VisitorSummary": {
            email?: string;
            full_name?: string;
            phone_number?: string;
            photo_url?: string;
        };
        /** @enum {string} */
        "models.VisitorVehicleType": "bike" | "car" | "auto" | "cab" | "truck" | "other";
        "models.NotificationDetailAPIResponse": {
            data?: components["schemas"]["models.Notification"];
            message?: string;
            success?: boolean;
        };
    };
    responses: never;
    parameters: never;
    requestBodies: {
        /** @description Visitor QR token */
        "models.QRTokenRequest": {
            content: {
                "application/json": components["schemas"]["models.QRTokenRequest"];
            };
        };
        /** @description Guard approval options */
        "models.GuardApproveEntryRequest": {
            content: {
                "application/json": components["schemas"]["models.GuardApproveEntryRequest"];
            };
        };
        /** @description Visitor invite request */
        "models.CreateVisitorInviteRequest": {
            content: {
                "application/json": components["schemas"]["models.CreateVisitorInviteRequest"];
            };
        };
        putV1AuthProfileAvatar: {
            content: {
                "multipart/form-data": {
                    /**
                     * Format: binary
                     * @description JPEG or PNG, maximum 5 MiB and 25 megapixels
                     */
                    file: string;
                };
            };
        };
        /** @description Join request */
        "models.JoinFlatMemberInviteRequest": {
            content: {
                "application/json": components["schemas"]["models.JoinFlatMemberInviteRequest"];
            };
        };
        /** @description Visitor details */
        "models.VisitorFormRequest": {
            content: {
                "application/json": components["schemas"]["models.VisitorFormRequest"];
            };
        };
    };
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export type operations = Record<string, never>;
