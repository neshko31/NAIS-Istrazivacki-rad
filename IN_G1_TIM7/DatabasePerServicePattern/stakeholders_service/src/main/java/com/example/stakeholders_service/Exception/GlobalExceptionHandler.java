package com.example.stakeholders_service.Exception;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(UserNotFoundException.class)
    public ResponseEntity<String> handleUserNotFound(UserNotFoundException ex) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(ex.getMessage());
    }

    @ExceptionHandler(AdminProfileAccessException.class)
    public ResponseEntity<String> handleAdminProfileAccess(AdminProfileAccessException ad){
        return ResponseEntity.status(HttpStatus.FORBIDDEN).body(ad.getMessage());
    }

    @ExceptionHandler(AdminBlocksOtherAdminException.class)
    public ResponseEntity<String> handleAdminBlocksOtherAdmin(AdminBlocksOtherAdminException ex) {
        return ResponseEntity.status(HttpStatus.FORBIDDEN).body(ex.getMessage());
    }

    @ExceptionHandler(BlockedProfileAccessException.class)
    public ResponseEntity<String> handleBlockedProfileAccess(BlockedProfileAccessException ex) {
        return ResponseEntity.status(HttpStatus.FORBIDDEN).body(ex.getMessage());
    }

    @ExceptionHandler(NonAdminProfileAccessException.class)
    public ResponseEntity<String> handleUserNotFound(NonAdminProfileAccessException ex) {
        return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(ex.getMessage());
    }

    @ExceptionHandler(NonAdminProfileBlockingException.class)
    public ResponseEntity<String> handleUserNotFound(NonAdminProfileBlockingException ex) {
        return ResponseEntity.status(HttpStatus.FORBIDDEN).body(ex.getMessage());
    }

    @ExceptionHandler(UserNotTouristException.class)
    public ResponseEntity<String> handleUserNotTourst(UserNotTouristException ex) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(ex.getMessage());
    }
}
