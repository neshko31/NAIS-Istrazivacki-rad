package com.example.stakeholders_service.Exception;

public class BlockedProfileAccessException extends RuntimeException{

    public BlockedProfileAccessException() { super("Nije moguce pristupiti blokiranom nalogu."); }
}
