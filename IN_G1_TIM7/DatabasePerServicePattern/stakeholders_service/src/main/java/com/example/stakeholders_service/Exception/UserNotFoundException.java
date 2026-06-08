package com.example.stakeholders_service.Exception;

import javax.management.RuntimeErrorException;

public class UserNotFoundException extends RuntimeException {

    public UserNotFoundException(Long id) {
        super("Korisnik sa ID-em " + id + " nije pronađen");
    }

}
