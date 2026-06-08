package com.example.stakeholders_service.Exception;

public class BlockedProfileBlockingException extends RuntimeException{

    public BlockedProfileBlockingException() { super("Nije moguce blokirati vec blokiran nalog korisnika."); }
}
