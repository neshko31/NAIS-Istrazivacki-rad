package com.example.stakeholders_service.Exception;

public class NonAdminProfileBlockingException extends RuntimeException{

    public NonAdminProfileBlockingException() { super("Samo administratori smeju da blokiraju korisnike."); }
}
