# Simplified Building Extraction Process Sequence Diagram

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant Backend
    participant Database
    
    User->>Frontend: Upload remote sensing image
    Frontend->>Backend: POST /api/extraction with image
    Backend->>Backend: Authenticate user (JWT)
    Backend->>Backend: Save uploaded image
    
    User->>Frontend: Select extraction model
    Frontend->>Backend: Process with selected model
    
    Note over Backend: Run building extraction
    Backend->>Backend: Process image with selected model
    Backend->>Backend: Generate building masks
    Backend->>Backend: Save extraction results
    
    Backend->>Database: Store extraction metadata
    Backend-->>Frontend: Return result URLs
    Frontend-->>User: Display extraction results
    
    User->>Frontend: Save project
    Frontend->>Backend: POST /api/projects
    Backend->>Database: Save project data
    Backend-->>Frontend: Return success
    Frontend-->>User: Confirm project saved
```

## Diagram Explanation

This simplified sequence diagram illustrates the essential flow of the building extraction process:

1. User uploads a remote sensing image through the frontend
2. User selects a model for building extraction
3. Backend processes the image using the selected model
4. Results are displayed to the user
5. User can save the project for future reference

The diagram focuses on the core user interactions and essential system operations without detailed implementation specifics. 