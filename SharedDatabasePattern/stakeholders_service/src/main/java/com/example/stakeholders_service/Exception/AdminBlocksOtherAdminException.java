package com.example.stakeholders_service.Exception;

public class AdminBlocksOtherAdminException extends RuntimeException{

    public AdminBlocksOtherAdminException() { super("Administratori se ne mogu blokirati medjusobno."); }
}
