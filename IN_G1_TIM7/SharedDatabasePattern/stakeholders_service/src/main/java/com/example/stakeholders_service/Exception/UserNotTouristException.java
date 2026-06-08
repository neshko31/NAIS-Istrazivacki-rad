package com.example.stakeholders_service.Exception;

public class UserNotTouristException extends RuntimeException {
    public UserNotTouristException() { super("Korisnik kojem ste pristupili nije turista!"); }
}
