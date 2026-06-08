package com.example.stakeholders_service.Exception;

public class AdminProfileAccessException extends RuntimeException{

    public AdminProfileAccessException() { super("Administratori nemaju profil."); }
}
