package com.example.stakeholders_service.Exception;

public class NonAdminProfileAccessException extends RuntimeException{

    public NonAdminProfileAccessException() { super("Samo administratori mogu videti sve profile."); }
}
