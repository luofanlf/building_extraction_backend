# JWT Login Sequence Diagram

```mermaid
sequenceDiagram
    participant Client
    participant Server as Server (Gin Framework)
    participant AuthService as Authentication Service
    participant Database as Database (PostgreSQL)
    
    Client->>Server: POST /api/login with username and password
    Server->>AuthService: Verify user credentials
    AuthService->>Database: Query user information
    Database-->>AuthService: Return user data
    
    alt User does not exist
        AuthService-->>Server: Return user not found error
        Server-->>Client: 401 Unauthorized (Authentication failed)
    else Incorrect password
        AuthService->>AuthService: Validate password hash
        AuthService-->>Server: Return password incorrect error
        Server-->>Client: 401 Unauthorized (Authentication failed)
    else Authentication successful
        AuthService->>AuthService: Generate JWT token
        AuthService-->>Server: Return JWT token
        Server-->>Client: 200 OK with JWT token
    end
    
    Client->>Client: Store JWT token
    
    Note over Client,Server: Subsequent requests
    
    Client->>Server: Request protected resource with JWT in Authorization header
    Server->>Server: Validate JWT (AuthMiddleware)
    
    alt Valid JWT
        Server->>Server: Parse JWT to get user information
        Server->>Server: Add user information to context
        Server->>Server: Continue processing request
        Server-->>Client: Return requested data
    else Invalid or expired JWT
        Server-->>Client: 401 Unauthorized
    end
```

## Diagram Explanation

1. **Login Process**:
   - Client sends username and password to the `/api/login` endpoint
   - Server passes credentials to the Authentication Service for validation
   - Authentication Service queries user information from the database
   - If validation is successful, a JWT token is generated and returned to the client
   - Client stores the JWT token for subsequent requests

2. **Protected Resource Access Flow**:
   - Client includes the JWT token in the Authorization header of the request
   - Server's AuthMiddleware intercepts the request to validate the JWT
   - If the JWT is valid, user information is extracted from the token and added to the request context
   - If the JWT is invalid or expired, a 401 Unauthorized error is returned

3. **JWT Structure**:
   - Header: Contains token type and encryption algorithm used
   - Payload: Contains user information (such as user ID, roles, etc.) and expiration time
   - Signature: Generated using the server's secret key to ensure the token has not been tampered with

This sequence diagram illustrates the main interaction steps in the JWT authentication process for your building extraction backend project. 